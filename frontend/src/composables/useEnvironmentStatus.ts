import { ref } from 'vue';
import {
  fetchEnvironmentStatus,
  openEnvironmentDiagnostic,
  type EnvironmentStatusReport,
} from '../services/environmentService';

const loading = ref(false);
const report = ref<EnvironmentStatusReport | null>(null);
const error = ref<string>('');

export function useEnvironmentStatus() {
  const load = async () => {
    loading.value = true;
    error.value = '';
    try {
      report.value = await fetchEnvironmentStatus();
    } catch (err) {
      error.value = err instanceof Error ? err.message : String(err);
    } finally {
      loading.value = false;
    }
  };

  const openDiagnostic = async () => {
    error.value = '';
    try {
      await openEnvironmentDiagnostic();
    } catch (err) {
      error.value = err instanceof Error ? err.message : String(err);
      throw err;
    }
  };

  return {
    loading,
    report,
    error,
    load,
    openDiagnostic,
  };
}
