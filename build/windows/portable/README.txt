QfPlus 便携版 / Portable build
================================

用法 / Usage
  1. 把整个 QfPlus 文件夹解压到任意位置（不要只解压 QfPlus.exe）。
     Extract the whole QfPlus folder somewhere, not just QfPlus.exe.
  2. 双击 QfPlus.exe 即可运行，无需安装、无需管理员权限。
     Run QfPlus.exe directly - no install, no admin rights.

目录结构 / Layout
  QfPlus.exe
  core\windows\x86_64\vfox.exe      版本管理引擎 / SDK version engine
  core\windows\x86_64\mihomo.exe    内置代理内核 / bundled proxy core

说明 / Notes
  配置与已安装的 SDK 写入 %APPDATA%\QfPlus，不会改动系统代理。
  Settings and installed SDKs live in %APPDATA%\QfPlus, and the built-in proxy
  only affects processes started by QfPlus itself.

  只有在应用内启用「SDK PATH」时才会写系统 PATH，那一步会请求管理员授权。
  Enabling "SDK PATH" in the app is the only thing that writes the system PATH,
  and that step asks for elevation.

  代理节点需要你自己导入订阅地址：程序内 设置 -> 网络代理。
  Bring your own subscription: Settings -> Proxy. Nothing is pre-configured.
