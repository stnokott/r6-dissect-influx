<script lang="ts">
	import { InlineLoading, InlineNotification } from "carbon-components-svelte";

	export let open: boolean;
	export let loadingDesc: string;

	// will show slotted component if done == true
	export let done = false;

	// will show loader with loadingDesc if errorTitle is null
	export let errorTitle: string | null = null;
	export let errorDetail: string | null = null;
</script>

<div id="root" style:display={open ? "flex" : "none"}>
	<div>
		{#if errorTitle === null}
			{#if !done}
				<div data-cy="loader">
					<InlineLoading description={loadingDesc} />
				</div>
			{:else}
				<slot />
			{/if}
		{:else}
			<div data-cy="error">
				<InlineNotification
					kind="error"
					title={errorTitle}
					subtitle={errorDetail ? errorDetail : ""}
					on:close={(e) => {
						e.preventDefault();
						open = false;
					}}
				/>
			</div>
		{/if}
	</div>
</div>

<style>
	#root {
		position: fixed;
		z-index: 1000;
		top: 0;
		left: 0;
		width: 100%;
		height: 100%;
		background-color: rgba(0, 0, 0, 0.666);

		display: flex;
		align-items: center;
		justify-content: center;
	}
</style>
