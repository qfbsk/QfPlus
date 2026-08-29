export namespace model {
	
	export class CoreInfo {
	    currentVersion: string;
	    bundledVersion: string;
	    usesLocalCore: boolean;
	    executablePath: string;
	    osArch: string;
	    latestVersion: string;
	    updateAvailable: boolean;
	    releaseNotes: string;
	    releaseUrl: string;
	    autoUpdate: boolean;
	    lastCheck: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new CoreInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.currentVersion = source["currentVersion"];
	        this.bundledVersion = source["bundledVersion"];
	        this.usesLocalCore = source["usesLocalCore"];
	        this.executablePath = source["executablePath"];
	        this.osArch = source["osArch"];
	        this.latestVersion = source["latestVersion"];
	        this.updateAvailable = source["updateAvailable"];
	        this.releaseNotes = source["releaseNotes"];
	        this.releaseUrl = source["releaseUrl"];
	        this.autoUpdate = source["autoUpdate"];
	        this.lastCheck = source["lastCheck"];
	        this.error = source["error"];
	    }
	}
	export class CoreRelease {
	    version: string;
	    title: string;
	    date: string;
	    notes: string;
	    url: string;
	    isCurrent: boolean;
	    downloaded: boolean;
	
	    static createFrom(source: any = {}) {
	        return new CoreRelease(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.title = source["title"];
	        this.date = source["date"];
	        this.notes = source["notes"];
	        this.url = source["url"];
	        this.isCurrent = source["isCurrent"];
	        this.downloaded = source["downloaded"];
	    }
	}
	export class DownloadPathInfo {
	    path: string;
	    defaultPath: string;
	    isDefault: boolean;
	    hasMigratableData: boolean;
	
	    static createFrom(source: any = {}) {
	        return new DownloadPathInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.path = source["path"];
	        this.defaultPath = source["defaultPath"];
	        this.isDefault = source["isDefault"];
	        this.hasMigratableData = source["hasMigratableData"];
	    }
	}
	export class EnvironmentSdkEntry {
	    name: string;
	    source: string;
	    pluginAdded: boolean;
	    current: string;
	    versions: string[];
	    path: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentSdkEntry(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.source = source["source"];
	        this.pluginAdded = source["pluginAdded"];
	        this.current = source["current"];
	        this.versions = source["versions"];
	        this.path = source["path"];
	        this.version = source["version"];
	    }
	}
	export class EnvironmentPlatform {
	    os: string;
	    arch: string;
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentPlatform(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.arch = source["arch"];
	    }
	}
	export class EnvironmentGenerator {
	    app: string;
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentGenerator(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.app = source["app"];
	        this.version = source["version"];
	    }
	}
	export class EnvironmentDocument {
	    schemaVersion: number;
	    kind: string;
	    // Go type: time
	    generatedAt: any;
	    generator: EnvironmentGenerator;
	    platform: EnvironmentPlatform;
	    vfoxHome: string;
	    sdks: EnvironmentSdkEntry[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentDocument(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schemaVersion = source["schemaVersion"];
	        this.kind = source["kind"];
	        this.generatedAt = this.convertValues(source["generatedAt"], null);
	        this.generator = this.convertValues(source["generator"], EnvironmentGenerator);
	        this.platform = this.convertValues(source["platform"], EnvironmentPlatform);
	        this.vfoxHome = source["vfoxHome"];
	        this.sdks = this.convertValues(source["sdks"], EnvironmentSdkEntry);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class EnvironmentImportItem {
	    id: string;
	    name: string;
	    version: string;
	    source: string;
	    path: string;
	    resolution: string;
	    fallbackVersion: string;
	    action: string;
	    skipReason: string;
	    skipMessage: string;
	    current: boolean;
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentImportItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.id = source["id"];
	        this.name = source["name"];
	        this.version = source["version"];
	        this.source = source["source"];
	        this.path = source["path"];
	        this.resolution = source["resolution"];
	        this.fallbackVersion = source["fallbackVersion"];
	        this.action = source["action"];
	        this.skipReason = source["skipReason"];
	        this.skipMessage = source["skipMessage"];
	        this.current = source["current"];
	    }
	}
	export class EnvironmentImportPlan {
	    schemaVersion: number;
	    // Go type: time
	    generatedAt: any;
	    sourceVfoxHome: string;
	    items: EnvironmentImportItem[];
	    fallbackAllowed: boolean;
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentImportPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.schemaVersion = source["schemaVersion"];
	        this.generatedAt = this.convertValues(source["generatedAt"], null);
	        this.sourceVfoxHome = source["sourceVfoxHome"];
	        this.items = this.convertValues(source["items"], EnvironmentImportItem);
	        this.fallbackAllowed = source["fallbackAllowed"];
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EnvironmentImportResult {
	    imported: number;
	    skipped: number;
	    failed: number;
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentImportResult(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.imported = source["imported"];
	        this.skipped = source["skipped"];
	        this.failed = source["failed"];
	        this.warnings = source["warnings"];
	    }
	}
	export class SdkVersion {
	    version: string;
	
	    static createFrom(source: any = {}) {
	        return new SdkVersion(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	    }
	}
	export class SdkInfo {
	    name: string;
	    source: string;
	    path: string;
	    versions: SdkVersion[];
	
	    static createFrom(source: any = {}) {
	        return new SdkInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.source = source["source"];
	        this.path = source["path"];
	        this.versions = this.convertValues(source["versions"], SdkVersion);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class EnvironmentInventory {
	    addedPlugins: string[];
	    installedSdks: SdkInfo[];
	    customSdksMap: Record<string, Array<SdkInfo>>;
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentInventory(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.addedPlugins = source["addedPlugins"];
	        this.installedSdks = this.convertValues(source["installedSdks"], SdkInfo);
	        this.customSdksMap = this.convertValues(source["customSdksMap"], Array<SdkInfo>, true);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	
	export class SdkCommandStatus {
	    sdkName: string;
	    alias: string;
	    source: string;
	    managedBy: string;
	    resolved: boolean;
	    exePath: string;
	    exeDir: string;
	    onUserPath: boolean;
	    onMachinePath: boolean;
	    version: string;
	    isCurrent: boolean;
	    state: string;
	    notes: string[];
	
	    static createFrom(source: any = {}) {
	        return new SdkCommandStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sdkName = source["sdkName"];
	        this.alias = source["alias"];
	        this.source = source["source"];
	        this.managedBy = source["managedBy"];
	        this.resolved = source["resolved"];
	        this.exePath = source["exePath"];
	        this.exeDir = source["exeDir"];
	        this.onUserPath = source["onUserPath"];
	        this.onMachinePath = source["onMachinePath"];
	        this.version = source["version"];
	        this.isCurrent = source["isCurrent"];
	        this.state = source["state"];
	        this.notes = source["notes"];
	    }
	}
	export class EnvironmentStatusReport {
	    // Go type: time
	    generatedAt: any;
	    vfoxHome: string;
	    vfoxInPath: boolean;
	    shimDir: string;
	    pathDrift: boolean;
	    items: SdkCommandStatus[];
	    warnings: string[];
	
	    static createFrom(source: any = {}) {
	        return new EnvironmentStatusReport(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.generatedAt = this.convertValues(source["generatedAt"], null);
	        this.vfoxHome = source["vfoxHome"];
	        this.vfoxInPath = source["vfoxInPath"];
	        this.shimDir = source["shimDir"];
	        this.pathDrift = source["pathDrift"];
	        this.items = this.convertValues(source["items"], SdkCommandStatus);
	        this.warnings = source["warnings"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class MigrationItem {
	    name: string;
	    kind: string;
	    willMove: boolean;
	    count: number;
	    sizeBytes: number;
	    reason?: string;
	
	    static createFrom(source: any = {}) {
	        return new MigrationItem(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.kind = source["kind"];
	        this.willMove = source["willMove"];
	        this.count = source["count"];
	        this.sizeBytes = source["sizeBytes"];
	        this.reason = source["reason"];
	    }
	}
	export class MigrationPlan {
	    sourcePath: string;
	    targetPath: string;
	    movableItems: MigrationItem[];
	    excludedItems: MigrationItem[];
	    totalCount: number;
	    totalSizeBytes: number;
	    blockingReason?: string;
	
	    static createFrom(source: any = {}) {
	        return new MigrationPlan(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.sourcePath = source["sourcePath"];
	        this.targetPath = source["targetPath"];
	        this.movableItems = this.convertValues(source["movableItems"], MigrationItem);
	        this.excludedItems = this.convertValues(source["excludedItems"], MigrationItem);
	        this.totalCount = source["totalCount"];
	        this.totalSizeBytes = source["totalSizeBytes"];
	        this.blockingReason = source["blockingReason"];
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	export class PlatformInfo {
	    os: string;
	    name: string;
	    coreOS: string;
	    coreArch: string;
	    downloadPath: string;
	    defaultDownloadPath: string;
	    vfoxPathTarget: string;
	    sdkPathTarget: string;
	    shellProfile: string;
	    requiresElevation: boolean;
	    restartHint: string;
	
	    static createFrom(source: any = {}) {
	        return new PlatformInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.os = source["os"];
	        this.name = source["name"];
	        this.coreOS = source["coreOS"];
	        this.coreArch = source["coreArch"];
	        this.downloadPath = source["downloadPath"];
	        this.defaultDownloadPath = source["defaultDownloadPath"];
	        this.vfoxPathTarget = source["vfoxPathTarget"];
	        this.sdkPathTarget = source["sdkPathTarget"];
	        this.shellProfile = source["shellProfile"];
	        this.requiresElevation = source["requiresElevation"];
	        this.restartHint = source["restartHint"];
	    }
	}
	export class PluginInfo {
	    name: string;
	    isAdded: boolean;
	    isOfficial: boolean;
	    url: string;
	
	    static createFrom(source: any = {}) {
	        return new PluginInfo(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.isAdded = source["isAdded"];
	        this.isOfficial = source["isOfficial"];
	        this.url = source["url"];
	    }
	}
	export class ProxyNode {
	    name: string;
	    type: string;
	    delay: number;
	
	    static createFrom(source: any = {}) {
	        return new ProxyNode(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.delay = source["delay"];
	    }
	}
	export class ProxyGroup {
	    name: string;
	    type: string;
	    now: string;
	    nodes: ProxyNode[];
	
	    static createFrom(source: any = {}) {
	        return new ProxyGroup(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.type = source["type"];
	        this.now = source["now"];
	        this.nodes = this.convertValues(source["nodes"], ProxyNode);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	export class ProxyQuickStatus {
	    enabled: boolean;
	    running: boolean;
	    hasConfig: boolean;
	    exitGroup: string;
	    exitNode: string;
	    delay: number;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxyQuickStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.running = source["running"];
	        this.hasConfig = source["hasConfig"];
	        this.exitGroup = source["exitGroup"];
	        this.exitNode = source["exitNode"];
	        this.delay = source["delay"];
	        this.error = source["error"];
	    }
	}
	export class ProxyStatus {
	    enabled: boolean;
	    running: boolean;
	    hasConfig: boolean;
	    mixedPort: number;
	    subscriptionUrl: string;
	    selectedGroup: string;
	    selectedNode: string;
	    error?: string;
	
	    static createFrom(source: any = {}) {
	        return new ProxyStatus(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.enabled = source["enabled"];
	        this.running = source["running"];
	        this.hasConfig = source["hasConfig"];
	        this.mixedPort = source["mixedPort"];
	        this.subscriptionUrl = source["subscriptionUrl"];
	        this.selectedGroup = source["selectedGroup"];
	        this.selectedNode = source["selectedNode"];
	        this.error = source["error"];
	    }
	}
	
	export class SdkVersionDetail {
	    version: string;
	    isCurrent: boolean;
	
	    static createFrom(source: any = {}) {
	        return new SdkVersionDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.version = source["version"];
	        this.isCurrent = source["isCurrent"];
	    }
	}
	export class SdkDetail {
	    name: string;
	    current: string;
	    versions: SdkVersionDetail[];
	
	    static createFrom(source: any = {}) {
	        return new SdkDetail(source);
	    }
	
	    constructor(source: any = {}) {
	        if ('string' === typeof source) source = JSON.parse(source);
	        this.name = source["name"];
	        this.current = source["current"];
	        this.versions = this.convertValues(source["versions"], SdkVersionDetail);
	    }
	
		convertValues(a: any, classs: any, asMap: boolean = false): any {
		    if (!a) {
		        return a;
		    }
		    if (a.slice && a.map) {
		        return (a as any[]).map(elem => this.convertValues(elem, classs));
		    } else if ("object" === typeof a) {
		        if (asMap) {
		            for (const key of Object.keys(a)) {
		                a[key] = new classs(a[key]);
		            }
		            return a;
		        }
		        return new classs(a);
		    }
		    return a;
		}
	}
	
	

}

