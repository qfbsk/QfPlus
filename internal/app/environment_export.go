package app

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// ExportEnvironmentJson saves the current environment document as a .json file.
func (a *App) exportEnvironmentJson() (string, error) {
	if a.ctx == nil {
		return "", fmt.Errorf("application context is not ready")
	}

	doc := a.environmentDocumentNow()
	defaultDir := ""
	if home, err := os.UserHomeDir(); err == nil {
		defaultDir = home
	}

	path, err := runtime.SaveFileDialog(a.ctx, runtime.SaveDialogOptions{
		Title:            "Export QfPlus environment",
		DefaultDirectory: defaultDir,
		DefaultFilename:  fmt.Sprintf("QfPlus-environment-%s", time.Now().Format("20060102-150405")),
		Filters: []runtime.FileFilter{
			{DisplayName: "JSON Files (*.json)", Pattern: "*.json"},
		},
		CanCreateDirectories: true,
	})
	if err != nil {
		return "", err
	}
	path = strings.TrimSpace(path)
	if path == "" {
		return "", nil
	}
	if !strings.HasSuffix(strings.ToLower(path), ".json") {
		path += ".json"
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return "", err
	}

	data, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal environment document: %w", err)
	}
	data = append(data, '\n')
	if err := os.WriteFile(path, data, 0644); err != nil {
		return "", err
	}
	return path, nil
}

// PreviewEnvironmentDocument returns the current environment as a JSON-ready document.
func (a *App) previewEnvironmentDocument() (EnvironmentDocument, error) {
	if a.ctx == nil {
		return EnvironmentDocument{}, fmt.Errorf("application context is not ready")
	}
	return a.environmentDocumentNow(), nil
}
