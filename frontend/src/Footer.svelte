<script lang="ts">
	import IconButton from "@smui/icon-button";

	export let openSettings: () => any;
</script>

<div id="root" class="footer">
	<div id="left">
		<IconButton
			class="material-icons"
			style="--wails-draggable:no-drag"
			on:click={openSettings}>settings</IconButton
		>
	</div>
	<pre>v0.0.1 - 12345abc</pre>
	<IconButton class="material-icons-outlined">info</IconButton>
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
