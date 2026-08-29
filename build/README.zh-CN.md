# 构建资源

[English](README.md)

本目录保存 QfPlus 的 Wails 构建资源、平台 manifest 和 Windows 安装脚本。

## 目录结构

```text
build/
  appicon.png                   Wails 使用的源图标
  screenshot.png                可选截图资源
  darwin/
    Info.plist                  macOS 生产 plist
    Info.dev.plist              macOS 开发 plist
  windows/
    icon.ico                    Windows 应用图标
    info.json                   Windows 版本/资源元数据
    wails.exe.manifest          Windows 应用 manifest
    installer/
      project.nsi               主 Windows 安装脚本
      project_386.nsi           32 位 Windows 安装脚本
      cleanup_qfplus.ps1         卸载清理 helper
      wails_tools.nsh           Wails/NSIS helper 宏
    portable/
      README.txt                便携版压缩包根目录的说明文件
```

## Windows 安装包

发布版本提供 Windows 安装包，全部在 Windows 机器上手工构建。amd64 安装包由 Wails 开启 NSIS 构建：

```bash
wails build -platform windows/amd64 -nsis -clean
```

386 安装包是可选产物，通过 `project_386.nsi` 直接调用 NSIS 构建，命令见 `docs/RELEASE.zh-CN.md`。

安装器清理行为很重要，因为 QfPlus 会创建托管 SDK 入口、Windows shim 和 PATH override 元数据。修改安装器文件时，需要验证真实安装和卸载流程。

## macOS 和 Linux

后端可以编译到 `linux/amd64`、`darwin/arm64` 等 Go 目标，仓库也保留了 macOS 的 plist 和根目录的 Linux 打包模板 `nfpm.yaml.tmpl`。但这些仍然必须在目标系统上构建：Wails 拒绝跨平台构建，在 Windows 上执行 `wails build -platform linux/amd64` 会停在 `Crosscompiling to Linux not currently supported.`。每个平台的产物都当成一次该平台本机构建，并准备好各自的 `core/<os>/<arch>` 二进制——Ubuntu 24.04 上要加 `-tags webkit2_41`，因为该版本已经不再有 WebKitGTK 4.0。手上没有对应机器时，`.github/workflows/build-artifacts.yml` 就是在匹配的 runner 上做这件事，详见 `docs/RELEASE.zh-CN.md`。

## 隐藏辅助窗口

卸载和清理 helper 不应显示可见 PowerShell 窗口。如果安装脚本需要启动 PowerShell，使用现有隐藏窗口方式，并用真实安装包确认。

## Core 二进制

安装器会把 `core/windows/x86_64/` 里的现有文件一并打包，所以打包前需要手动放好 `vfox.exe` 和 `mihomo.exe`，没有任何自动下载。`core/` 不在 `build/` 目录下，并且已被 Git 忽略。

不要把下载的 core 二进制提交进仓库。

## 安全编辑规则

- 除非产品行为需要，不要随意修改 Wails 生成的默认文件。
- 安装器行为变化需要同步更新 `docs/RELEASE.md` 和 `docs/RELEASE.zh-CN.md`。
- 修改 `build/windows/installer/` 时，测试安装和卸载。
- 避免把图标或 manifest 变化和无关代码行为混在一起。
