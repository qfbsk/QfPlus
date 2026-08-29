# QfPlus

[中文文档](READMEcn.md)

<p align="center">
  <img src="build/appicon.png" alt="QfPlus" width="128" />
</p>

<p align="center">
  <strong>vfox GUI Manager · Wails + Vue 3</strong>
</p>

<p align="center">
  <a href="https://github.com/qfbsk/QfPlus/">https://github.com/qfbsk/QfPlus/</a>
</p>

<p align="center">
  Manage vfox SDKs, plugins, system SDKs, and PATH integration from a desktop GUI.
</p>

---

## Screenshot

![QfPlus screenshot](image.png)

## Downloads

Download published builds from the Releases page of this repository.

QfPlus is built with Wails and ships Windows amd64 builds:

| Platform | Artifact |
| --- | --- |
| Windows amd64, installer | `QfPlus-windows-amd64-installer.exe` |
| Windows amd64, portable | `QfPlus-windows-amd64-portable.zip` |

The portable archive needs no installation and no administrator rights: extract the whole `QfPlus` folder and run `QfPlus.exe`. Settings and downloaded SDKs still live in `%APPDATA%\QfPlus`.

QfPlus is written in Go and Vue and runs cross-platform, and the repository keeps the packaging assets for every target: NSIS scripts for Windows, `build/darwin` plists for macOS, and `nfpm.yaml.tmpl` for Linux packages. Windows amd64 is built and verified here; Wails refuses to cross-compile, so macOS and Linux are compiled on their own system — either by hand on that hardware or through the manual-trigger `Build non-Windows artifacts` workflow, which downloads the matching `core/<os>/<arch>` binaries and hands the packages back as build artifacts.

## Features

- SDK management: view, install, uninstall, switch, and unuse vfox-managed SDK versions.
- Custom SDKs: add system-installed SDKs, detect versions from executables, and switch them through the same managed entrypoints.
- Plugin marketplace: browse available vfox plugins, add plugins, remove plugins, and keep custom SDK paths when needed.
- Release-aware version search: versions with no published releases are shown as available to retry instead of locking the SDK.
- Built-in proxy: import your own subscription in Settings, pick a node group or node, and measure latency. No server, node or subscription is bundled or pre-configured, and the proxy only applies to processes QfPlus starts itself - system proxy settings are never touched.
- PATH integration: add or remove vfox from user PATH, and manage SDK command overrides where supported.
- Windows compatibility: handle App Execution Alias conflicts, junctions, hidden PowerShell helper windows, and installer cleanup.
- Download directory migration: review a plan of what will be copied versus listed-only, then move the vfox home/download location with progress feedback.
- Bilingual UI and docs: Chinese and English resources are kept in parallel.

## Documentation

| Topic | English | Chinese |
| --- | --- | --- |
| Documentation index | [docs/README.md](docs/README.md) | [docs/README.zh-CN.md](docs/README.zh-CN.md) |
| Architecture | [docs/ARCHITECTURE.md](docs/ARCHITECTURE.md) | [docs/ARCHITECTURE.zh-CN.md](docs/ARCHITECTURE.zh-CN.md) |
| Development | [docs/DEVELOPMENT.md](docs/DEVELOPMENT.md) | [docs/DEVELOPMENT.zh-CN.md](docs/DEVELOPMENT.zh-CN.md) |
| Release | [docs/RELEASE.md](docs/RELEASE.md) | [docs/RELEASE.zh-CN.md](docs/RELEASE.zh-CN.md) |
| Code style | [docs/CODE_STYLE.en.md](docs/CODE_STYLE.en.md) | [docs/CODE_STYLE.md](docs/CODE_STYLE.md) |
| Frontend | [frontend/README.md](frontend/README.md) | [frontend/README.zh-CN.md](frontend/README.zh-CN.md) |
| Build assets | [build/README.md](build/README.md) | [build/README.zh-CN.md](build/README.zh-CN.md) |

## Architecture Overview

