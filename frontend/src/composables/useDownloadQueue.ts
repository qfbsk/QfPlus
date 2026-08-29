import { computed, ref } from 'vue';
import { onRuntimeEvent } from '../services/runtimeService';

export type DownloadTaskKind = 'plugin';
export type DownloadTaskStatus = 'queued' | 'running' | 'done' | 'error';

export interface DownloadTask {
  id: string;
  kind: DownloadTaskKind;
  title: string;
  subtitle: string;
  /** total number of files/units in this job (import items, or 1 for a plugin) */
  unitsTotal: number;
  /** how many units have finished so far */
  unitsDone: number;
  /** 0..100 for the current unit; null = indeterminate (no byte progress available) */
  unitPercent: number | null;
  /** human readable speed like "1.2 MB/s", or '' when unknown */
  speed: string;
  /** coarse phase label: download / install / add_plugin / activate / register / skip */
  phase: string;
  status: DownloadTaskStatus;
  error: string;
  finishedAt: number;
}

export interface EnqueueOptions {
  kind: DownloadTaskKind;
  title: string;
  unitsTotal?: number;
  run: (task: DownloadTask) => Promise<void>;
}

// Module-level singletons so every caller shares one queue.
const tasks = ref<DownloadTask[]>([]);
const resolvers = new Map<string, () => void>();
const runners = new Map<string, (task: DownloadTask) => Promise<void>>();
let processing = false;
let seq = 0;
const AUTO_HIDE_MS = 6000;

function removeTask(id: string) {
  const idx = tasks.value.findIndex((t) => t.id === id);
  if (idx >= 0) tasks.value.splice(idx, 1);
}

function scheduleAutoHide(id: string) {
  window.setTimeout(() => {
    const task = tasks.value.find((t) => t.id === id);
    if (task && (task.status === 'done' || task.status === 'error')) {
      removeTask(id);
    }
  }, AUTO_HIDE_MS);
}

async function processQueue() {
  if (processing) return;
  processing = true;
  try {
    for (;;) {
      const next = tasks.value.find((t) => t.status === 'queued');
      if (!next) break;
      next.status = 'running';
      try {
        const run = runners.get(next.id);
        if (run) await run(next);
        if (next.status === 'running') next.status = 'done';
      } catch (err) {
        next.status = 'error';
        next.error =
          err instanceof Error ? err.message : typeof err === 'string' ? err : 'download failed';
        next.unitPercent = null;
      } finally {
        next.finishedAt = Date.now();
        if (next.status === 'done') {
          next.unitPercent = 100;
          next.unitsDone = next.unitsTotal;
        }
        const resolve = resolvers.get(next.id);
        if (resolve) {
          resolvers.delete(next.id);
          resolve();
        }
        runners.delete(next.id);
        scheduleAutoHide(next.id);
      }
    }
  } finally {
    processing = false;
  }
}

export function useDownloadQueue() {
  const visible = computed(() => tasks.value.length > 0);
  const activeCount = computed(
    () => tasks.value.filter((t) => t.status === 'queued' || t.status === 'running').length
  );
  const finishedCount = computed(
    () => tasks.value.filter((t) => t.status === 'done' || t.status === 'error').length
  );

  // Append to the tail so earlier jobs always run first (first-come-first-served).
  const enqueue = (opts: EnqueueOptions): Promise<void> => {
    const id = `dl-${++seq}`;
    const task: DownloadTask = {
      id,
      kind: opts.kind,
      title: opts.title,
      subtitle: '',
      unitsTotal: opts.unitsTotal && opts.unitsTotal > 0 ? opts.unitsTotal : 1,
      unitsDone: 0,
      unitPercent: null,
      speed: '',
      phase: '',
      status: 'queued',
      error: '',
      finishedAt: 0,
    };
    tasks.value.push(task);
    const done = new Promise<void>((resolve) => {
      resolvers.set(id, resolve);
    });
    runners.set(id, opts.run);
    void processQueue();
    return done;
  };

  const dismiss = (id: string) => {
    resolvers.delete(id);
    runners.delete(id);
    removeTask(id);
  };

  const clearFinished = () => {
    tasks.value = tasks.value.filter((t) => t.status === 'queued' || t.status === 'running');
  };

  return { tasks, visible, activeCount, finishedCount, enqueue, dismiss, clearFinished };
}

// Re-export so callers can subscribe to progress events without extra imports.
export { onRuntimeEvent };
