package app

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"context"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	stdruntime "runtime"
)

const (
	coreAtomURL       = "https://github.com/version-fox/vfox/releases.atom"
	coreDownloadBase  = "https://github.com/version-fox/vfox/releases/download/"
	coreCheckInterval = 24 * time.Hour
)

var coreVersionRegex = regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)

// coreExecutableName returns the engine binary name for the current platform.
func coreExecutableName() string {
	return getVfoxExeName()
}

// bundledCoreDir is the directory shipped with the app and never overwritten
// the updater, so a bad download can always be rolled back.
func (a *App) bundledCoreDir() string {
	return a.getCoreDir()
}

func (a *App) activeCoreDir() string {
	if dir := markerCoreDir(); dir != "" {
		return dir
	}
	return a.getCoreDir()
}

// coreVersionForDir reports the engine version living in dir by running it.
func (a *App) coreVersionForDir(dir string) string {
	exe := filepath.Join(dir, coreExecutableName())
	if !coreFileExists(exe) {
		return ""
	}
	out, err := a.runCoreBinaryVersion(exe)
	if err != nil {
		return ""
	}
	return out
}

func (a *App) runCoreBinaryVersion(exe string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, exe, "--version")
	hideWindow(cmd)
	// Same sanitization as regular vfox calls: without __VFOX_SHELL the engine
	// tries to spawn an interactive shell and can deadlock against the GUI.
	cmd.Env = a.getCleanedEnvForVfox()
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}
	match := coreVersionRegex.FindStringSubmatch(string(out))
	if match == nil {
		return "", fmt.Errorf("unable to parse version from %q", strings.TrimSpace(string(out)))
	}
	return match[1], nil
}

// getCurrentCoreVersion asks the active engine for its own version so the GUI
// never reports a stale number after an in-place upgrade.
func (a *App) getCurrentCoreVersion() (string, error) {
	out, err := a.runVfoxCommand("--version")
	if err != nil {
		return "", err
	}
	match := coreVersionRegex.FindStringSubmatch(out)
	if match == nil {
		return "", fmt.Errorf("unable to parse core version from %q", strings.TrimSpace(out))
	}
	return match[1], nil
}

// localCoreVersions lists engine versions already downloaded to the
// user-writable store, newest first.
func localCoreVersions() []string {
	root := localCoreRoot()
	entries, err := os.ReadDir(root)
	if err != nil {
		return nil
	}
	var versions []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		if coreFileExists(filepath.Join(root, entry.Name(), coreExecutableName())) {
			versions = append(versions, entry.Name())
		}
	}
	sort.SliceStable(versions, func(i, j int) bool {
		return compareCoreVersions(versions[i], versions[j]) > 0
	})
	return versions
}

// coreHTTPClient routes the app's own outbound requests through mihomo when
// built-in proxy is live. Environment-variable injection only covers child
// processes, so the transport has to be configured explicitly here.
func (a *App) coreHTTPClient(timeout time.Duration) *http.Client {
	transport := http.DefaultTransport.(*http.Transport).Clone()
	if proxyURL, ok := a.vfoxProxyEnvURL(); ok {
		if parsed, err := url.Parse(proxyURL); err == nil {
			transport.Proxy = http.ProxyURL(parsed)
		}
	}
	return &http.Client{Timeout: timeout, Transport: transport}
}

type atomFeed struct {
	Entries []atomEntry `xml:"entry"`
}

type atomEntry struct {
	ID      string `xml:"id"`
	Title   string `xml:"title"`
	Updated string `xml:"updated"`
	Link    struct {
		Href string `xml:"href,attr"`
	} `xml:"link"`
	Content string `xml:"content"`
}

