import { onMounted, onUnmounted, ref } from 'vue';
import { onRuntimeEvent } from '../services/runtimeService';

const ONBOARDING_SEEN_KEY = 'qfplus-onboarding-seen';
const FORCE_QUERY_FLAG = 'onboarding';

const MIN_VISIBLE_MS = 1150;
const READY_TIMEOUT_MS = 3000;
const READY_ANIMATION_MS = 2010;
const LEAVE_ANIMATION_MS = 260;

export type OnboardingPhase = 'loading' | 'ready' | 'leaving';

export const shouldShowOnboarding = () => {
  if (new URLSearchParams(window.location.search).get(FORCE_QUERY_FLAG) === '1') {
    return true;
  }
  return localStorage.getItem(ONBOARDING_SEEN_KEY) !== '1';
};

export const useOnboarding = (onDismissed: () => void) => {
  const phase = ref<OnboardingPhase>('loading');
  const isAnimating = ref(true);

  const timerIds: number[] = [];
  const later = (delayMs: number, run: () => void) => {
    timerIds.push(window.setTimeout(run, delayMs));
  };

  let hasMinimumDelayElapsed = false;
  let hasSdkScanFinished = false;
  let offSdkScanReady: (() => void) | null = null;

  const markReady = () => {
    if (phase.value !== 'loading') {
      return;
    }
    phase.value = 'ready';
    later(READY_ANIMATION_MS, () => {
      isAnimating.value = false;
    });
  };

  const tryMarkReady = () => {
    if (hasMinimumDelayElapsed && hasSdkScanFinished) {
      markReady();
    }
  };

  onMounted(() => {
    offSdkScanReady = onRuntimeEvent<void>('system-sdks-ready', () => {
      hasSdkScanFinished = true;
      tryMarkReady();
    });
    later(MIN_VISIBLE_MS, () => {
      hasMinimumDelayElapsed = true;
      tryMarkReady();
    });
    later(READY_TIMEOUT_MS, markReady);
  });

  onUnmounted(() => {
    if (offSdkScanReady) {
      offSdkScanReady();
      offSdkScanReady = null;
    }
    timerIds.splice(0).forEach((timerId) => window.clearTimeout(timerId));
  });

  const requestEnter = () => {
    if (phase.value === 'leaving') {
      return;
    }
    phase.value = 'leaving';
    later(LEAVE_ANIMATION_MS, () => {
      localStorage.setItem(ONBOARDING_SEEN_KEY, '1');
      onDismissed();
    });
  };

  return { phase, isAnimating, requestEnter };
};
