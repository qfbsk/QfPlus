import type { enShell } from '../en/shell';

export const zhShell: Record<keyof typeof enShell, string> = {
  // Navigation
  'nav.installed': '已安装',
  'nav.display': '已安装',
  'nav.environment': '环境',
  'nav.market': '插件市场',
  'nav.settings': '设置',
  'nav.about': '关于',
  'nav.diagnostic': '诊断',

  // Home proxy shortcut
  'home.proxy.on': '代理已开启',
  'home.proxy.off': '代理已关闭',
  'home.proxy.toggle': '一键开关代理',
  'home.proxy.test_delay': '点击测速',
  'home.proxy.testing': '测速中',

  // Terminal
  'terminal.title': '终端',
  'terminal.clear': '清空终端',
  'terminal.expand': '展开终端',
  'terminal.collapse': '收起终端',
  'terminal.empty': '暂无终端输出。',
  'terminal.aria': '终端输出',

  // Common
  'common.error': '错误',
  'common.success': '成功',
  'common.notification': '通知',
  'common.confirm': '确认',
  'common.cancel': '取消',
  'common.unknown': '未知',
  'common.expand': '展开',
  'common.collapse': '收起',
  'common.loading': '加载中...',
  'common.copy_path': '复制路径',

  // Toasts and tasks
  'toast.starting': '正在开始...',
  'toast.completed': '已成功完成！',
  'toast.task_failed': '任务失败。',
  'toast.please_wait': '请稍后',
  'toast.phase.download': '下载中',
  'toast.phase.install': '安装中',
  'toast.installing_after_download': '下载完成，正在安装...',
  'task.plugin.add': '正在添加插件：{name}',
  'task.plugin.remove': '正在移除插件：{name}',
  'task.version.install': '正在安装 {name}@{version}',
  'task.version.uninstall': '正在卸载 {name}@{version}',
  'task.version.switch': '正在将 {name} 切换到 {version}',
  'task.version.unset': '正在取消使用 {name}',
  'task.custom.use': '正在使用 {name}（{path}）',

  // Platform sentences shared by Settings and the SDK pages
  'platform.admin.required': ' 需要管理员授权。',
  'platform.restart.windows': '修改 PATH 后请打开新的终端窗口。',
  'platform.restart.darwin': '请打开新的终端窗口，或运行 `source ~/.zprofile`。',
  'platform.restart.linux': '请打开新的终端窗口，或运行 `source ~/.profile`。',
  'platform.restart.default': '修改 PATH 后请打开新的终端窗口。',
};
