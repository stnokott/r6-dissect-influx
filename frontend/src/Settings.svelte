<script lang="ts">
	import {
		Modal,
		Button,
		ButtonSet,
		Grid,
		Row,
		Column,
		TextInput,
		NumberInput,
	} from "carbon-components-svelte";

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
>
	<Grid>
		<Row>
			<Column>
				<TextInput
					bind:value={gameDir}
					labelText="Game directory"
					helperText="Folder should contain RainbowSix.exe"
				/>
			</Column>
			<Column>
				<ButtonSet>
					<Button>Open</Button>
					<Button kind="secondary">Autodetect</Button>
				</ButtonSet>
			</Column>
		</Row>
		<Row>
			<Column>
				<TextInput
					bind:value={influxDBHost}
					labelText="InfluxDB Host"
					helperText="IP or hostname, without http(s)"
				/>
			</Column>
			<Column>
				<NumberInput bind:value={influxDBPort} label="Port" />
			</Column>
		</Row>
	</Grid>
</Modal>
