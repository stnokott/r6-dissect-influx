import type { RoundInfo } from "../game";
import type { RoundStatus } from "./MatchProgressStep.svelte";

export class Round {
	data: RoundInfo;
	status: RoundStatus;

	constructor(data: RoundInfo) {
		this.data = data;
		this.status = "waiting";
	}
}
