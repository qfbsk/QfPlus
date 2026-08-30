package app

import "QfPlus/internal/model"

type AppConfig = model.AppConfig
type CoreInfo = model.CoreInfo
type CoreRelease = model.CoreRelease
type DownloadPathInfo = model.DownloadPathInfo
type GitHubSource = model.GitHubSource
type GitHubSourceConfig = model.GitHubSourceConfig
type GitHubSourceSettings = model.GitHubSourceSettings
type MigrationProgress = model.MigrationProgress
type MigrationPlan = model.MigrationPlan
type MigrationItem = model.MigrationItem
type MigrationItemKind = model.MigrationItemKind
type PlatformInfo = model.PlatformInfo
type PluginInfo = model.PluginInfo
type ProxyConfig = model.ProxyConfig
type ProxyGroup = model.ProxyGroup
type ProxyNode = model.ProxyNode
type ProxyQuickStatus = model.ProxyQuickStatus
type ProxyStatus = model.ProxyStatus
type SdkDetail = model.SdkDetail
type SdkEnvironmentImportResult = model.SdkEnvironmentImportResult
type SdkInfo = model.SdkInfo
type SdkVersion = model.SdkVersion
type SdkVersionDetail = model.SdkVersionDetail
type EnvironmentStatusReport = model.EnvironmentStatusReport
type SdkCommandStatus = model.SdkCommandStatus
type EnvironmentDocument = model.EnvironmentDocument
type EnvironmentImportPlan = model.EnvironmentImportPlan
type EnvironmentImportItem = model.EnvironmentImportItem
type EnvironmentImportResult = model.EnvironmentImportResult
type EnvironmentImportProgress = model.EnvironmentImportProgress
type EnvironmentInventory = model.EnvironmentInventory

const (
	EnvironmentResolutionExact            = model.EnvironmentResolutionExact
	EnvironmentResolutionAlreadyInstalled = model.EnvironmentResolutionAlreadyInstalled
	EnvironmentResolutionUnavailable      = model.EnvironmentResolutionUnavailable
	EnvironmentResolutionFallback         = model.EnvironmentResolutionFallback
	EnvironmentResolutionPathMissing      = model.EnvironmentResolutionPathMissing
	EnvironmentResolutionInvalidName      = model.EnvironmentResolutionInvalidName
	EnvironmentResolutionNotExported        = model.EnvironmentResolutionNotExported
)
