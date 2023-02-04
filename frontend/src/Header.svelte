<script lang="ts">
	import { Button, Tile } from "carbon-components-svelte";
	import Add from "carbon-icons-svelte/lib/Add.svelte";

	import { GetWindowTitle } from "../wailsjs/go/main/App.js";

	let title: string = "";

	GetWindowTitle().then((v) => (title = v));
</script>

<Tile style="--wails-draggable:drag">
	{title}
	<Button
		on:click={window.runtime.WindowMinimise}
		icon={Add}
		style="--wails-draggable:no-drag">Greet</Button
	>
	<Button
		on:click={window.runtime.WindowToggleMaximise}
		icon={Add}
		style="--wails-draggable:no-drag">Greet</Button
	>
	<Button
		on:click={window.runtime.Quit}
		icon={Add}
		style="--wails-draggable:no-drag">Greet</Button
	>
</Tile>
