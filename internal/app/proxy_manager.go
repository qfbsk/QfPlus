package app

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func (a *App) getMihomoExecutable() (string, error) {
	name := "mihomo"
	if isWindows() {
		name = "mihomo.exe"
	}
	if path := findCoreFile(name); path != "" {
		return path, nil
	}
	return "", fmt.Errorf("mihomo executable not found at %s", filepath.Join(a.getCoreDir(), name))
}

func isWindows() bool {
	return os.PathSeparator == '\\'
}

func (a *App) proxyPIDFile() string {
	return filepath.Join(a.proxyDir(), "mihomo.pid")
}

func (a *App) ensureAPISecret() (string, error) {
	config, err := a.readAppConfig()
	if err != nil {
		return "", err
	}
	if config.Proxy.APISecret != "" {
		return config.Proxy.APISecret, nil
	}
	buf := make([]byte, 16)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	secret := hex.EncodeToString(buf)
	config.Proxy.APISecret = secret
	if config.Proxy.MixedPort == 0 {
		config.Proxy.MixedPort = 17890
	}
	if config.Proxy.APIPort == 0 {
		config.Proxy.APIPort = 19090
	}
	return secret, a.saveAppConfig(config)
}

func (a *App) startMihomo() error {
	a.proxyMu.Lock()
	defer a.proxyMu.Unlock()

	if a.proxyCmd != nil && a.proxyCmd.Process != nil {
		return nil
	}

	a.killOrphanMihomo()

	exePath, err := a.getMihomoExecutable()
	if err != nil {
		return err
	}

	config, err := a.readAppConfig()
	if err != nil {
		return err
	}
	mixedPort := config.Proxy.MixedPort
	if mixedPort == 0 {
		mixedPort = 17890
	}
	apiPort := config.Proxy.APIPort
	if apiPort == 0 {
		apiPort = 19090
	}

	mixedPort = a.findFreePort(mixedPort)
	apiPort = a.findFreePort(apiPort)
	if mixedPort != config.Proxy.MixedPort || apiPort != config.Proxy.APIPort {
		config.Proxy.MixedPort = mixedPort
		config.Proxy.APIPort = apiPort
		_ = a.saveAppConfig(config)
		a.rewriteProxyConfigPorts(mixedPort, apiPort)
	}

	cmd := exec.Command(exePath, "-d", a.proxyDir())
	hideWindow(cmd)
	cmd.Stdout = nil
	cmd.Stderr = nil
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start mihomo: %w", err)
	}

	pidFile := a.proxyPIDFile()
	_ = os.WriteFile(pidFile, []byte(strconv.Itoa(cmd.Process.Pid)), 0644)

	secret := config.Proxy.APISecret
	if !a.waitForAPI(apiPort, secret, 10*time.Second) {
		_ = cmd.Process.Kill()
		os.Remove(pidFile)
		return fmt.Errorf("mihomo API did not become ready within 10s on port %d", apiPort)
	}

	a.proxyCmd = cmd
	go a.waitMihomoExit(cmd)
	return nil
}

func (a *App) waitMihomoExit(cmd *exec.Cmd) {
	_ = cmd.Wait()
	a.proxyMu.Lock()
	if a.proxyCmd == cmd {
		a.proxyCmd = nil
	}
	a.proxyMu.Unlock()
}

func (a *App) stopMihomo() {
	a.proxyMu.Lock()
	defer a.proxyMu.Unlock()

	if a.proxyCmd != nil && a.proxyCmd.Process != nil {
		_ = a.proxyCmd.Process.Kill()
		_ = a.proxyCmd.Wait()
		a.proxyCmd = nil
	}
	os.Remove(a.proxyPIDFile())
}

func (a *App) killOrphanMihomo() {
	data, err := os.ReadFile(a.proxyPIDFile())
	if err != nil {
		return
	}
	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil || pid <= 0 {
		return
	}
	if proc, err := os.FindProcess(pid); err == nil {
		_ = proc.Kill()
	}
	os.Remove(a.proxyPIDFile())
	time.Sleep(500 * time.Millisecond)
}

func (a *App) isMihomoRunning() bool {
	a.proxyMu.Lock()
	cmd := a.proxyCmd
	a.proxyMu.Unlock()
	return cmd != nil && cmd.Process != nil
}

func (a *App) findFreePort(start int) int {
	for port := start; port < start+100; port++ {
		ln, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", port))
		if err == nil {
			ln.Close()
			return port
		}
	}
	return start
}

func (a *App) waitForAPI(port int, secret string, timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	client := &http.Client{Timeout: 2 * time.Second}
	url := fmt.Sprintf("http://127.0.0.1:%d/version", port)
	for time.Now().Before(deadline) {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer "+secret)
		resp, err := client.Do(req)
		if err == nil {
			resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				return true
			}
		}
		time.Sleep(300 * time.Millisecond)
	}
	return false
}

func (a *App) rewriteProxyConfigPorts(mixedPort, apiPort int) {
	data, err := os.ReadFile(a.proxyConfigFile())
	if err != nil {
		return
	}
	content := string(data)
	content = replaceYAMLInt(content, "mixed-port", mixedPort)
	content = replaceYAMLStr(content, "external-controller", fmt.Sprintf("127.0.0.1:%d", apiPort))
	_ = os.WriteFile(a.proxyConfigFile(), []byte(content), 0644)
}

func replaceYAMLInt(content, key string, value int) string {
	return replaceYAMLValue(content, key, strconv.Itoa(value))
}

func replaceYAMLStr(content, key, value string) string {
	return replaceYAMLValue(content, key, value)
}

func replaceYAMLValue(content, key, value string) string {
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, key+":") {
			lines[i] = key + ": " + value
			break
		}
	}
	return strings.Join(lines, "\n")
}

func (a *App) vfoxProxyEnvURL() (string, bool) {
	config, err := a.readAppConfig()
	if err != nil || !config.Proxy.Enabled {
		return "", false
	}
	if !a.isMihomoRunning() {
		return "", false
	}
	port := config.Proxy.MixedPort
	if port == 0 {
		port = 17890
	}
	return fmt.Sprintf("http://127.0.0.1:%d", port), true
}
