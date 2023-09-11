<script lang="ts">
	import { InlineLoading, InlineNotification } from "carbon-components-svelte";
	import { afterUpdate, onMount } from "svelte";

	import {
		GetEventNames,
		IsConfigComplete,
		StartRoundWatcher,
		StopRoundWatcher,
	} from "./../../wailsjs/go/main/App";
	import type { EventNames } from "../app";
	import type { RoundInfo } from "../game";
	import type { InfluxEvent } from "../db";
	import type { MatchListAPI } from "./matchlist";
	import MatchItem from "./MatchItem.svelte";
	import { Round } from "./matchitem";

	let errorTitle = "Error";
	let error: string | null;

	export const matchListAPI: MatchListAPI = {
		async onConfigChanged() {
			try {
				if (roundWatcherRunning) {
					await StopRoundWatcher();
				}
				let configComplete = await IsConfigComplete();
				if (configComplete) {
					await StartRoundWatcher();
				}
			} catch (e) {
				errorTitle = "Round Watcher error:";
				error = e;
			}
		},
	};

	let roundWatcherRunning = false;
	let matchesContainer: HTMLElement;
	let matchInfos: Map<string, Round[]> = new Map();

	async function onNewRound(r: RoundInfo) {
		let matchInfo = matchInfos.get(r.MatchID);
		if (matchInfo) {
			matchInfos.set(r.MatchID, [...matchInfo, new Round(r)]);
		} else {
			matchInfos.set(r.MatchID, [new Round(r)]);
		}
		matchInfos = matchInfos;
	}

	function onRoundWatcherStarted() {
		roundWatcherRunning = true;
	}

	function onRoundWatcherStopped() {
		roundWatcherRunning = false;
	}

	function onRoundWatcherError(err: string | null) {
		errorTitle = "Error parsing round:";
		error = err;
	}

	function onRoundPush(e: InfluxEvent) {
		// find entry for match by ID
		let matchInfo = matchInfos.get(e.MatchID);
		if (matchInfo) {
			if (e.Err) {
				matchInfo[e.RoundIndex].setPushError(e.Err);
			} else {
				matchInfo[e.RoundIndex].pushStatus = "done";
			}

			matchInfos.set(e.MatchID, matchInfo);
			matchInfos = matchInfos;
		}
	}

	onMount(() => {
		GetEventNames().then((e: EventNames) => {
			window.runtime.EventsOn(e.NewRound, onNewRound);
			window.runtime.EventsOn(e.RoundWatcherStarted, onRoundWatcherStarted);
			window.runtime.EventsOn(e.RoundWatcherStopped, onRoundWatcherStopped);
			window.runtime.EventsOn(e.RoundWatcherError, onRoundWatcherError);
			window.runtime.EventsOn(e.RoundPush, onRoundPush);
		});
	});

	afterUpdate(async () => {
		if (matchesContainer) {
			let lastItem =
				matchesContainer.children[matchesContainer.children.length - 1];
			lastItem.scrollIntoView({ behavior: "smooth" });
		}
	});
</script>

{#if roundWatcherRunning}
	<div id="match-container" bind:this={matchesContainer}>
		{#each [...matchInfos] as [matchID, roundInfos] (matchID)}
			<MatchItem rounds={roundInfos} />
		{:else}
			<div class="placeholder-center">
				<div class="placeholder-center-content">
					<InlineLoading description="Waiting for matches..." />
				</div>
			</div>
		{/each}
	</div>
{:else}
	<div class="placeholder-center">
		<div class="placeholder-center-content">
			<InlineLoading description="Waiting for Round Watcher..." />
		</div>
	</div>
{/if}

<div id="error-container">
	{#if error}
		<InlineNotification
			kind="error"
			title={errorTitle}
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

	.placeholder-center {
		height: 100%;
		position: relative;
	}

	.placeholder-center-content {
		position: absolute;

		left: 50%;
		top: 50%;
		transform: translate(-50%, -50%);
	}
</style>
