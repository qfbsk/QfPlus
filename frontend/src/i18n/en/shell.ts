// Shell strings: navigation chrome, home proxy shortcut, terminal dock,
// generic labels, task toasts and the platform sentences shared by several pages.
export const enShell = {
  // Navigation
  'nav.installed': 'Installed',
  'nav.display': 'Installed',
  'nav.environment': 'Environment',
  'nav.market': 'Plugin Market',
  'nav.settings': 'Settings',
  'nav.about': 'About',
  'nav.diagnostic': 'Diagnostic',

  // Home proxy shortcut
  'home.proxy.on': 'Proxy on',
  'home.proxy.off': 'Proxy off',
  'home.proxy.toggle': 'Toggle proxy',
  'home.proxy.test_delay': 'Click to test latency',
  'home.proxy.testing': 'testing',

  // Terminal
  'terminal.title': 'Terminal',
  'terminal.clear': 'Clear terminal',
  'terminal.expand': 'Expand terminal',
  'terminal.collapse': 'Collapse terminal',
  'terminal.empty': 'No terminal output yet.',
  'terminal.aria': 'Terminal output',

  // Common
  'common.error': 'Error',
  'common.success': 'Success',
  'common.notification': 'Notification',
  'common.confirm': 'Confirm',
  'common.cancel': 'Cancel',
  'common.unknown': 'unknown',
  'common.expand': 'Expand',
  'common.collapse': 'Collapse',
  'common.loading': 'Loading...',
  'common.copy_path': 'Copy path',

  // Toasts and tasks
  'toast.starting': 'Starting...',
  'toast.completed': 'Completed successfully!',
  'toast.task_failed': 'Task failed.',
  'toast.please_wait': 'Please wait',
  'toast.phase.download': 'Downloading',
  'toast.phase.install': 'Installing',
  'toast.installing_after_download': 'Download complete, installing...',
  'task.plugin.add': 'Adding plugin: {name}',
  'task.plugin.remove': 'Removing plugin: {name}',
  'task.version.install': 'Installing {name}@{version}',
  'task.version.uninstall': 'Uninstalling {name}@{version}',
  'task.version.switch': 'Switching {name} to {version}',
  'task.version.unset': 'Unsetting {name}',
  'task.custom.use': 'Using {name} ({path})',

  // Platform sentences shared by Settings and the SDK pages
  'platform.admin.required': ' Administrator approval is required.',
  'platform.restart.windows': 'Open a new terminal after changing PATH.',
  'platform.restart.darwin': 'Open a new terminal or run `source ~/.zprofile`.',
  'platform.restart.linux': 'Open a new terminal or run `source ~/.profile`.',
  'platform.restart.default': 'Open a new terminal after changing PATH.',
};
