<script lang="ts">
	import { Button, Tag } from "carbon-components-svelte";
	import Information from "carbon-icons-svelte/lib/Information.svelte";
	import Settings from "carbon-icons-svelte/lib/Settings.svelte";
	import ConnectionSignal from "carbon-icons-svelte/lib/ConnectionSignal.svelte";
	import ConnectionSignalOff from "carbon-icons-svelte/lib/ConnectionSignalOff.svelte";
	import { onMount } from "svelte";

	import { GetVersion } from "../wailsjs/go/main/App.js";
	import type { db } from "./settings/types";

	let buildInfo: string = "";
	export let openSettings: () => any;
	export let connectionDetails: db.ConnectionDetails = null;

	onMount(() => {
		GetVersion().then((bi) => {
			buildInfo = `${bi.Version} - ${bi.Commit}`;
		});
	});
</script>

<div id="root" class="footer">
	<div id="left">
		<Button
			on:click={openSettings}
			icon={Settings}
			iconDescription="Settings"
			tooltipPosition="right"
			size="field"
			kind="secondary"
		/>
		<div id="connection-details">
			{#if connectionDetails}
				<ConnectionSignal size={16} />
				<Tag size="sm">{connectionDetails.Name}</Tag>
				<Tag size="sm"
					>{connectionDetails.Version} - {connectionDetails.Commit}</Tag
				>
			{:else}
				<ConnectionSignalOff size={16} />
				<Tag size="sm">Not connected</Tag>
			{/if}
		</div>
	</div>
	<pre>{buildInfo}</pre>
	<Button
		icon={Information}
		iconDescription="Information"
		tooltipPosition="left"
		size="field"
		kind="ghost"
	/>
</div>

<style>
	#root {
		width: 100%;

		display: flex;
		flex-direction: row;
		flex-wrap: nowrap;
		justify-content: flex-end;
		gap: 0.7em;
		align-items: center;
	}

	#left {
		flex-grow: 1;
		color: var(--surface);
		display: flex;
	}

	#connection-details {
		padding-left: 1rem;

		display: inline-flex;
		align-items: center;
		gap: 0.3rem;
	}
</style>