// fetchCoreReleases reads the public Atom feed instead of the REST API: the
// API is rate-limited per IP and the shared proxy exit IP is often exhausted.
func (a *App) fetchCoreReleases() ([]CoreRelease, error) {
	req, err := http.NewRequest("GET", a.applyGitHubSource(coreAtomURL), nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "QfPlus/1.0")
	resp, err := a.coreHTTPClient(30 * time.Second).Do(req)
	if err != nil {
		return nil, fmt.Errorf("无法连接 GitHub：%w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub 返回状态码 %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, 2<<20))
	if err != nil {
		return nil, err
	}
	releases, err := parseCoreAtomFeed(body)
	if err != nil {
		return nil, err
	}
	for i := range releases {
		releases[i].URL = a.applyGitHubSource(releases[i].URL)
	}
	return releases, nil
}

// parseCoreAtomFeed converts the GitHub releases Atom XML into release DTOs,
// newest first, dropping prerelease/duplicate entries.
func parseCoreAtomFeed(body []byte) ([]CoreRelease, error) {
	var feed atomFeed
	if err := xml.Unmarshal(body, &feed); err != nil {
		return nil, fmt.Errorf("无法解析版本列表：%w", err)
	}

	var releases []CoreRelease
	seen := map[string]bool{}
	for _, entry := range feed.Entries {
		match := coreVersionRegex.FindStringSubmatch(entry.Title)
		if match == nil {
			match = coreVersionRegex.FindStringSubmatch(entry.ID)
			if match == nil {
				continue
			}
		}
		version := match[1]
		if seen[version] {
			continue
		}
		seen[version] = true
		releases = append(releases, CoreRelease{
			Version: version,
			Title:   strings.TrimSpace(entry.Title),
			Date:    formatReleaseDate(entry.Updated),
			Notes:   plainTextReleaseNotes(entry.Content),
			URL:     entry.Link.Href,
		})
	}
	if len(releases) == 0 {
		return nil, fmt.Errorf("未从 GitHub 获取到任何版本")
	}
	sort.SliceStable(releases, func(i, j int) bool {
		return compareCoreVersions(releases[i].Version, releases[j].Version) > 0
	})
	return releases, nil
}

func formatReleaseDate(raw string) string {
	if t, err := time.Parse(time.RFC3339, raw); err == nil {
		return t.Format("2006-01-02")
	}
	return strings.TrimSpace(raw)
}

var (
	htmlTagRegex     = regexp.MustCompile(`<[^>]*>`)
	htmlNewlineRegex = regexp.MustCompile(`\n{3,}`)
	htmlSpaceRegex   = regexp.MustCompile(`[ \t]+`)
)

// plainTextReleaseNotes turns the escaped HTML body of an Atom entry into a
// compact changelog block that renders well inside the settings card.
func plainTextReleaseNotes(content string) string {
	text := strings.ReplaceAll(content, "<br>", "\n")
	text = strings.ReplaceAll(text, "<br/>", "\n")
	text = strings.ReplaceAll(text, "</p>", "\n")
	text = strings.ReplaceAll(text, "</li>", "\n")
	text = strings.ReplaceAll(text, "<li>", "- ")
	text = htmlTagRegex.ReplaceAllString(text, "")
	text = unescapeHTMLEntities(text)
	lines := strings.Split(text, "\n")
	var kept []string
	for _, line := range lines {
		line = htmlSpaceRegex.ReplaceAllString(strings.TrimSpace(line), " ")
		if line != "" {
			kept = append(kept, line)
		}
	}
	notes := strings.Join(kept, "\n")
	notes = htmlNewlineRegex.ReplaceAllString(notes, "\n\n")
	if len(notes) > 4000 {
		notes = notes[:4000] + "…"
	}
	return notes
}

func unescapeHTMLEntities(text string) string {
	replacer := strings.NewReplacer(
		"&lt;", "<", "&gt;", ">", "&quot;", `"`, "&#39;", "'", "&apos;", "'",
		"&nbsp;", " ", "&amp;", "&",
	)
	return replacer.Replace(text)
}

// compareCoreVersions returns >0 when a is newer than b.
func compareCoreVersions(a, b string) int {
	pa := strings.Split(strings.TrimPrefix(a, "v"), ".")
	pb := strings.Split(strings.TrimPrefix(b, "v"), ".")
	for i := 0; i < len(pa) || i < len(pb); i++ {
		var na, nb int
		if i < len(pa) {
			na = numericPrefix(pa[i])
		}
		if i < len(pb) {
			nb = numericPrefix(pb[i])
		}
		if na != nb {
			if na > nb {
				return 1
			}
			return -1
		}
	}
	return 0
}

func numericPrefix(s string) int {
	digits := strings.TrimRightFunc(s, func(r rune) bool { return !('0' <= r && r <= '9') })
	if digits == "" {
		return 0
	}
	n, err := strconv.Atoi(digits)
	if err != nil {
		return 0
	}
	return n
}

// coreAssetName mirrors the goreleaser archive template of upstream vfox.
func coreAssetName(version string) (string, string, error) {
	osName := stdruntime.GOOS
	format := "tar.gz"
	switch osName {
	case "darwin":
		osName = "macos"
	case "windows":
		format = "zip"
	case "linux":
	default:
		return "", "", fmt.Errorf("不支持的平台 %s", osName)
	}

	var arch string
	switch stdruntime.GOARCH {
	case "amd64":
		arch = "x86_64"
	case "386":
		arch = "i386"
	case "arm64":
		arch = "aarch64"
	case "arm":
		arch = "armv7"
	case "loong64":
		arch = "loong64"
	default:
		arch = stdruntime.GOARCH
	}

	name := fmt.Sprintf("vfox_%s_%s_%s.%s", version, osName, arch, format)
	return name, format, nil
}

func (a *App) coreUpdateDir() string {
	return dataPath("core-update")
}

// downloadCoreArtifact fetches a release archive and returns the extracted
// engine binary path inside a temporary directory.
func (a *App) downloadCoreArtifact(version string) (string, error) {
	assetName, format, err := coreAssetName(version)
	if err != nil {
		return "", err
	}
	downloadURL := a.applyGitHubSource(coreDownloadBase + "v" + version + "/" + assetName)

	workDir := a.coreUpdateDir()
	if err := os.MkdirAll(workDir, 0755); err != nil {
		return "", err
	}
	archivePath := filepath.Join(workDir, assetName)
	defer os.Remove(archivePath)

	a.emitEvent("core-update-progress", "downloading", version)
	req, err := http.NewRequest("GET", downloadURL, nil)
	if err != nil {
		return "", err
	}
	req.Header.Set("User-Agent", "QfPlus/1.0")
	resp, err := a.coreHTTPClient(5 * time.Minute).Do(req)
	if err != nil {
		return "", fmt.Errorf("下载失败：%w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("下载失败，GitHub 返回状态码 %d", resp.StatusCode)
	}

	file, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	if _, err := io.Copy(file, resp.Body); err != nil {
		file.Close()
		return "", fmt.Errorf("下载中断：%w", err)
	}
	file.Close()

	extractDir := filepath.Join(workDir, "extract-"+version)
	_ = os.RemoveAll(extractDir)
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return "", err
	}

	var binaryPath string
	switch format {
	case "zip":
		binaryPath, err = extractCoreFromZip(archivePath, extractDir)
	default:
		binaryPath, err = extractCoreFromTarGz(archivePath, extractDir)
	}
	if err != nil {
		return "", err
	}
	if binaryPath == "" {
		return "", fmt.Errorf("安装包中未找到 vfox 可执行文件")
	}
	a.emitEvent("core-update-progress", "installing", version)
	return binaryPath, nil
}

func extractCoreFromZip(archivePath, extractDir string) (string, error) {
	reader, err := zip.OpenReader(archivePath)
	if err != nil {
		return "", fmt.Errorf("无法打开安装包：%w", err)
	}
	defer reader.Close()

	var binaryPath string
	for _, entry := range reader.File {
		name := filepath.ToSlash(entry.Name)
		if entry.FileInfo().IsDir() || filepath.Base(name) != coreExecutableName() {
			continue
		}
		src, err := entry.Open()
		if err != nil {
			return "", err
		}
		target := filepath.Join(extractDir, coreExecutableName())
		dst, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			src.Close()
			return "", err
		}
		if _, err := io.Copy(dst, src); err != nil {
			dst.Close()
			src.Close()
			return "", err
		}
		dst.Close()
		src.Close()
		binaryPath = target
	}
	return binaryPath, nil
}

func extractCoreFromTarGz(archivePath, extractDir string) (string, error) {
	file, err := os.Open(archivePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	gz, err := gzip.NewReader(file)
	if err != nil {
		return "", fmt.Errorf("无法解压缩包：%w", err)
	}
	defer gz.Close()

	var binaryPath string
	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if header.Typeflag != tar.TypeReg {
			continue
		}
		base := filepath.Base(filepath.ToSlash(header.Name))
		if base != coreExecutableName() {
			continue
		}
		target := filepath.Join(extractDir, base)
		dst, err := os.OpenFile(target, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0755)
		if err != nil {
			return "", err
		}
		if _, err := io.Copy(dst, reader); err != nil {
			dst.Close()
			return "", err
		}
		dst.Close()
		binaryPath = target
	}
	return binaryPath, nil
}

// installCoreBinary places a freshly extracted engine binary into the
// user-writable version store and activates it.
func (a *App) installCoreBinary(version, binaryPath string) error {
	targetDir := localCoreVersionDir(version)
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}
	targetPath := filepath.Join(targetDir, coreExecutableName())

	data, err := os.ReadFile(binaryPath)
	if err != nil {
		return err
	}
	tmp := targetPath + ".new"
	if err := os.WriteFile(tmp, data, 0755); err != nil {
		return err
	}
	// A previous run of the same version may still hold a handle on Windows.
	var lastErr error
	for attempt := 0; attempt < 6; attempt++ {
		if err := os.Rename(tmp, targetPath); err == nil {
			lastErr = nil
			break
		} else {
			lastErr = err
			time.Sleep(400 * time.Millisecond)
		}
	}
	if lastErr != nil {
		_ = os.Remove(tmp)
		return fmt.Errorf("写入内核文件失败：%w", lastErr)
	}
	_ = os.Chmod(targetPath, 0755)
	return setMarkerCoreDir(targetDir)
}

