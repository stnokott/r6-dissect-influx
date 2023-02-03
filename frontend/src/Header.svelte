<script lang="ts">
	import TopAppBar, { Row, Section, Title } from "@smui/top-app-bar";
	import IconButton from "@smui/icon-button";
</script>

<TopAppBar variant="fixed" dense color="primary" style="--wails-draggable:drag">
	<Row>
		<Section>
			<Title>Standard</Title>
		</Section>
		<Section align="end" toolbar>
			<IconButton
				class="material-icons-outlined"
				aria-label="Minimize"
				on:click={window.runtime.WindowMinimise}
				style="--wails-draggable:no-drag">minimize</IconButton
			>
			<IconButton
				class="material-icons-outlined"
				aria-label="Maximize"
				on:click={window.runtime.WindowToggleMaximise}
				style="--wails-draggable:no-drag">crop_square</IconButton
			>
			<IconButton
				class="material-icons-outlined"
				aria-label="Close"
				on:click={window.runtime.Quit}
				style="--wails-draggable:no-drag">close</IconButton
			>
		</Section>
	</Row>
</TopAppBar>

<style>
	/* Hide everything above this component. */
	:global(#smui-app),
	:global(body),
	:global(html) {
		display: block !important;
		height: auto !important;
		width: auto !important;
		position: static !important;
	}
</style>
