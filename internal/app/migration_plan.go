package app

import (
	"os"
	"path/filepath"
	"sort"
	"strings"

	"QfPlus/internal/model"
)

// planEnvironmentMigration builds a read-only preview of what a storage
// migration would copy and what it would leave behind. It never mutates state;
// the actual copy is performed by migrateVfoxHomeData after the user confirms.
func (a *App) planEnvironmentMigration(targetPath string) (model.MigrationPlan, error) {
	source := a.getVfoxHome()
	normalized, err := normalizeDownloadPath(targetPath)
	if err != nil {
		return model.MigrationPlan{}, err
	}
	plan := model.MigrationPlan{
		SourcePath: source,
		TargetPath: normalized,
	}

	if strings.TrimSpace(source) == "" {
		plan.BlockingReason = "current storage path is not resolved"
		return plan, nil
	}
	if samePath(source, normalized) {
		plan.BlockingReason = "target is the same as the current storage path"
		return plan, nil
	}
	if isPathWithin(normalized, source) {
		plan.BlockingReason = "new storage path cannot be inside the current storage directory"
		return plan, nil
	}

	entries, err := os.ReadDir(source)
	if err != nil {
		if os.IsNotExist(err) {
			return plan, nil
		}
		return plan, err
	}

	plan.MovableItems = make([]model.MigrationItem, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if name == "gui-non-vfox-sdks.json" {
			continue
		}
		count, size, err := scanTreeSize(filepath.Join(source, name))
		if err != nil {
			return plan, err
		}
		plan.MovableItems = append(plan.MovableItems, model.MigrationItem{
			Name:      name,
			Kind:      classifyMigrationEntry(name),
			WillMove:  true,
			Count:     count,
			SizeBytes: size,
		})
		plan.TotalCount += count
		plan.TotalSizeBytes += size
	}

	// Third-party / custom SDKs live at external paths recorded in
	// gui-non-vfox-sdks.json. They are listed for transparency but never copied
	// and never deleted — the new storage root only references them.
	excluded := a.getNonVfoxSdksMap()
	names := make([]string, 0, len(excluded))
	for name := range excluded {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		list := excluded[name]
		if len(list) == 0 {
			continue
		}
		plan.ExcludedItems = append(plan.ExcludedItems, model.MigrationItem{
			Name:     name,
			Kind:     migrationExcludedKind(list[0].Source),
			WillMove: false,
			Count:    len(list),
			Reason:   "external path, listed only — not copied or deleted",
		})
	}

	sort.Slice(plan.MovableItems, func(i, j int) bool {
		return plan.MovableItems[i].Name < plan.MovableItems[j].Name
	})
	return plan, nil
}

func classifyMigrationEntry(name string) model.MigrationItemKind {
	switch strings.ToLower(name) {
	case "cache":
		return model.MigrationItemKindVersion
	case "plugin":
		return model.MigrationItemKindPlugin
	case "sdks":
		return model.MigrationItemKindVersion
	case "path-shims":
		return model.MigrationItemKindShim
	case ".vfox.toml", "hijacked_paths.json":
		return model.MigrationItemKindMetadata
	}
	if strings.HasPrefix(name, ".") {
		return model.MigrationItemKindMetadata
	}
	return model.MigrationItemKindOther
}

func migrationExcludedKind(source string) model.MigrationItemKind {
	if source == "system" {
		return model.MigrationItemKindThirdParty
	}
	return model.MigrationItemKindCustom
}
