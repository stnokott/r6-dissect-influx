<script lang="ts">
	import { Button, TextInput, Tile } from "carbon-components-svelte";
	import Information from "carbon-icons-svelte/lib/Information.svelte";
	import Settings from "carbon-icons-svelte/lib/Settings.svelte";
	import { onMount } from "svelte";

	import { GetVersion } from "../wailsjs/go/main/App.js";

	let buildInfo: string = "";
	export let openSettings: () => any;

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

	#root > * {
		display: inline-block;
	}

	#left {
		flex-grow: 1;

		color: var(--surface);
	}
</style>
