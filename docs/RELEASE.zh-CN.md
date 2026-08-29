# 发布指南

[English](RELEASE.md)

QfPlus 的 Windows 发布是在 Windows 机器上手工完成的：产物本地构建，再手动上传到 GitHub 的 Releases 页面。这台机器出不了的平台由 `.github/workflows/build-artifacts.yml` 负责，它只在有人点了按钮时才跑，跑完只留下一堆可下载的构建产物。任何东西都不会自动发布，也没有 push 或 tag 会触发构建。

## 发布范围

| 平台 | 产物 | 构建方式 |
| --- | --- | --- |
| Windows amd64 安装包 | `QfPlus-windows-amd64-installer.exe` | 本地 `wails build -nsis` |
| Windows amd64 便携版 | `QfPlus-windows-amd64-portable.zip` | 本地重新打包 |
| Linux amd64 安装包 | `QfPlus-linux-amd64-installer.deb` | 手动触发 workflow |
| Linux amd64 便携版 | `QfPlus-linux-amd64-portable.tar.gz` | 手动触发 workflow |
| macOS arm64 安装包 | `QfPlus-darwin-arm64-installer.dmg` | 手动触发 workflow |
| macOS arm64 便携版 | `QfPlus-darwin-arm64-portable.zip` | 手动触发 workflow |

两个 Windows 产物出自同一次 `wails build`，因此应用版本和内置内核版本始终一致。

## 内置了什么

安装包会把 `core/windows/x86_64` 目录里的全部内容打进去，便携版则把同一个目录放在可执行文件旁边。构建过程不会下载任何东西，所以下面两个文件必须提前放好：

| 二进制 | 当前内置版本 | 来源 |
| --- | --- | --- |
| `vfox.exe` | 1.0.11 | https://github.com/version-fox/vfox/releases |
| `mihomo.exe` | v1.18.8 | https://github.com/MetaCubeX/mihomo/releases |

升级任一版本前需要确认：

- 上游 release 已存在，Windows archive 可以下载。
- 替换 `core/windows/x86_64/` 下对应文件。
- 完成下面的发布前检查，包括真实安装验证。
- Release notes 里说明 core 版本变化。

## 构建安装包

```bash
wails build -platform windows/amd64 -nsis -clean
```

产物是 `build/bin/QfPlus.exe` 和 `build/bin/QfPlus-amd64-installer.exe`，把后者改名为 `QfPlus-windows-amd64-installer.exe`。

`-clean` 只清空 `build/bin`，不会动 `core/`。

## 打包便携版

便携版必须把 core 目录放在可执行文件旁边，因为 `getCoreDir()` 会优先按可执行文件所在目录解析 `core/<os>/<arch>`：

```text
QfPlus/
  QfPlus.exe
  README.txt
  core/windows/x86_64/vfox.exe
  core/windows/x86_64/mihomo.exe
```

其中的 `README.txt` 已在仓库中维护，路径是 `build/windows/portable/README.txt`，打包时直接复制到压缩包根目录，不要每次重新写一份。

只解压 `QfPlus.exe` 会导致版本管理引擎和代理内核都不可用。

压缩时要打包这个文件夹本身，结果命名为 `QfPlus-windows-amd64-portable.zip`。请使用写入正斜杠路径的压缩工具，Windows PowerShell 的 `Compress-Archive` 写入的是反斜杠，部分解压工具会把它当成一个超长文件名。

`tar -a -cf QfPlus-windows-amd64-portable.zip QfPlus` 这条命令只有 Windows 自带的 `C:\Windows\System32\tar.exe` 支持。Git Bash 会把它自己的 `tar` 排在 `PATH` 前面，而 GNU tar 并不实现 `-a`：它产出的是一个不带压缩的 tar 包，只是文件名叫 `.zip`，任何正规解压工具都会拒绝打开。即使在 Git Bash 里执行，也要写 Windows 版 tar 的完整路径，并且上传前先确认产物确实是 zip。

## 可选：32 位安装包

`build/windows/installer/project_386.nsi` 仍然可以产出 386 安装包，但除非你另外把 `windows/x86` 版本的 `vfox.exe` 和 `mihomo.exe` 放进 `core/windows/x86`，否则它不属于发布内容：

