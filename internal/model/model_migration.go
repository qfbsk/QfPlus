package model

type MigrationProgress struct {
	Stage              string `json:"stage"`
	Current            string `json:"current"`
	Completed          int    `json:"completed"`
	Total              int    `json:"total"`
	Percent            int    `json:"percent"`
	EstimatedRemaining int    `json:"estimatedRemaining"`
}

// MigrationItemKind classifies one entry in a storage migration preview.
type MigrationItemKind string

const (
	MigrationItemKindVersion    MigrationItemKind = "version"
	MigrationItemKindPlugin     MigrationItemKind = "plugin"
	MigrationItemKindJunction   MigrationItemKind = "junction"
	MigrationItemKindShim       MigrationItemKind = "shim"
	MigrationItemKindSelection  MigrationItemKind = "selection"
	MigrationItemKindMetadata   MigrationItemKind = "metadata"
	MigrationItemKindThirdParty MigrationItemKind = "third_party"
	MigrationItemKindCustom     MigrationItemKind = "custom"
	MigrationItemKindOther      MigrationItemKind = "other"
)

// MigrationItem is one row of the migration preview: either something that will
// be copied (WillMove=true) or something that is only listed (WillMove=false).
type MigrationItem struct {
	Name       string           `json:"name"`
	Kind       MigrationItemKind `json:"kind"`
	WillMove   bool             `json:"willMove"`
	Count      int              `json:"count"`
	SizeBytes  int64            `json:"sizeBytes"`
	Reason     string           `json:"reason,omitempty"`
}

// MigrationPlan is a read-only preview of what a storage migration would do.
// It never mutates state; the actual copy is performed by migrateVfoxHomeData.
type MigrationPlan struct {
	SourcePath     string          `json:"sourcePath"`
	TargetPath     string          `json:"targetPath"`
	MovableItems   []MigrationItem `json:"movableItems"`
	ExcludedItems  []MigrationItem `json:"excludedItems"`
	TotalCount     int             `json:"totalCount"`
	TotalSizeBytes int64           `json:"totalSizeBytes"`
	BlockingReason string          `json:"blockingReason,omitempty"`
}
