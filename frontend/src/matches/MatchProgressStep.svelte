<script lang="ts" context="module">
	const roundStatusList = ["waiting", "done", "error"] as const;
	export type RoundPushStatus = (typeof roundStatusList)[number];
</script>

<script lang="ts">
	import {
		CheckmarkFilled,
		PendingFilled,
		WarningFilled,
	} from "carbon-icons-svelte";
	import type { Round } from "./matchitem";
	import Attack from "./icons/Attack.svelte";
	import Defense from "./icons/Defense.svelte";
	import { Tooltip } from "carbon-components-svelte";

	export let round: Round;
</script>

<div class="step {round.data.Won ? 'won' : 'lost'}">
	{#if round.data.Teams[round.data.TeamIndex].Role === "Attack"}
		<Attack size={24} class="icon attack" />
	{:else}
		<Defense size={24} class="icon defense" />
	{/if}
	<div class="status">
		{#if round.pushStatus === "waiting"}
			<PendingFilled size={16} />
		{:else if round.pushStatus === "done"}
			<CheckmarkFilled size={16} class="success" />
		{:else}
			<Tooltip>
				<WarningFilled slot="icon" size={16} class="error" />
				<p>{round.pushError}</p>
			</Tooltip>
		{/if}
	</div>
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

	// use :global because Svelte doesn't recognise Carbon icons' SVG class (.icon)
	.step > :global(.icon) {
		display: block;
		margin: auto;
	}

	.status {
		position: absolute;
		bottom: 35%;
		right: -15%;
		transform: translateY(50%);
		height: 16px;
		width: 16px;
	}

	.status > :global(:not(.error)) {
		fill: themes.$text-on-color;
		stroke: themes.$text-inverse;
	}

	.status :global(.error) {
		fill: themes.$support-error-inverse;
		stroke: themes.$text-on-color;
	}

	.status > :global(.success) {
		fill: themes.$text-on-color;
		stroke: themes.$support-success;
	}

	.status :global(.bx--tooltip__trigger) {
		margin: 0;
	}
</style>
