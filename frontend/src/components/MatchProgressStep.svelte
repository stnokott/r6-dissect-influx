<script lang="ts">
	import { Checkmark, Close } from "carbon-icons-svelte";

	export let won: boolean;
</script>

<div class="step {won ? 'won' : 'lost'}">
	{#if won}
		<Checkmark size={24} class="icon" />
	{:else}
		<Close size={24} class="icon" />
	{/if}
</div>

<style lang="scss">
	@use "@carbon/themes/scss/themes" as *;
	@use "@carbon/themes" with (
		$theme: $g90
	);

	$size: 34px;

	.step {
		/* border: 3px solid themes.$layer-accent-03; */
		border-radius: 100%;
		width: $size;
		height: $size;
		line-height: $size;
		text-align: center;
		font-family: sans-serif;
		font-size: 14px;
		position: relative;
		z-index: 1;

		display: flex;
		align-items: center;
	}

	.won {
		background-color: themes.$support-success;
	}

	.lost {
		background-color: themes.$layer-accent-03;
	}

	// use :global because Svelte doesn't recognise Carbon icons' SVG class
	.step > :global(.icon) {
		display: block;
		margin: auto;
	}
</style>
