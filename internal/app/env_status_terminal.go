package app

import (
	"fmt"
	"os"
	"strings"
	"time"
)

// emitDiagnosticLine streams one line into the in-app terminal via the same
// "vfox-log" event channel used by vfox tasks. Empty lines keep a single space
// so the terminal renderer does not drop them.
func (a *App) emitDiagnosticLine(line string) {
	if strings.TrimSpace(line) == "" {
		a.emitEvent("vfox-log", " ")
		return
	}
	a.emitEvent("vfox-log", line)
}

// gatherPluginNames returns the set of SDK names the diagnostic should report
// on: every vfox-managed plugin plus locally registered custom SDKs. This is
// shared with collectEnvironmentStatus so the two never drift.
func (a *App) gatherPluginNames() map[string]string {
	sdks, _ := a.getAllSdks()
	custom := a.getNonVfoxSdksMap()

	pluginNames := map[string]string{} // name -> source
	for _, sdk := range sdks {
		if sdk.Source == "system" {
			continue
		}
		pluginNames[sdk.Name] = sdk.Source
	}
	for name := range custom {
		if _, ok := pluginNames[name]; !ok {
			pluginNames[name] = "custom"
		}
	}
	return pluginNames
}

// streamDiagnosticToTerminal builds the read-only diagnostic report and streams
// it into the in-app terminal as it is computed. Instead of collecting
// everything first and dumping one large block, it prints the header
// immediately and then emits each command's result the moment its probe
// finishes. Because every version probe spawns a real subprocess, the emitted
// events are spread across real time and the terminal fills in progressively
// rather than appearing as one long burst after a silent pause.
//
// The report is read-only: it never mutates PATH, the registry, or any file.
func (a *App) streamDiagnosticToTerminal() error {
	userPath := envUserPath()
	machinePath := envMachinePath()
	shimDir := a.getVfoxHomePath("path-shims")

	// Header + fast metadata first so the terminal reacts instantly.
	a.emitDiagnosticLine("")
	a.emitDiagnosticLine("========================================")
	a.emitDiagnosticLine("  QfPlus 环境诊断报告")
	a.emitDiagnosticLine("========================================")
	a.emitDiagnosticLine("生成时间: " + time.Now().Format("2006-01-02 15:04:05"))
	a.emitDiagnosticLine("VfoxHome: " + a.getVfoxHome())
	a.emitDiagnosticLine("ShimDir:  " + shimDir)
	a.emitDiagnosticLine("")
	a.emitDiagnosticLine("正在扫描 SDK 命令，请稍候…")
	a.emitDiagnosticLine("")

	selections, _ := a.readVfoxGlobalSelections()
	addedPlugins, _ := a.getAddedPlugins()
	probeEnv := livePathEnv(os.Environ(), userPath, machinePath)

	pluginNames := a.gatherPluginNames()

	var broken []string

	// Probe each command and stream its result immediately. The probe itself
	// blocks on a subprocess, which is exactly what paces the streaming.
	for name, source := range pluginNames {
		isCurrent := selections[name] != ""
		for _, alias := range pluginCommandAliases(name) {
			item := a.probeCommand(name, alias, source, isCurrent, probeEnv, userPath, machinePath, addedPlugins)

			a.emitDiagnosticLine("----------------------------------------")
			a.emitDiagnosticLine(fmt.Sprintf("命令: %s  (插件: %s, 来源: %s)", alias, name, source))
			a.emitDiagnosticLine("状态: " + item.State)
			if item.Resolved {
				a.emitDiagnosticLine("路径: " + item.ExePath)
				a.emitDiagnosticLine("目录: " + item.ExeDir)
				if item.Version != "" {
					a.emitDiagnosticLine("版本: " + item.Version)
				}
				a.emitDiagnosticLine(fmt.Sprintf("在用户 PATH: %s   在系统 PATH: %s",
					boolText(item.OnUserPath), boolText(item.OnMachinePath)))
			}
			for _, note := range item.Notes {
				a.emitDiagnosticLine("备注: " + note)
			}
			if !item.Resolved || item.State == "broken" {
				broken = append(broken, alias)
			}
		}
	}

	// PATH scopes + summary last (these reads are instant).
	inPath := a.vfoxShimDirOnPath(userPath, machinePath)
	pathDrift := pathContainsDir(machinePath, shimDir) &&
		!pathContainsDir(os.Getenv("PATH"), shimDir)

	a.emitDiagnosticLine("")
	a.emitDiagnosticLine("当前用户 PATH:")
	a.emitDiagnosticLine(userPath)
	a.emitDiagnosticLine("")
	a.emitDiagnosticLine("当前系统 PATH:")
	a.emitDiagnosticLine(machinePath)
	a.emitDiagnosticLine("")
	a.emitDiagnosticLine("Vfox 在 PATH 中: " + boolText(inPath))
	a.emitDiagnosticLine("PATH 漂移: " + boolText(pathDrift))

	if len(broken) > 0 {
		a.emitDiagnosticLine("")
		a.emitDiagnosticLine("警告:")
		for _, b := range broken {
			a.emitDiagnosticLine("- " + b + " 未就绪")
		}
	}

	a.emitDiagnosticLine("")
	a.emitDiagnosticLine("----------------------------------------")
	a.emitDiagnosticLine("[诊断完成] 如果命令缺失或 PATH 不符，请关闭并重新打开终端。")
	a.emitDiagnosticLine("")

	return nil
}

func boolText(b bool) string {
	if b {
		return "是"
	}
	return "否"
}
