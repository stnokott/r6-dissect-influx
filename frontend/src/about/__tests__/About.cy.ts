import type { AppInfo, EventNames, ReleaseInfo } from "../../app"
import About from "../About.svelte"

const mockAppInfo: AppInfo = {
	ProjectName: "Foo Bar",
	Version: "v0.0.1",
	Commit: "abde123",
	GithubURL: "github.com/stnokott/r6-dissect"
}
const mockEventNames: EventNames = {
	NewRound: "NEW_ROUND",
	RoundWatcherStarted: "ROUND_WATCHER_STOPPED",
	RoundWatcherError: "ROUND_WATCHER_ERROR",
	RoundWatcherStopped: "ROUND_WATCHER_STOPPED",
	RoundPush: "ROUND_PUSH",
	LatestReleaseInfo: "LATEST_RELEASE_INFO",
	LatestReleaseInfoErr: "LATEST_RELEASE_INFO_ERR",
	UpdateProgress: "UPDATE_PROGRESS",
	UpdateErr: "UPDATE_ERR"
}

beforeEach(() => {
	cy.window().then((w) => {
		w["go"] = {
			main: {
				App: {
					GetAppInfo: async () => mockAppInfo,
					GetEventNames: async () => mockEventNames,
					RequestLatestReleaseInfo: async () => { },
					StartUpdate: async () => { },
				}
			}
		}
		cy.spy(w["go"]["main"]["App"], "GetAppInfo").as("SpyGetAppInfo")
		cy.spy(w["go"]["main"]["App"], "GetEventNames").as("SpyGetEventNames")
		cy.spy(w["go"]["main"]["App"], "RequestLatestReleaseInfo").as("SpyRequestLatestReleaseInfo")
		cy.spy(w["go"]["main"]["App"], "StartUpdate").as("SpyStartUpdate")
		w["runtime"] = {
			EventsOn: () => { }
		}
		cy.spy(w["runtime"], "EventsOn").as("SpyEventsOn")
	})
})

describe('Settings', () => {
	it('should not be visible by default', () => {
		cy.mount(About)
		cy.get("div[data-cy-root] > *").should("not.be.visible")
	})
	it('should be visible if opened', () => {
		cy.mount(About, { props: { open: true } })
		cy.get("div[data-cy-root] > *").should("be.visible")
	})
	it("should display app info", () => {
		cy.mount(About, { props: { open: true } })
		cy.get("#tag-current-version").contains(mockAppInfo.Version)
		cy.get("#tag-current-version").contains(mockAppInfo.Commit)
	})
	function testReleaseInfo(releaseInfo: ReleaseInfo) {
		cy.window().then((w) => {
			const eventCbs: { [k: string]: (d: any) => void } = {};
			w["runtime"]["EventsOn"].restore()
			cy.stub(w["runtime"], "EventsOn").as("StubEventsOn").callsFake((eventName: string, cb: (d: any) => void) => {
				eventCbs[eventName] = cb;
			})
			cy.mount(About, { props: { open: true } })
			cy.get("@StubEventsOn").should("be.calledWith", mockEventNames.LatestReleaseInfo).then(() => {
				// simulate new release information
				eventCbs[mockEventNames.LatestReleaseInfo](releaseInfo)
			})
		})
	}
	it("should process newer release info", () => {
		const releaseInfo: ReleaseInfo = {
			Version: "1.0.0",
			IsNewer: true,
			Commitish: "456def",
			PublishedAt: "2023-01-01 00:00:00",
			Changelog: "Added foo and bar"
		}
		testReleaseInfo(releaseInfo)
		cy.get("#tag-latest-version").contains(releaseInfo.Version)
		cy.get("#update-buttons-container > button").contains("Apply").should("not.be.disabled")
	})
	it("should process older release info", () => {
		const releaseInfo: ReleaseInfo = {
			Version: "0.0.9",
			IsNewer: false,
			Commitish: "456def",
			PublishedAt: "2022-01-01 00:00:00",
			Changelog: "Added foo and bar"
		}
		testReleaseInfo(releaseInfo)
		cy.get("#tag-latest-version").contains(releaseInfo.Version)
		cy.get("#update-buttons-container button").contains("Apply").should("not.exist")
	})
	it("should request updates when button pressed", () => {
		cy.mount(About, { props: { open: true } })
		cy.get("#update-buttons-container button").contains("Check").click().then(() => {
			cy.get("@SpyRequestLatestReleaseInfo").should("have.been.calledOnce")
			cy.get("#update-buttons-container button").contains("Check").should("be.disabled")
		})
	})
	it("should initiate update", () => {
		const releaseInfo: ReleaseInfo = {
			Version: "1.0.0",
			IsNewer: true,
			Commitish: "456def",
			PublishedAt: "2023-01-01 00:00:00",
			Changelog: "Added foo and bar"
		}
		testReleaseInfo(releaseInfo)
		cy.get('div[data-cy="loader"').should("not.be.visible")
		cy.get("#update-buttons-container button").contains("Apply").click().then(() => {
			cy.get("@SpyStartUpdate").should("have.been.calledOnce")
			cy.get("#update-buttons-container button").contains("Apply").should("be.disabled")
			cy.get('div[data-cy="loader"').should("be.visible")
		})
	})
})
