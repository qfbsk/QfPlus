export const enEnvironment = {
  // Page
  'environment.title': 'Environment',
  'environment.updated_at': 'Updated {time}',

  // Info box
  'environment.info.title': 'About this page',
  'environment.info.desc': 'Read-only view of how each SDK command resolves on PATH, its version, and ownership. Click "Open diagnostic console" to inspect the full PATH in a separate window.',
  'environment.info.note': 'No check modifies PATH or any file.',

  // Status card
  'environment.status.title': 'Status',
  'environment.status.refresh': 'Refresh status',
  'environment.status.detect': 'Open diagnostic console',
  'environment.status.loading': 'Loading status...',
  'environment.status.empty': 'No status yet. Click refresh to scan.',
  'environment.status.vfox_in_path': 'vfox shim dir on PATH',
  'environment.status.path_drift': 'PATH drift detected',
  'environment.status.vfox_home': 'Vfox home',
  'environment.status.shim_dir': 'Shim directory',

  // Item states
  'environment.state.ok': 'OK',
  'environment.state.managed': 'Managed',
  'environment.state.broken': 'Broken',
  'environment.state.missing': 'Missing',
  'environment.state.unmanaged': 'Unmanaged',
  'environment.state.unknown': 'Unknown',
  'environment.state.tooltip.ok': 'Resolved on PATH and owned by vfox',
  'environment.state.tooltip.managed': 'Resolved on PATH through QfPlus shim',
  'environment.state.tooltip.broken': 'Found on PATH but version mismatch or not owned',
  'environment.state.tooltip.missing': 'Not found on PATH',
  'environment.state.tooltip.unmanaged': 'Resolved on PATH but not owned by vfox',
  'environment.state.tooltip.unknown': 'Not scanned yet',

  // Status item detail (circle + exclamation)
  'environment.detail.toggle': 'View details',
  'environment.detail.reason': 'Reason',
  'environment.detail.executable': 'Executable',
  'environment.detail.version': 'Version',
  'environment.detail.managed_by': 'Managed by',
  'environment.detail.source': 'Source',
  'environment.detail.on_path': 'On PATH',
  'environment.detail.path_user': 'User PATH',
  'environment.detail.path_machine': 'Machine PATH',
  'environment.detail.path_both': 'User + Machine PATH',
  'environment.detail.path_none': 'Not on PATH',
  'environment.detail.unknown': 'Unknown',

  // Shared download queue (plugin installs)
  'download_queue.title': 'Downloads',
  'download_queue.plugin_title': 'Install plugin',
  'download_queue.queued': 'Queued',
  'download_queue.running': 'Downloading',
  'download_queue.done': 'Done',
  'download_queue.error': 'Failed',
  'download_queue.files': 'File {done}/{total}',
  'download_queue.speed': '{speed}',
  'download_queue.clear': 'Clear finished',
  'download_queue.dismiss': 'Dismiss',
};
