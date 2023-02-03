<script lang="ts">
	import Button from "@smui/button";
	import { Label } from "@smui/common";
	import Textfield from "@smui/textfield";
	import HelperText from "@smui/textfield/helper-text";
	import Card, { Actions, Content } from "@smui/card";

	import { Greet } from "../wailsjs/go/main/App.js";

	let resultText: string = "Please enter your name below 👇";
	let name: string = "";

	function greet(): void {
		Greet(name).then((result) => (resultText = result));
	}
</script>

<Card>
	<Content>
		<Textfield variant="filled" bind:value={name} label="Name">
			<HelperText slot="helper">Name to greet</HelperText>
		</Textfield>
		<pre>{resultText}</pre>
	</Content>
	<Actions fullBleed>
		<Button on:click={greet}>
			<Label>Greet</Label>
			<i class="material-icons-outlined" aria-hidden="true">arrow_forward</i>
		</Button>
	</Actions>
</Card>
