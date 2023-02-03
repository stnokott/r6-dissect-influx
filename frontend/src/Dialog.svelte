<script lang="ts">
	import Button, { Label } from "@smui/button";
	import Dialog, { Title, Content, Actions } from "@smui/dialog";

	export let open = false;
</script>

<Dialog
	bind:open
	aria-labelledby="simple-title"
	aria-describedby="simple-content"
>
	<Title>Dialog Title</Title>
	<Content>Super awesome dialog body text?</Content>
	<Actions>
		<Button>
			<Label>No</Label>
		</Button>
		<Button>
			<Label>Yes</Label>
		</Button>
	</Actions>
</Dialog>
