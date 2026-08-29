package app

import (
	"runtime"
	"sort"
	"strings"
	"time"

	"QfPlus/internal/model"
)

// buildEnvironmentDocument turns the internal snapshot into the public JSON document.
func buildEnvironmentDocument(snapshot sdkEnvironmentExport) model.EnvironmentDocument {
	doc := model.EnvironmentDocument{
		SchemaVersion: 1,
		Kind:          "qfplus.environment",
		GeneratedAt:   snapshot.GeneratedAt,
		Generator: model.EnvironmentGenerator{
			App:     "QfPlus",
			Version: productVersion(),
		},
		Platform: model.EnvironmentPlatform{
			OS:   runtime.GOOS,
			Arch: runtime.GOARCH,
		},
		VfoxHome: snapshot.Platform.DownloadPath,
		Warnings: append([]string(nil), snapshot.Warnings...),
	}

	seen := make(map[string]bool)
	for _, sdk := range snapshot.VfoxSdks {
		name := strings.ToLower(strings.TrimSpace(sdk.Name))
		if name == "" {
			continue
		}
		id := name + "vfox"
		if seen[id] {
			continue
		}
		seen[id] = true

		versions := make([]string, 0, len(sdk.Versions))
		for _, v := range sdk.Versions {
			if v.Version != "" {
				versions = append(versions, v.Version)
			}
		}
		current := sdk.Detail.Current
		if current == "" && len(versions) == 1 {
			current = versions[0]
		}
		doc.Sdks = append(doc.Sdks, model.EnvironmentSdkEntry{
			Name:        sdk.Name,
			Source:      string(model.EnvironmentSdkSourceVfox),
			PluginAdded: true,
			Current:     current,
			Versions:    versions,
		})
	}

	for _, sdk := range snapshot.SystemSdks {
		name := strings.ToLower(strings.TrimSpace(sdk.Name))
		if name == "" {
			continue
		}
		id := name + "system"
		if seen[id] {
			continue
		}
		seen[id] = true
		version := ""
		if len(sdk.Versions) > 0 {
			version = sdk.Versions[0].Version
		}
		doc.Sdks = append(doc.Sdks, model.EnvironmentSdkEntry{
			Name:    sdk.Name,
			Source:  string(model.EnvironmentSdkSourceSystem),
			Path:    sdk.Path,
			Version: version,
		})
	}

	names := make([]string, 0, len(snapshot.CustomSdks))
	for name := range snapshot.CustomSdks {
		names = append(names, name)
	}
	sort.Strings(names)
	for _, name := range names {
		id := strings.ToLower(strings.TrimSpace(name)) + "custom"
		if seen[id] {
			continue
		}
		seen[id] = true
		for _, sdk := range snapshot.CustomSdks[name] {
			version := ""
			if len(sdk.Versions) > 0 {
				version = sdk.Versions[0].Version
			}
			doc.Sdks = append(doc.Sdks, model.EnvironmentSdkEntry{
				Name:    sdk.Name,
				Source:  string(model.EnvironmentSdkSourceCustom),
				Path:    sdk.Path,
				Version: version,
			})
		}
	}

	return doc
}

// productVersion returns the application version recorded in build metadata.
func productVersion() string {
	return "1.0.0"
}

// environmentDocumentNow returns the current environment as a model.EnvironmentDocument.
func (a *App) environmentDocumentNow() model.EnvironmentDocument {
	return buildEnvironmentDocument(a.collectSdkEnvironmentExport(time.Now()))
}
