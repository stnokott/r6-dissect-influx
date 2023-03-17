import MatchProgressIndicator from "../MatchProgressIndicator.svelte";
import MatchProgressIndicatorWrapper from "./MatchProgressIndicatorWrapper.svelte";
import { createRoundInfo } from "./util";

describe('MatchProgressIndicator', () => {
  it('should be empty by default', () => {
    cy.mount(MatchProgressIndicator, { props: { fullWidth: false } })
    cy.get("#progress-bar").should("be.empty")
  })
  it('should pass style props down', () => {
    cy.mount(MatchProgressIndicator, { props: { style: "margin-bottom: 0" } })
    cy.get("#progress").should("have.css", "margin-bottom", "0px")
  })
  it('should render the correct amount of rounds', () => {
    const round1 = createRoundInfo(false, "ATTACK");
    const round2 = createRoundInfo(true, "ATTACK");
    const round3 = createRoundInfo(true, "DEFENSE");
    cy.mount(MatchProgressIndicatorWrapper, { props: { Component: MatchProgressIndicator, roundInfos: [round1, round2, round3] } })
    cy.get("#progress-bar .step").should("have.length", 3)
  })
})
