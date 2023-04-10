import type { AppInfo, EventNames, ReleaseInfo } from "../../app"
import type { ConnectionDetails } from "../../db"
import Footer from "../Footer.svelte"

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
					StartReleaseWatcher: async () => { },
				}
			}
		}
		cy.spy(w["go"]["main"]["App"], "GetAppInfo").as("SpyGetAppInfo")
		cy.spy(w["go"]["main"]["App"], "GetEventNames").as("SpyGetEventNames")
		cy.spy(w["go"]["main"]["App"], "StartReleaseWatcher").as("SpyStartReleaseWatcher")
		w["runtime"] = {
			EventsOn: () => { }
		}
		cy.spy(w["runtime"], "EventsOn").as("SpyEventsOn")
	})
})

describe('Footer', () => {
	it('should be visible', () => {
		cy.mount(Footer)
		cy.get("#root").should("be.visible")
	})
	it('should show app info', () => {
		cy.mount(Footer)
		cy.get("#root")
			.contains(mockAppInfo.Version)
			.contains(mockAppInfo.Commit)
	})
	it("should query app info", () => {
		cy.mount(Footer)
		cy.get("@SpyGetAppInfo").should("be.called")
	})
	it("should request release information", () => {
		cy.mount(Footer)
		cy.get("@SpyStartReleaseWatcher").should("be.called")
		cy.get("@SpyGetEventNames").then((SpyGetEventNames) => {
			cy.get("@SpyEventsOn").should("be.calledAfter", SpyGetEventNames)
			cy.get("@SpyEventsOn").should("be.calledWith", mockEventNames.LatestReleaseInfo)
			cy.get("@SpyEventsOn").should("be.calledWith", mockEventNames.LatestReleaseInfoErr)
		})
	})
	it("should notify about new release", () => {
		cy.window().then((w) => {
			const eventCbs: { [k: string]: (d: any) => void } = {};
			// save callbacks in dictionary
			w["runtime"]["EventsOn"].restore()
			cy.stub(w["runtime"], "EventsOn").as("StubEventsOn").callsFake((eventName: string, cb: (d: any) => void) => {
				eventCbs[eventName] = cb;
			})
			cy.mount(Footer)
			cy.get("@StubEventsOn").should("be.calledWith", mockEventNames.LatestReleaseInfo).then(() => {
				// simulate new release information
				const releaseInfo: ReleaseInfo = {
					Version: "1.0.0",
					IsNewer: true,
					Commitish: "456def",
					PublishedAt: "2023-01-01 00:00:00",
					Changelog: "Added foo and bar"
				}
				eventCbs[mockEventNames.LatestReleaseInfo](releaseInfo)
				cy.get("#root").contains("Update available")
			})
		})
	})
	it("should show db connection details", () => {
		const connectionDetails: ConnectionDetails = {
			Name: "Foo Bar DB",
			Version: "1.2.3",
			Commit: "aaaaa111"
		}
		const promConnectionDetails = async () => connectionDetails
		cy.mount(Footer, { props: { promConnectionDetails: promConnectionDetails() } })
		cy.get("#connection-details")
			.contains(connectionDetails.Name)
			.contains(connectionDetails.Version)
			.contains(connectionDetails.Commit)
	})
})
