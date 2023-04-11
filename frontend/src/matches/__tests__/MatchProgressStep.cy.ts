import MatchProgressStep from "../MatchProgressStep.svelte";
import { createRound } from "./util";

describe('MatchProgressIndicator', () => {
  it('should have correct class based on win/loss', () => {
    const round = createRound(true, "Attack");
    cy.mount(MatchProgressStep, { props: { round: round } });
    cy.get(".step").should("have.class", "won");
  })

  it('should render attack icon correctly', () => {
    const round = createRound(true, "Attack");
    cy.mount(MatchProgressStep, { props: { round: round } });
    cy.get(".step").find(".icon.attack").should("exist");
  })

  it('should render defense icon correctly', () => {
    const round = createRound(true, "Defense");
    cy.mount(MatchProgressStep, { props: { round: round } });
    cy.get(".step").find(".icon.defense").should("exist");
  })

  it('should have default state indicator', () => {
    const round = createRound(true, "Defense");
    cy.mount(MatchProgressStep, { props: { round: round } });
    cy.get(".status").should("exist");
    cy.get(".status > svg").should("exist");
  })

  it('should have defined state indicator', () => {
    const round = createRound(true, "Defense");
    round.status = "done";
    cy.mount(MatchProgressStep, { props: { round: round } });
    cy.get(".status").should("exist");
    cy.get(".status > svg").should("exist");
  })
})
