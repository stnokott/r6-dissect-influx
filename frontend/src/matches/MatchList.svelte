<script lang="ts">
	import { InlineNotification } from "carbon-components-svelte";
	import { onMount } from "svelte";
	import MatchItem from "./MatchItem.svelte";

	import {
		GetEventNames,
		StartRoundWatcher,
	} from "./../../wailsjs/go/main/App";
	import type { app } from "../index";
	import type { matches } from "./matches";

	let error: string;
	let matchInfos: Map<string, Array<matches.RoundInfo>> = new Map();

	async function onNewRound(r: matches.RoundInfo) {
		if (matchInfos.has(r.MatchID)) {
			matchInfos.set(r.MatchID, [...matchInfos.get(r.MatchID), r]);
		} else {
			matchInfos.set(r.MatchID, [r]);
		}
		matchInfos = matchInfos;
	}

	async function onRoundWatcherStopped() {}

	async function onRoundWatcherError(err: any) {
		error = err;
	}

	onMount(() => {
		GetEventNames().then((e: app.EventNames) => {
			window.runtime.EventsOn(e.NewRound, onNewRound);
			window.runtime.EventsOn(e.RoundWatcherStopped, onRoundWatcherStopped);
			window.runtime.EventsOn(e.RoundWatcherError, onRoundWatcherError);
			//TODO: only start if config ready
			StartRoundWatcher().catch((e) => console.log(e));
		});
	});
</script>

<div id="match-container">
	{#each [...matchInfos] as [matchID, roundInfos] (matchID)}
		<MatchItem {roundInfos} />
	{:else}
		No content
	{/each}
</div>

<div id="error-container">
	{#if error}
		<InlineNotification
			kind="error"
			title="Error"
			timeout={5000}
			subtitle={error}
			on:close={(e) => {
				e.preventDefault();
				error = null;
			}}
		/>
	{/if}
</div>

<style>
	#match-container {
		display: flex;
		flex-direction: column;
		gap: 0.5rem;
	}

	#error-container {
		position: fixed;
		right: 0;
		bottom: 0;

		margin: 0 2rem 2rem 0;
	}
</style>
