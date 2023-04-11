<script>
	import MatchProgressStep from "../MatchProgressStep.svelte";

	export let Component;
	export let rounds;
</script>

<svelte:component this={Component}>
	{#each rounds as round}
		<MatchProgressStep {round} />
	{/each}
</svelte:component>