// activateBundledCore clears the override so the shipped engine is used again.
func (a *App) activateBundledCore() error {
	return setMarkerCoreDir("")
}

// updateCore downloads and installs a specific engine version.
func (a *App) updateCore(version string) error {
	version = strings.TrimPrefix(strings.TrimSpace(version), "v")
	if version == "" {
		return fmt.Errorf("版本号不能为空")
	}
	if _, _, err := coreAssetName(version); err != nil {
		return err
	}

	binaryPath, err := a.downloadCoreArtifact(version)
	if err != nil {
		a.emitEvent("core-update-progress", "failed", version)
		return err
	}
	defer os.RemoveAll(a.coreUpdateDir())

	if err := a.installCoreBinary(version, binaryPath); err != nil {
		a.emitEvent("core-update-progress", "failed", version)
		return err
	}

	current, err := a.getCurrentCoreVersion()
	if err != nil {
		return fmt.Errorf("内核已替换但无法验证：%w", err)
	}
	if compareCoreVersions(current, version) != 0 {
		return fmt.Errorf("内核版本验证失败，期望 %s，实际 %s", version, current)
	}

	config, _ := a.readAppConfig()
	config.Core.CurrentKnown = current
	_ = a.saveAppConfig(config)

	a.emitEvent("core-update-progress", "done", version)
	a.emitEvent("core-status-changed")
	return nil
}

// autoUpdateCore silently upgrades the engine once per day when enabled.
func (a *App) autoUpdateCore() {
	config, err := a.readAppConfig()
	if err != nil || !config.Core.AutoUpdate {
		return
	}
	if last, err := time.Parse(time.RFC3339, config.Core.LastCheck); err == nil &&
		time.Since(last) < coreCheckInterval {
		return
	}

	releases, err := a.fetchCoreReleases()
	if err != nil {
		a.emitEvent("vfox-log", "[CORE] Auto-update check failed: "+err.Error())
		return
	}
	now := time.Now().Format(time.RFC3339)
	config, _ = a.readAppConfig()
	config.Core.LastCheck = now
	config.Core.LatestKnown = releases[0].Version
	_ = a.saveAppConfig(config)
	a.emitEvent("core-status-changed")

	current, err := a.getCurrentCoreVersion()
	if err != nil {
		return
	}
	if compareCoreVersions(releases[0].Version, current) <= 0 {
		return
	}
	if err := a.updateCore(releases[0].Version); err != nil {
		a.emitEvent("vfox-log", "[CORE] Auto-update failed: "+err.Error())
		a.emitEvent("core-update-error", err.Error())
		return
	}
	a.emitEvent("core-update-done", releases[0].Version)
}
