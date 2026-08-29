//go:build windows

package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// A shim forwards to one executable, so the probe matrix below decides both
// which lines the script tries at run time and which aliases deserve a script
// at all. windowsShimStore checks it before writing, otherwise a plugin whose
// name is not a command ("golang" shipping go.exe) would leave a shim that can
// only ever print its own failure message.
var (
	windowsShimProbeDirs = []string{"", "bin", "Scripts", "sbin"}
	windowsShimProbeExts = []string{".exe", ".cmd", ".bat"}
)

func windowsShimRelPath(dir string, alias string, ext string) string {
	if dir == "" {
		return alias + ext
	}
	return dir + `\` + alias + ext
}

func windowsShimTargetExists(sdkRoot string, alias string) bool {
	if strings.TrimSpace(sdkRoot) == "" {
		return false
	}
	for _, ext := range windowsShimProbeExts {
		for _, dir := range windowsShimProbeDirs {
			candidate := filepath.Join(sdkRoot, dir, alias+ext)
			if info, err := os.Stat(candidate); err == nil && !info.IsDir() {
				return true
			}
		}
	}
	return false
}

func windowsShimProbes(alias string) string {
	var probes strings.Builder
	for _, ext := range windowsShimProbeExts {
		invoke := ""
		if ext != ".exe" {
			invoke = "call "
		}
		for _, dir := range windowsShimProbeDirs {
			rel := windowsShimRelPath(dir, alias, ext)
			fmt.Fprintf(&probes, "if exist \"%%SDK_ROOT%%\\%s\" (%s\"%%SDK_ROOT%%\\%s\" %%* & exit /b)\n", rel, invoke, rel)
		}
	}
	return probes.String()
}

func windowsShimScript(pluginName string, alias string, sdkPath string) string {
	alias = windowsSafeShimName(alias)
	return fmt.Sprintf(`@echo off
setlocal
set "SDK_ROOT=%[1]s"
set "ALIAS_NAME=%[2]s"
%[4]sfor /f "delims=" %%%%I in ('where "%%ALIAS_NAME%%" 2^>nul') do (
  if /I not "%%%%~fI"=="%%~f0" (
    if /I not "%%%%~dpI"=="%%~dp0" (
      if /I "%%%%~xI"==".cmd" (call "%%%%~fI" %%* & exit /b)
      if /I "%%%%~xI"==".bat" (call "%%%%~fI" %%* & exit /b)
      "%%%%~fI" %%*
      exit /b
    )
  )
)
echo QfPlus: %[2]s for %[3]s is not available under %%SDK_ROOT%%, and no fallback %[2]s was found on PATH. 1>&2
exit /b 9009
`, sdkPath, alias, pluginName, windowsShimProbes(alias))
}
