<script lang="ts">
	import { Button, Tile } from "carbon-components-svelte";
	import Minimize from "carbon-icons-svelte/lib/ChevronMini.svelte";
	import Close from "carbon-icons-svelte/lib/Close.svelte";
	import Unmaximize from "carbon-icons-svelte/lib/Minimize.svelte";
	import Maximize from "carbon-icons-svelte/lib/Maximize.svelte";

	import { GetWindowTitle } from "../wailsjs/go/main/App.js";

	let title: string = "";
	let isMaximized = false;
	window.runtime.WindowIsMaximised().then((v: boolean) => (isMaximized = v));

	GetWindowTitle().then((v) => (title = v));

	function toggleMaximise(): void {
		if (isMaximized) {
			window.runtime.WindowUnmaximise();
		} else {
			window.runtime.WindowMaximise();
		}
		isMaximized = !isMaximized;
	}
</script>

<Tile style="--wails-draggable:drag;">
	<div id="root">
		<div id="title">
			<span>{title}</span>
		</div>
		<Button
			on:click={window.runtime.WindowMinimise}
			icon={Minimize}
			iconDescription="Minimize"
			size="small"
			kind="tertiary"
			style="--wails-draggable:no-drag"
		/>
		<Button
			on:click={toggleMaximise}
			icon={isMaximized ? Unmaximize : Maximize}
			iconDescription={isMaximized ? "Unmaximize" : "Maximize"}
			size="small"
			kind="tertiary"
			style="--wails-draggable:no-drag"
		/>
		<Button
			on:click={window.runtime.Quit}
			icon={Close}
			iconDescription="Close"
			size="small"
			kind="danger-tertiary"
			style="--wails-draggable:no-drag"
		/>
	</div>
</Tile>

<style>
	#root {
		display: flex;
		flex-direction: row;
	}

	#title {
		flex-grow: 1;
		position: relative;
	}

	#title > span {
		position: absolute;
		top: 50%;
		left: 0;
		transform: translateY(-50%);
		font-size: 1.5em;
		font-weight: bold;

		cursor: default;
	}
</style>
