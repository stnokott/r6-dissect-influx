import MatchProgressStep from "../MatchProgressStep.svelte";
import { createRoundInfo } from "./util";

describe('MatchProgressIndicator', () => {
  it('should have correct class based on win/loss', () => {
    const roundInfo = createRoundInfo(true, "Attack");
    cy.mount(MatchProgressStep, { props: { roundInfo: roundInfo } });
    cy.get(".step").should("have.class", "won");
  })

  it('should render attack icon correctly', () => {
    const roundInfo = createRoundInfo(true, "Attack");
    cy.mount(MatchProgressStep, { props: { roundInfo: roundInfo } });
    cy.get(".step").find(".icon.attack").should("exist");
  })

  it('should render defense icon correctly', () => {
    const roundInfo = createRoundInfo(true, "Defense");
    cy.mount(MatchProgressStep, { props: { roundInfo: roundInfo } });
    cy.get(".step").find(".icon.defense").should("exist");
  })

  it('should have default state indicator', () => {
    const roundInfo = createRoundInfo(true, "Defense");
    cy.mount(MatchProgressStep, { props: { roundInfo: roundInfo } });
    cy.get(".status").should("exist");
    cy.get(".status > svg").should("exist");
  })

  it('should have defined state indicator', () => {
    const roundInfo = createRoundInfo(true, "Defense");
    cy.mount(MatchProgressStep, { props: { roundInfo: roundInfo, status: "done" } });
    cy.get(".status").should("exist");
    cy.get(".status > svg").should("exist");
  })
})
