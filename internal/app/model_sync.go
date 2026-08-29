package app

import (
	"time"

	"QfPlus/internal/model"
)

type sdkEnvironmentExport struct {
	GeneratedAt  time.Time
	Platform     model.PlatformInfo
	VfoxInPath   bool
	PathOverride bool
	VfoxSdks     []sdkEnvironmentVfoxSdk
	SystemSdks   []model.SdkInfo
	CustomSdks   map[string][]model.SdkInfo
	Warnings     []string
}

type sdkEnvironmentVfoxSdk struct {
	Name             string
	Versions         []model.SdkVersion
	Detail           model.SdkDetail
	VersionPaths     map[string]string
	ActiveCustomPath string
}
