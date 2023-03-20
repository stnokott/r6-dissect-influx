import type { ConnectionDetails } from "../../db"
import type { config } from "../../../wailsjs/go/models"
import Settings from "../Settings.svelte"

const mockSavedConfig: config.Config = {
	game: {
		install_dir: "/user/foo"
	},
	influx_db: {
		host: "myhost",
		port: 8086,
		org: "myorg",
		bucket: "mybucket",
		token: "mytoken"
	},
	convertValues: () => { }
}
const mockUserChosenGameDir = ""
const mockAutodetectedGameDir = ""

const mockConnectionDetails: ConnectionDetails = {
	Name: "Foo Bar DB",
	Version: "1.2.3",
	Commit: "aaaaa111"
}

beforeEach(() => {
	cy.window().then((w) => {
		w["go"] = {
			main: {
				App: {
					GetConfig: async () => mockSavedConfig,
					SaveAndValidateConfig: async () => mockConnectionDetails,
					OpenGameDirDialog: async () => mockUserChosenGameDir,
					AutodetectGameDir: async () => mockAutodetectedGameDir,
					ValidateGameDir: async (_s: string) => { throw "fails on purpose" },
					ValidateInfluxHost: async (_s: string) => { throw "fails on purpose" },
					ValidateInfluxPort: async (_s: string) => { throw "fails on purpose" },
					ValidateInfluxOrg: async (_s: string) => { throw "fails on purpose" },
					ValidateInfluxBucket: async (_s: string) => { throw "fails on purpose" },
					ValidateInfluxToken: async (_s: string) => { throw "fails on purpose" },
				}
			}
		}
		cy.spy(w["go"]["main"]["App"], "GetConfig").as("SpyGetConfig")
		cy.spy(w["go"]["main"]["App"], "SaveAndValidateConfig").as("SpySaveAndValidateConfig")
		cy.spy(w["go"]["main"]["App"], "OpenGameDirDialog").as("SpyOpenGameDirDialog")
		cy.spy(w["go"]["main"]["App"], "AutodetectGameDir").as("SpyAutodetectGameDir")
		cy.spy(w["go"]["main"]["App"], "ValidateGameDir").as("SpyValidateGameDir")
		cy.spy(w["go"]["main"]["App"], "ValidateInfluxHost").as("SpyValidateInfluxHost")
		cy.spy(w["go"]["main"]["App"], "ValidateInfluxPort").as("SpyValidateInfluxPort")
		cy.spy(w["go"]["main"]["App"], "ValidateInfluxOrg").as("SpyValidateInfluxOrg")
		cy.spy(w["go"]["main"]["App"], "ValidateInfluxBucket").as("SpyValidateInfluxBucket")
		cy.spy(w["go"]["main"]["App"], "ValidateInfluxToken").as("SpyValidateInfluxToken")
	})
})

describe('Settings', () => {
	it('should not be visible by default', () => {
		cy.mount(Settings)
		cy.get("div[data-cy-root] > *").should("not.be.visible")
	})
	it('should be visible if opened', () => {
		cy.mount(Settings, { props: { open: true } })
		cy.get("div[data-cy-root] > *").should("be.visible")
	})
	it('should load config', () => {
		cy.mount(Settings, { props: { open: true } })
		cy.get("@SpyGetConfig").should("be.called").then(() => {
			cy.get("input").should("have.attr", "data-invalid", "true")
			cy.get("input").eq(0).should("have.value", mockSavedConfig.game.install_dir)
			cy.get("input").eq(1).should("have.value", mockSavedConfig.influx_db.host)
			cy.get("input").eq(2).should("have.value", mockSavedConfig.influx_db.port.toString())
			cy.get("input").eq(3).should("have.value", mockSavedConfig.influx_db.org)
			cy.get("input").eq(4).should("have.value", mockSavedConfig.influx_db.bucket)
			cy.get("input").eq(5).should("have.value", mockSavedConfig.influx_db.token)
		})
	})
	it("should call validators", () => {
		cy.mount(Settings, { props: { open: true } })
		cy.get("@SpyValidateGameDir").should("be.called")
		cy.get("input").eq(1).type("a").then(() => {
			// once at initialization
			// once when config is loaded
			// once when we type
			cy.get("@SpyValidateInfluxHost").should("be.calledThrice")
		})
	})
	it("should disable confirm button when validation fails", () => {
		cy.mount(Settings, { props: { open: true } })
		cy.get('div[class$="footer"] > button[class*="primary"]').should("be.disabled")
	})
	it("should disable confirm button even when one input is valid", () => {
		cy.window().then((w) => {
			// set validator that marks the field as valid
			w["go"]["main"]["App"]["ValidateGameDir"].restore()
			cy.stub(w["go"]["main"]["App"], "ValidateGameDir").callsFake(async (_s: string) => true)
			cy.mount(Settings, { props: { open: true } })

			// gameDir input should be valid
			cy.get("input").eq(0).should("not.have.attr", "data-invalid", "true")
			// next input should still be invalid
			cy.get("input").eq(1).should("have.attr", "data-invalid", "true")
			// submit button should be disabled
			cy.get('div[class$="footer"] > button[class*="primary"]').should("be.disabled")
		})
	})
	it("should enable confirm button when all inputs are valid", () => {
		cy.window().then((w) => {
			// set validator that marks the field as valid
			w["go"]["main"]["App"]["ValidateGameDir"].restore()
			w["go"]["main"]["App"]["ValidateInfluxHost"].restore()
			w["go"]["main"]["App"]["ValidateInfluxPort"].restore()
			w["go"]["main"]["App"]["ValidateInfluxOrg"].restore()
			w["go"]["main"]["App"]["ValidateInfluxBucket"].restore()
			w["go"]["main"]["App"]["ValidateInfluxToken"].restore()
			cy.stub(w["go"]["main"]["App"], "ValidateGameDir").callsFake(async (_s: string) => true)
			cy.stub(w["go"]["main"]["App"], "ValidateInfluxHost").callsFake(async (_s: string) => true)
			cy.stub(w["go"]["main"]["App"], "ValidateInfluxPort").callsFake(async (_s: string) => true)
			cy.stub(w["go"]["main"]["App"], "ValidateInfluxOrg").callsFake(async (_s: string) => true)
			cy.stub(w["go"]["main"]["App"], "ValidateInfluxBucket").callsFake(async (_s: string) => true)
			cy.stub(w["go"]["main"]["App"], "ValidateInfluxToken").callsFake(async (_s: string) => true)
			cy.mount(Settings, { props: { open: true } })

			// every input should be valid
			cy.get("input").each((el) => {
				cy.wrap(el).should("not.have.attr", "data-invalid", "true")
			})
			// submit button should be enabled 
			cy.get('div[class$="footer"] > button[class*="primary"]').should("not.be.disabled")
		})
	})
})
