# Release Guide

[中文](RELEASE.zh-CN.md)

QfPlus Windows releases are assembled by hand on a Windows machine: the binaries are produced locally and uploaded to the GitHub Releases page manually. The targets this machine cannot build come from `.github/workflows/build-artifacts.yml`, which only runs when someone presses the button and only leaves downloadable artifacts behind. Nothing publishes on its own, and no tag or push triggers a build.

## Release Scope

| Platform | Artifact | Built by |
| --- | --- | --- |
| Windows amd64, installer | `QfPlus-windows-amd64-installer.exe` | local `wails build -nsis` |
| Windows amd64, portable | `QfPlus-windows-amd64-portable.zip` | local repack |
| Linux amd64, installer | `QfPlus-linux-amd64-installer.deb` | dispatch workflow |
| Linux amd64, portable | `QfPlus-linux-amd64-portable.tar.gz` | dispatch workflow |
| macOS arm64, installer | `QfPlus-darwin-arm64-installer.dmg` | dispatch workflow |
| macOS arm64, portable | `QfPlus-darwin-arm64-portable.zip` | dispatch workflow |

The two Windows artifacts come from the same `wails build` output, so they always carry the same application version and the same core binaries.

## What Gets Bundled

The installer bundles everything found in `core/windows/x86_64`, and the portable archive ships the same directory next to the executable. Nothing is downloaded during the build, so the two files below must already be in place:

| Binary | Current bundled version | Source |
| --- | --- | --- |
| `vfox.exe` | 1.0.11 | https://github.com/version-fox/vfox/releases |
| `mihomo.exe` | v1.18.8 | https://github.com/MetaCubeX/mihomo/releases |

Before bumping either version:

- Confirm the upstream release exists and the Windows archive is downloadable.
- Replace the file under `core/windows/x86_64/`.
- Run the pre-release checks below, including a real install.
- Mention the core version change in the release notes.

## Build the Installer

```bash
wails build -platform windows/amd64 -nsis -clean
```

This leaves `build/bin/QfPlus.exe` and `build/bin/QfPlus-amd64-installer.exe`. Rename the installer to `QfPlus-windows-amd64-installer.exe`.

`-clean` only clears `build/bin`; it never touches `core/`.

## Build the Portable Archive

The portable archive must keep the core directory next to the executable, because `getCoreDir()` resolves `core/<os>/<arch>` relative to the executable first:

```text
QfPlus/
  QfPlus.exe
  README.txt
  core/windows/x86_64/vfox.exe
  core/windows/x86_64/mihomo.exe
```

`README.txt` is tracked at `build/windows/portable/README.txt` - copy it into the archive root instead of writing a new one per release.

Extracting only `QfPlus.exe` leaves neither the SDK engine nor the proxy core available.

Zip the folder itself, not its contents, and name the result `QfPlus-windows-amd64-portable.zip`. Use an archiver that writes forward-slash entry paths - Windows PowerShell's `Compress-Archive` writes backslashes that some extractors turn into one giant filename.

`tar -a -cf QfPlus-windows-amd64-portable.zip QfPlus` works only with the archiver shipped with Windows, `C:\Windows\System32\tar.exe`. Git Bash puts its own `tar` earlier on `PATH`, and GNU tar does not implement `-a`: it writes an uncompressed tar archive with a `.zip` name, which every real extractor rejects. Pass the Windows binary by full path even from Git Bash, then confirm the result is a zip before uploading.

## Optional: 32-bit Installer

`build/windows/installer/project_386.nsi` still produces a 386 installer, but it is not part of a release unless you also place `vfox.exe` and `mihomo.exe` for `windows/x86` under `core/windows/x86` and run:

```bash
wails build -platform windows/386 -clean
makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\QfPlus.exe build\windows\installer\project_386.nsi
```

The define name is shared with the 64-bit script and points at the freshly built executable.

## macOS and Linux Builds

Wails refuses to cross-compile — `wails build -platform linux/amd64` stops with `Crosscompiling to Linux not currently supported.` before it invokes the Go toolchain, and a macOS build additionally needs the Xcode SDK. Those targets have to be compiled on their own system, which is what `.github/workflows/build-artifacts.yml` is for: it runs each target on a matching runner, so no cross-compiling is involved.

To produce them, open **Actions → Build non-Windows artifacts → Run workflow**. The job downloads the pinned `vfox` and `mihomo` builds for each platform, compiles, packages the installer and portable form, and uploads both as workflow artifacts. Download them from the run page and file them into `release/` yourself; the workflow never creates or edits a GitHub Release.

Two limits worth knowing before you advertise those builds:

- The matrix covers `linux/amd64` and `darwin/arm64`. Intel Macs are not built, and adding them means a runner that still ships an x86_64 macOS SDK.
- A workflow run proves the code compiles and packages on that platform; it does not prove the app behaves there. Work through the checklist below on a real install before calling it a release. The macOS output is also unsigned and unnotarized, so Gatekeeper blocks the first launch until the user overrides it.

If you have the hardware, the same result is reachable by hand on that platform:

1. Place the matching core binaries under `core/<os>/<arch>/` (`vfox` and `mihomo`, from their upstream releases).
2. Run `wails build -platform <os>/<arch> -clean` — add `-tags webkit2_41` on Ubuntu 24.04 — then the platform packaging step (`nfpm` reads `nfpm.yaml` generated from the template on Linux, `hdiutil` builds the dmg on macOS).
3. Work through the checklist below on a real install — SDK install/use/unuse, custom SDKs, plugin add/remove, the built-in proxy, and PATH integration all behave differently per platform.

Name the artifacts the same way as the Windows ones, as in the scope table above, and attach them to the same release.

## Pre-Release Checklist

Run these checks before packaging:

```bash
go test ./...
npm --prefix frontend run build
git diff --check
```

Manual checks:

- App launches without a console window for installer/uninstaller helper operations.
- SDK list, detail, version search, install, uninstall, use, and unuse work.
- SDKs with no published releases show a retryable "no release version" state instead of locking the action.
- Custom SDK add, detect, use, and remove work.
- Plugin add and remove keep or delete custom SDK paths according to user choice.
- Download path migration confirmation is centered like other floating windows.
- Proxy settings import a subscription, and the proxy only affects processes QfPlus starts.
- Uninstall removes app data, shims, and PATH/override residue expected by the installer cleanup design.
- English and Chinese UI text both render correctly.

## Uploading

Create the release on GitHub, then attach both artifacts. Tagging is optional and only marks the source revision - no pipeline reacts to it.

## Installer Assets

Windows installer files live under `build/windows/installer/`:

| File | Purpose |
| --- | --- |
| `project.nsi` | Main Windows installer script. |
| `project_386.nsi` | 32-bit installer script. |
| `cleanup_qfplus.ps1` | Cleanup helper used by uninstall logic. |
| `wails_tools.nsh` | Installer helper macros and Wails integration. |

The installer should not leave visible PowerShell windows during uninstall. If cleanup behavior changes, verify it with a real installed build instead of only running from source.

## Release Notes

Release notes should include:

- User-visible fixes and features.
- SDK/plugin/PATH/migration behavior changes.
- vfox and mihomo core versions.
- Known platform limitations.
- Upgrade or uninstall notes when cleanup behavior changes.

Avoid listing internal refactors unless they affect maintainability, testing, or future contributor workflow.
