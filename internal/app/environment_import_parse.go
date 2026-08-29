package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"QfPlus/internal/model"
)

const environmentDocumentKind = "qfplus.environment"

// EnvironmentImportRow is the raw data extracted from an import source.
type EnvironmentImportRow struct {
	Kind    string
	Name    string
	Version string
	Path    string
	Current bool
}

// parseEnvironmentDocument reads a JSON environment file and returns validated rows.
func parseEnvironmentDocument(path string) ([]EnvironmentImportRow, model.EnvironmentDocument, []string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, model.EnvironmentDocument{}, nil, err
	}

	// Strip UTF-8 BOM if present.
	data = bytes.TrimPrefix(data, []byte{0xef, 0xbb, 0xbf})

	var doc model.EnvironmentDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, model.EnvironmentDocument{}, nil, fmt.Errorf("invalid JSON: %w", err)
	}

	warnings := make([]string, 0)
	if doc.Kind != environmentDocumentKind {
		return nil, model.EnvironmentDocument{}, nil, fmt.Errorf("unexpected document kind %q, expected %q", doc.Kind, environmentDocumentKind)
	}
	if doc.SchemaVersion > 1 {
		return nil, model.EnvironmentDocument{}, nil, fmt.Errorf("unsupported environment document schema version %d", doc.SchemaVersion)
	}
	if doc.SchemaVersion < 1 {
		warnings = append(warnings, "Document schema version is missing; assuming v1.")
	}

	rows := make([]EnvironmentImportRow, 0, len(doc.Sdks))
	for _, sdk := range doc.Sdks {
		name := strings.TrimSpace(sdk.Name)
		if name == "" {
			warnings = append(warnings, "Skipped an SDK entry with no name.")
			continue
		}
		source := strings.ToLower(strings.TrimSpace(sdk.Source))
		switch source {
		case string(model.EnvironmentSdkSourceVfox):
			for _, v := range uniqueStrings(sdk.Versions) {
				rows = append(rows, EnvironmentImportRow{
					Kind:    "vfox",
					Name:    name,
					Version: v,
					Current: sdk.Current != "" && v == sdk.Current,
				})
			}
			// If only Current is set and no versions, import the current version.
			if len(sdk.Versions) == 0 && sdk.Current != "" {
				rows = append(rows, EnvironmentImportRow{
					Kind:    "vfox",
					Name:    name,
					Version: sdk.Current,
					Current: true,
				})
			}
		case string(model.EnvironmentSdkSourceCustom):
			rows = append(rows, EnvironmentImportRow{
				Kind:    "custom",
				Name:    name,
				Version: sdk.Version,
				Path:    sdk.Path,
				Current: false,
			})
		case string(model.EnvironmentSdkSourceSystem):
			// System SDKs are informational; they are not imported automatically.
		default:
			warnings = append(warnings, fmt.Sprintf("Skipped %s with unknown source %q.", name, sdk.Source))
		}
	}

	return rows, doc, warnings, nil
}

// readEnvironmentDocumentText reads the raw file for inspection (e.g. secret scanning).
func readEnvironmentDocumentText(path string) (string, error) {
	f, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer f.Close()
	b := new(strings.Builder)
	if _, err := io.Copy(b, f); err != nil {
		return "", err
	}
	return b.String(), nil
}
