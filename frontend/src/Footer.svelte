<script lang="ts">
	import { Button, TextInput, Tile } from "carbon-components-svelte";
	import Add from "carbon-icons-svelte/lib/Add.svelte";

	import { GetVersion } from "../wailsjs/go/main/App.js";

	let buildInfo: string = "";
	export let openSettings: () => any;

	GetVersion().then((bi) => {
		buildInfo = `${bi.Version} - ${bi.Commit}`;
	});
</script>

<div id="root" class="footer">
	<div id="left">
		<Button on:click={openSettings} icon={Add}>Greet</Button>
	</div>
	<pre>{buildInfo}</pre>
	<Button icon={Add} style="--wails-draggable:no-drag">Greet</Button>
</div>

<style>
	#root {
		width: 100%;

		display: flex;
		flex-direction: row;
		flex-wrap: nowrap;
		justify-content: flex-end;
		/*align-content: space-between;*/
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
