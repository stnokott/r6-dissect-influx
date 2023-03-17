<script lang="ts">
	export let fullWidth = false;
</script>

<div
	id="progress"
	style:width={fullWidth ? "100%" : "auto"}
	style={$$props.style}
>
	<ul id="progress-bar">
		<slot />
	</ul>
</div>

<style lang="scss">
	@use "@carbon/themes/scss/themes" as *;
	@use "@carbon/themes" with (
		$theme: $g90
	);

	#progress {
		position: relative;
	}

	#progress-bar {
		margin: 0;
		padding: 0;
		list-style: none;
		display: flex;
		justify-content: space-between;
	}

	#progress-bar::before {
		content: "";
		background-color: themes.$layer-accent-03;
		position: absolute;
		top: 50%;
		left: 0;
		transform: translateY(-50%);
		height: 6px;
		width: 100%;
		z-index: -1;
	}
</style>
