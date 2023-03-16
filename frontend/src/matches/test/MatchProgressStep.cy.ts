import MatchProgressStep from "../MatchProgressStep.svelte";
import { createRoundInfo } from "./util";

describe('MatchProgressIndicator', () => {
  it('should have correct class based on win/loss', () => {
    const roundInfo = createRoundInfo(true, "ATTACK");
    cy.mount(MatchProgressStep, { props: { roundInfo: roundInfo } });
    cy.get(".step").should("have.class", "won");
  })

  it('should render attack icon correctly', () => {
    const roundInfo = createRoundInfo(true, "ATTACK");
    cy.mount(MatchProgressStep, { props: { roundInfo: roundInfo } });
    cy.get(".step").find(".icon.attack").should("exist");
  })

  it('should render defense icon correctly', () => {
    const roundInfo = createRoundInfo(true, "DEFENSE");
    console.log(roundInfo);
    cy.mount(MatchProgressStep, { props: { roundInfo: roundInfo } });
    cy.get(".step").find(".icon.defense").should("exist");
  })
})
