import type { RoundInfo } from "../game";
import type { RoundPushStatus } from "./MatchProgressStep.svelte";

export class Round {
	data: RoundInfo;
	pushStatus: RoundPushStatus;
	pushError: string;

	constructor(data: RoundInfo) {
		this.data = data;
		this.pushStatus = "waiting";
		this.pushError = "";
	}

	public setPushError(e: any) {
		this.pushStatus = "error";
		if (typeof e === "string") {
			this.pushError = e;
		} else if ("Message" in e) {
			this.pushError = e.Message;
		} else {
			this.pushError = JSON.stringify(e, null, 0);
		}
		console.log(e);
	}
}
