# Build Assets

[中文](README.zh-CN.md)

This directory contains Wails build assets, platform manifests, and Windows installer scripts for QfPlus.

## Directory Layout

```text
build/
  appicon.png                   Source icon used by Wails
  screenshot.png                Optional screenshot asset
  darwin/
    Info.plist                  macOS production plist
    Info.dev.plist              macOS development plist
  windows/
    icon.ico                    Windows application icon
    info.json                   Windows version/resource metadata
    wails.exe.manifest          Windows app manifest
    installer/
      project.nsi               Main Windows installer script
      project_386.nsi           32-bit Windows installer script
      cleanup_qfplus.ps1         Uninstall cleanup helper
      wails_tools.nsh           Wails/NSIS helper macros
    portable/
      README.txt                Readme shipped in the portable archive root
```

## Windows Installer

Releases ship Windows installers, built by hand on a Windows machine. The amd64 installer is built by Wails with NSIS enabled:

```bash
wails build -platform windows/amd64 -nsis -clean
```

The 386 installer is optional and is produced by calling NSIS directly with `project_386.nsi`; see `docs/RELEASE.md` for the command.

Installer cleanup behavior matters because QfPlus creates managed SDK entrypoints, Windows shims, and PATH override metadata. When changing installer files, verify a real install and uninstall flow.

## macOS and Linux

The backend compiles for `linux/amd64`, `darwin/arm64` and the other Go targets, and this directory keeps the macOS plists plus the `nfpm.yaml.tmpl` Linux package template at the repository root. Building them still has to happen on the target system: Wails refuses to cross-compile, so `wails build -platform linux/amd64` on Windows stops with `Crosscompiling to Linux not currently supported.` Treat each platform build as a build on that platform, with its own `core/<os>/<arch>` binaries in place — on Ubuntu 24.04 add `-tags webkit2_41`, since that release dropped WebKitGTK 4.0. When you do not have the hardware, `.github/workflows/build-artifacts.yml` does exactly that on matching runners; see `docs/RELEASE.md`.

## Hidden Helper Windows

Uninstall and cleanup helpers should not show a visible PowerShell window. If installer scripts launch PowerShell, use the existing hidden-window approach and confirm it with an installed build.

## Core Binaries

The installer bundles whatever sits in `core/windows/x86_64/`, so place `vfox.exe` and `mihomo.exe` there before packaging. Nothing downloads them for you. The `core/` directory is intentionally outside this `build/` directory and is ignored by Git.

Do not commit downloaded core binaries.

## Safe Editing Rules

- Keep Wails-generated default files unless a product behavior requires changes.
- Document installer behavior changes in `docs/RELEASE.md` and `docs/RELEASE.zh-CN.md`.
- Test both install and uninstall when touching `build/windows/installer/`.
- Avoid changing app icons or manifests together with unrelated code behavior.
