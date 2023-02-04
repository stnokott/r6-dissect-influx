<script lang="ts">
	import { Button, TextInput, Tile } from "carbon-components-svelte";
	import Add from "carbon-icons-svelte/lib/Add.svelte";

	let resultText: string = "Please enter your name below 👇";
	let name: string = "";

	function greet(): void {
		resultText = name;
	}
</script>

<Tile>
	<TextInput bind:value={name} labelText="Name" helperText="Name to greet" />
	<Button on:click={greet} icon={Add}>Greet</Button>
	{resultText}
</Tile>