```text
+------------------------------+
| Frontend: Vue 3 + TypeScript |
| components / composables     |
| services / i18n / styles     |
+--------------+---------------+
               |
| Wails generated bindings
               v
+--------------+---------------+
| Backend: internal/app        |
| facade / sdk / plugin / environment |
| config / path / platform     |
+--------------+---------------+
               |
               | command execution
               v
+--------------+---------------+
| vfox CLI core                |
| bundled in release builds    |
+------------------------------+
```

The repository root now only keeps the Wails entrypoint in `main.go`. Backend application logic lives under `internal/app`, DTOs live under `internal/model`, and pure vfox output parsing lives under `internal/parser`. Public Wails methods stay in `internal/app/app_facade_*.go`; domain workflows live in focused files such as `sdk_use.go`, `plugin_remove.go`, `migration_run.go`, and `environment_import.go`.

The frontend follows a component/composable/service direction:

```text
components -> composables -> services -> frontend/wailsjs
```

Components handle rendering and user events, composables hold UI state and workflows, and services are the only layer that imports generated Wails bindings.

## Runtime Core

QfPlus drives the vfox engine and the mihomo proxy core as external binaries under `core/`. Nothing downloads them for you - place them there by hand:

```text
core/
  windows/
    x86_64/vfox.exe          # from https://github.com/version-fox/vfox/releases
    x86_64/mihomo.exe        # from https://github.com/MetaCubeX/mihomo/releases
    x86/vfox.exe
    x86/mihomo.exe
```

Both binaries are looked up in the same directory at runtime, so a missing `mihomo.exe` disables the built-in proxy while the rest of the app keeps working. `core/` is ignored by Git so local binaries are not committed.

## Development

Requirements:

| Tool | Version |
| --- | --- |
| Go | 1.23+ |
| Node.js | 22+ |
| Wails CLI | v2 |
| NSIS | 3.x, Windows installer only |

Install frontend dependencies:

```bash
npm --prefix frontend install
```

Run the desktop app:

```bash
wails dev
```

Verify the project:

```bash
go test ./...
npm --prefix frontend run build
```

Build locally:

```bash
wails build -clean
```

Build a Windows amd64 installer:

```bash
wails build -platform windows/amd64 -nsis -clean
```

## Project Structure

```text
QfPlus/
  main.go                        Wails startup and app binding entrypoint
  internal/
    app/                         Wails facade and stateful backend workflows
    model/                       DTOs shared with generated frontend bindings
    parser/                      pure parsers for vfox command output
    pathutil/                    shared path comparison and PATH cleanup helpers
    storage/                     shared JSON file persistence helpers
  frontend/                      Vue 3 frontend
  build/                         icons, manifests, installer scripts
  docs/                          bilingual project documentation
  tools/
    genicon/                     renders build/appicon.png, icon.ico and the UI mark
```

Regenerate every icon artifact after changing the mark geometry:

```bash
go run ./tools/genicon
```

## Release Process

Releases are assembled by hand on a Windows machine with the toolchain above. Build the installer:

```bash
wails build -platform windows/amd64 -nsis -clean
```

That leaves `build/bin/QfPlus.exe` and `build/bin/QfPlus-amd64-installer.exe`; rename the latter to `QfPlus-windows-amd64-installer.exe`.

Then assemble the portable folder so the core sits next to the executable, and zip the folder itself:

```text
QfPlus/
  QfPlus.exe
  README.txt
  core/windows/x86_64/vfox.exe
  core/windows/x86_64/mihomo.exe
```

Upload both files to the GitHub Releases page. `docs/RELEASE.md` keeps the full checklist and the pre-release verification steps.

## Acknowledgments

Thanks to the [vfox](https://github.com/version-fox/vfox) project for the cross-platform version management engine that powers QfPlus.

The GUI was designed by QfPlus, using vFoxG (by HuajiFruit) as a reference for page structure and interaction patterns.

Thanks to 千问3.8-Flash (Qwen), which was of great help throughout this work - reading the codebase, drafting and verifying changes, and writing the documentation.

Author: 清枫不识客
