import * as runtime from "../wailsjs/runtime/runtime";

export { };

declare global {
	interface Window {
		runtime: runtime;
	}
}

window.runtime = window.runtime || {};

export namespace app {
	export class AppInfo {
		ProjectName: string
		Version: string
		Commit: string
		GithubURL: string
	}
	export class ReleaseInfo {
		IsNewer: boolean
		Version: string
		Commitish: string
		PublishedAt: string
		Changelog: string
	}

	export class EventNames {
		NewRound: string
		RoundWatcherError: string
		RoundWatcherStopped: string
		LatestReleaseInfo: string
		LatestReleaseInfoErr: string
		UpdateProgress: string
		UpdateErr: string
	}
}

export namespace db {
	export class ConnectionDetails {
		Name: string
		Version: string
		Commit: string
	}
}
