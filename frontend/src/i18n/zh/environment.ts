import type { enEnvironment } from '../en/environment';

export const zhEnvironment: Record<keyof typeof enEnvironment, string> = {
  // Page
  'environment.title': '环境',
  'environment.updated_at': '更新于 {time}',

  // Info box
  'environment.info.title': '关于本页',
  'environment.info.desc': '这里只读展示每个 SDK 命令在 PATH 上的解析结果、版本与归属。点“打开诊断控制台”可在独立窗口查看完整 PATH。',
  'environment.info.note': '所有检测都不会修改 PATH 或任何文件。',

  // Status card
  'environment.status.title': '状态',
  'environment.status.refresh': '刷新状态',
  'environment.status.detect': '打开诊断控制台',
  'environment.status.loading': '正在加载状态...',
  'environment.status.empty': '暂无状态，点击刷新开始扫描。',
  'environment.status.vfox_in_path': 'vfox shim 目录已在 PATH',
  'environment.status.path_drift': '检测到 PATH 漂移',
  'environment.status.vfox_home': 'Vfox 主目录',
  'environment.status.shim_dir': 'Shim 目录',

  // Item states
  'environment.state.ok': '正常',
  'environment.state.managed': '已接管',
  'environment.state.broken': '异常',
  'environment.state.missing': '缺失',
  'environment.state.unmanaged': '未接管',
  'environment.state.unknown': '未检测',
  'environment.state.tooltip.ok': 'PATH 可解析且由 vfox 管理',
  'environment.state.tooltip.managed': '通过 QfPlus shim 在 PATH 上解析',
  'environment.state.tooltip.broken': 'PATH 上有命中但版本不匹配或未被接管',
  'environment.state.tooltip.missing': 'PATH 上未找到',
  'environment.state.tooltip.unmanaged': 'PATH 上有命中，但不是由 vfox 接管',
  'environment.state.tooltip.unknown': '尚未扫描',

  // Status item detail (circle + exclamation)
  'environment.detail.toggle': '查看详情',
  'environment.detail.reason': '原因',
  'environment.detail.executable': '可执行文件',
  'environment.detail.version': '版本',
  'environment.detail.managed_by': '接管方',
  'environment.detail.source': '来源',
  'environment.detail.on_path': '是否在 PATH',
  'environment.detail.path_user': '用户 PATH',
  'environment.detail.path_machine': '系统 PATH',
  'environment.detail.path_both': '用户 + 系统 PATH',
  'environment.detail.path_none': '不在 PATH',
  'environment.detail.unknown': '未知',

  // 共享下载队列（插件安装）
  'download_queue.title': '下载队列',
  'download_queue.plugin_title': '安装插件',
  'download_queue.queued': '排队中',
  'download_queue.running': '下载中',
  'download_queue.done': '已完成',
  'download_queue.error': '失败',
  'download_queue.files': '文件 {done}/{total}',
  'download_queue.speed': '{speed}',
  'download_queue.clear': '清空已完成',
  'download_queue.dismiss': '关闭',
};
