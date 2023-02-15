<script lang="ts">
	import {
		Button,
		Tag,
		TagSkeleton,
		TooltipDefinition,
	} from "carbon-components-svelte";
	import { createEventDispatcher } from "svelte";
	import Information from "carbon-icons-svelte/lib/Information.svelte";
	import Settings from "carbon-icons-svelte/lib/Settings.svelte";
	import { onMount } from "svelte";

	import { GetAppInfo } from "../../wailsjs/go/main/App.js";
	import type { app, db } from "../index";
	import { Cloud, CloudOffline } from "carbon-icons-svelte";
	import About from "../about/About.svelte";

	const dispatch = createEventDispatcher<{ openSettings: void }>();

	let aboutOpen = false;

	let buildInfo: string = "";

	export let promConnectionDetails: Promise<db.ConnectionDetails> = null;

	onMount(() => {
		GetAppInfo().then((bi: app.AppInfo) => {
			buildInfo = `${bi.Version} - ${bi.Commit}`;
		});
	});
</script>

<div id="root" class="footer">
	<div id="left">
		<Button
			on:click={() => dispatch("openSettings")}
			icon={Settings}
			iconDescription="Settings"
			tooltipPosition="right"
			size="field"
			kind="secondary"
		/>
		<div id="connection-details">
			{#if promConnectionDetails}
				{#await promConnectionDetails}
					<TagSkeleton size="sm" />
				{:then connectionDetails}
					<Tag size="sm" icon={Cloud}
						>{connectionDetails.Name} - {connectionDetails.Version} - {connectionDetails.Commit}</Tag
					>
				{:catch err}
					<TooltipDefinition direction="top" align="start" tooltipText={err}>
						<Tag size="sm" icon={CloudOffline} type="red" interactive
							>Connection error</Tag
						>
					</TooltipDefinition>
				{/await}
			{:else}
				<Tag size="sm" icon={CloudOffline}>Not connected</Tag>
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
		on:click={() => (aboutOpen = true)}
	/>
</div>
<About bind:open={aboutOpen} />

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
