import { onMounted, ref, watch } from 'vue';
import { setWindowLightTheme } from '../services/runtimeService';

type UseTerminalDockOptions = {
  onShowTerminal?: () => void;
};

const storageKey = 'qfplus-show-terminal';

export const useTerminalDock = (options: UseTerminalDockOptions = {}) => {
  // Renaming the executable moved the WebView2 profile, so the vfoxG-era key can
  // still be the only one holding a value.
  const stored = localStorage.getItem(storageKey) ?? localStorage.getItem('vfox-show-terminal');
  const showTerminal = ref(stored === 'true');

  watch(showTerminal, (visible) => {
    localStorage.setItem(storageKey, String(visible));
    if (visible) {
      options.onShowTerminal?.();
    }
  });

  onMounted(() => {
    setWindowLightTheme();
  });

  return { showTerminal };
};
