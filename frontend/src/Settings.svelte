<script lang="ts">
	import {
		Modal,
		Button,
		Form,
		Grid,
		Tile,
		Row,
		Column,
		TextInput,
		NumberInput,
	} from "carbon-components-svelte";
	import Folder from "carbon-icons-svelte/lib/Folder.svelte";

	export let open = false;

	let gameDir: string = "";
	let influxDBHost: string = "";
	let influxDBPort: number = 8086;
	let influxDBOrg: string = "";
	let influxDBBucket: string = "";
	let influxDBToken: string = "";
</script>

<Modal
	bind:open
	modalHeading="Settings"
	primaryButtonText="Save"
	secondaryButtonText="Cancel"
	hasForm={true}
	hasScrollingContent={true}
	on:click:button--secondary={() => (open = false)}
>
	<Form>
		<Tile light style="margin-bottom: 1rem;">
			<Grid narrow padding>
				<Row><h5>Game</h5></Row>
				<Row>
					<Column>
						<TextInput
							bind:value={gameDir}
							labelText="Directory"
							helperText="Folder should contain RainbowSix.exe"
							required
						/>
					</Column>
					<Column>
						<div id="game-dir-buttons">
							<Button icon={Folder} size="field" iconDescription="Open" />
							<Button kind="secondary" size="field">Autodetect</Button>
						</div>
					</Column>
				</Row>
			</Grid>
		</Tile>
		<Tile light>
			<Grid narrow padding>
				<Row><h5>InfluxDB</h5></Row>
				<Row>
					<Column>
						<TextInput
							bind:value={influxDBHost}
							labelText="Host"
							helperText="IP or hostname, without http(s)"
							required
						/>
					</Column>
					<Column>
						<NumberInput bind:value={influxDBPort} label="Port" required />
					</Column>
				</Row>
				<Row>
					<Column>
						<TextInput
							bind:value={influxDBOrg}
							labelText="Organization"
							required
						/>
					</Column>
					<Column>
						<TextInput
							bind:value={influxDBBucket}
							labelText="Bucket"
							required
						/>
					</Column>
				</Row>
				<Row>
					<TextInput bind:value={influxDBToken} labelText="Token" required />
				</Row>
			</Grid>
		</Tile>
	</Form>
</Modal>

<style>
	#game-dir-buttons {
		/* accounting for input label */
		margin-top: calc(1rem + 0.5rem);
	}
</style>
