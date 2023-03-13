<script>
	import MatchProgressStep from "../MatchProgressStep.svelte";

	export let Component;
	export let roundInfos;
</script>

<svelte:component this={Component}>
	{#each roundInfos as roundInfo}
		<MatchProgressStep {roundInfo} />
	{/each}
</svelte:component>