```bash
wails build -platform windows/386 -clean
makensis -DARG_WAILS_AMD64_BINARY=..\..\bin\QfPlus.exe build\windows\installer\project_386.nsi
```

这个宏名与 64 位脚本共用，指向刚构建出来的可执行文件。

## macOS 和 Linux 产物

Wails 拒绝跨平台构建——`wails build -platform linux/amd64` 会在调用 Go 工具链之前直接停在 `Crosscompiling to Linux not currently supported.`，macOS 构建还额外需要 Xcode SDK。这两个平台必须在各自系统上编译，`.github/workflows/build-artifacts.yml` 就是为此而生：每个目标跑在匹配的 runner 上，全程不涉及交叉编译。

要出这些产物，进 **Actions → Build non-Windows artifacts → Run workflow**。任务会下载各平台固定版本的 `vfox` 和 `mihomo`，编译、打成安装包和便携版两种形态，再把结果作为构建产物上传。你从运行页面下载后自己放进 `release/`；这个 workflow 不会创建也不会改动 GitHub Release。

两点需要提前说清楚：

- 矩阵只覆盖 `linux/amd64` 和 `darwin/arm64`。不构建 Intel Mac，要加的话得找仍然提供 x86_64 macOS SDK 的 runner。
- workflow 跑通只能证明代码在该平台能编译、能打包，不能证明程序在该平台跑得对。要当成正式发布，还得按下面的检查清单在真实安装环境走一遍。另外 macOS 产物既没签名也没公证，用户首次打开会被 Gatekeeper 拦一次。

如果你手上有对应机器，同样的结果也能手工做出来：

1. 把对应平台的 core 二进制（`vfox` 和 `mihomo`，取自各自上游 release）放进 `core/<os>/<arch>/`。
2. 执行 `wails build -platform <os>/<arch> -clean`（Ubuntu 24.04 上要加 `-tags webkit2_41`），再做该平台的打包步骤（Linux 上 `nfpm` 读取由模板生成的 `nfpm.yaml`，macOS 上用 `hdiutil` 出 dmg）。
3. 按下面的检查清单在真实安装环境走一遍——SDK 安装/应用/取消使用、自定义 SDK、插件增删、内置代理和 PATH 接管在各平台行为都不同。

产物命名与上表保持一致，附加到同一个 Release 上。

## 发布前检查清单

打包前运行：

```bash
go test ./...
npm --prefix frontend run build
git diff --check
```

手动检查：

- 应用启动正常，安装/卸载辅助操作不会弹出可见控制台窗口。
- SDK 列表、详情、版本搜索、安装、卸载、使用、取消使用正常。
- 没有已发布版本的 SDK 显示可重试的“无发布版本”状态，不锁死操作按钮。
- 自定义 SDK 添加、检测、使用、移除正常。
- 插件添加和移除按用户选择保留或删除自定义 SDK 路径。
- 下载目录迁移确认窗口像其它悬浮窗口一样居中显示。
- 代理设置能导入订阅，并且代理只影响 QfPlus 自己启动的进程。
- 卸载会清理安装器设计范围内的 app data、shim、PATH/override 残留。
- 中文和英文 UI 文案显示正常。

## 上传

先在 GitHub 上创建 Release，再把两个产物附加上去。打 tag 是可选项，只用于标记源码版本——没有任何流水线会响应它。

## 安装器资源

Windows 安装器文件位于 `build/windows/installer/`：

| 文件 | 用途 |
| --- | --- |
| `project.nsi` | 主 Windows 安装脚本。 |
| `project_386.nsi` | 32 位安装脚本。 |
| `cleanup_qfplus.ps1` | 卸载逻辑使用的清理 helper。 |
| `wails_tools.nsh` | 安装器 helper 宏和 Wails 集成。 |

卸载时不应出现可见 PowerShell 窗口。清理行为变化时，需要用真实安装包验证，不能只从源码运行。

## Release Notes

Release notes 应包含：

- 用户可见的修复和功能。
- SDK、插件、PATH、迁移行为变化。
- vfox 和 mihomo core 版本。
- 已知平台限制。
- 清理行为变化时的升级或卸载说明。

除非内部重构影响可维护性、测试或后续协作流程，否则不需要列入发布说明。
