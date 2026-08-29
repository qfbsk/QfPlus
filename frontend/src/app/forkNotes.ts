export type ForkNoteGroupId = 'design' | 'engine' | 'network' | 'platform';

export interface ForkNoteGroup {
  id: ForkNoteGroupId;
  changes: string[];
}

/**
 * The single source of truth for "what this version changed on top of upstream".
 * Rendered by both the onboarding sheet and the About page; text lives in
 * i18n as `about.group.<id>` and `about.change.<id>`.
 */
export const forkNoteGroups: ForkNoteGroup[] = [
  {
    id: 'design',
    changes: [
      'design.minimal',
      'design.tokens',
      'design.onboarding',
    ],
  },
  {
    id: 'engine',
    changes: [
      'engine.core_manager',
      'engine.auto_update',
      'engine.streaming',
      'engine.single_lock',
      'engine.market',
      'engine.profile',
    ],
  },
  {
    id: 'network',
    changes: [
      'network.proxy',
      'network.scope',
      'network.ua',
      'network.quick',
      'network.no_secrets',
    ],
  },
  {
    id: 'platform',
    changes: [
      'platform.path',
      'platform.metadata',
      'platform.shims',
      'platform.repair',
      'platform.storage',
      'platform.diagnostics',
      'platform.packaging',
      'platform.i18n',
    ],
  },
];

export const forkChangeCount = forkNoteGroups.reduce(
  (total, group) => total + group.changes.length,
  0,
);
