<script lang="ts">
	import type { RoundInfo } from "../game";
	import Attack from "./icons/Attack.svelte";
	import Defense from "./icons/Defense.svelte";

	export let roundInfo: RoundInfo;
</script>

<div class="step {roundInfo.Won ? 'won' : 'lost'}">
	{#if roundInfo.Teams[roundInfo.TeamIndex].Role === "ATTACK"}
		<Attack size={24} class="icon attack" />
	{:else}
		<Defense size={24} class="icon defense" />
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
