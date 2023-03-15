import LoadingOverlay from "../LoadingOverlay.svelte"

describe('LoadingOverlay', () => {
	it('should be hidden if open is false', () => {
		cy.mount(LoadingOverlay, { props: { open: false, loadingDesc: "foo bar" } })
		cy.get("#root").should("not.be.visible")
	})
	it('should be visible if open is true', () => {
		cy.mount(LoadingOverlay, { props: { open: true, loadingDesc: "foo bar" } })
		cy.get("#root").should("be.visible")
	})
	it('should show loading description if not erroneous', () => {
		cy.mount(LoadingOverlay, { props: { open: true, loadingDesc: "foo bar" } })
		cy.get("[data-cy=loader]").contains(/^foo bar$/)
	})
	it('should not show loading description if erroneous', () => {
		cy.mount(LoadingOverlay, { props: { open: true, loadingDesc: "foo bar", errorTitle: "my error" } })
		cy.get("[data-cy=loader]").should("not.exist")
		cy.get("[data-cy=error]").should("be.visible")
	})
	it('should show error details if erroneous', () => {
		cy.mount(LoadingOverlay, { props: { open: true, loadingDesc: "foo bar", errorTitle: "my error title", errorDetail: "my error details" } })
		cy.get("[data-cy=loader]").should("not.exist")
		cy.get("[data-cy=error]").contains(/^my error title$/).should("be.visible")
		cy.get("[data-cy=error]").contains(/^my error details$/).should("be.visible")
	})
	it('should not have close button if loading', () => {
		cy.mount(LoadingOverlay, { props: { open: true, loadingDesc: "foo bar" } })
		cy.get("[data-cy=error] button[title^=Close]").should("not.exist")
	})
	it('should have close button if erroneous', () => {
		cy.mount(LoadingOverlay, { props: { open: true, loadingDesc: "foo bar", errorTitle: "my error title", errorDetail: "my error details" } })
		cy.get("[data-cy=error] button[title^=Close]").should("be.visible")
	})
	it('should close when close button clicked', () => {
		cy.mount(LoadingOverlay, { props: { open: true, loadingDesc: "foo bar", errorTitle: "my error title", errorDetail: "my error details" } })
		cy.get("#root").should("be.visible")
		cy.get("[data-cy=error] button[title^=Close]").click()
		cy.get("#root").should("not.be.visible")
	})
})
