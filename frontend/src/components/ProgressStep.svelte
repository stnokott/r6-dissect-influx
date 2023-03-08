<script lang="ts">
	import { Checkmark, Close } from "carbon-icons-svelte";

	export let won: boolean;
</script>

<li class="step {won ? 'won' : 'lost'}">
	{#if won}
		<Checkmark size={16} />
	{:else}
		<Close size={16} />
	{/if}
</li>

<style lang="scss">
	@use "@carbon/themes/scss/themes" as *;
	@use "@carbon/themes" with (
		$theme: $g90
	);

	$size: 35px;

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
	}

	.won {
		background-color: themes.$support-success;
	}

	.lost {
		background-color: themes.$layer-accent-03;
	}
</style>
