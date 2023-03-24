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
		PasswordInput,
		NumberInput,
		InlineNotification,
	} from "carbon-components-svelte";
	import { createEventDispatcher } from "svelte";
	import Folder from "carbon-icons-svelte/lib/Folder.svelte";
	import {
		GetConfig,
		SaveAndValidateConfig,
		OpenGameDirDialog,
		AutodetectGameDir,
		ValidateGameDir,
		ValidateInfluxHost,
		ValidateInfluxPort,
		ValidateInfluxOrg,
		ValidateInfluxBucket,
		ValidateInfluxToken,
	} from "../../wailsjs/go/main/App";
	import { config } from "../../wailsjs/go/models";
	import type { ConnectionDetails } from "../db";
	import LoadingOverlay from "../components/LoadingOverlay.svelte";

	export let open = false;

	const dispatch = createEventDispatcher<{
		connected: ConnectionDetails;
		changed: void;
	}>();

	let errorTitle: string | null;
	let errorDetail: string | null;

	let gameDir = "";
	let influxHost = "";
	let influxPort = 8086;
	let influxOrg = "";
	let influxBucket = "";
	let influxToken = "";

	let gameDirValidationErr = "";
	let influxHostValidationErr = "";
	let influxPortValidationErr = "";
	let influxOrgValidationErr = "";
	let influxBucketValidationErr = "";
	let influxTokenValidationErr = "";

	$: configInvalid =
		gameDirValidationErr !== "" ||
		influxHostValidationErr !== "" ||
		influxPortValidationErr !== "" ||
		influxOrgValidationErr !== "" ||
		influxBucketValidationErr !== "" ||
		influxTokenValidationErr !== "";

	function openGameDirDialog(): void {
		OpenGameDirDialog().then((d) => {
			if (d !== "") {
				gameDir = d;
			}
		});
	}

	let autodetectRunning = false;
	let autodetectError: string | null = null;

	async function autodetectGameDir(): Promise<void> {
		autodetectRunning = true;
		try {
			gameDir = await AutodetectGameDir();
			autodetectError = null;
		} catch (e) {
			autodetectError = e;
		} finally {
			autodetectRunning = false;
		}
	}

	function handleValidationPromise(
		p: Promise<void>,
		validationErrSetter: (e: string) => void
	): void {
		p.then(() => validationErrSetter("")).catch((e) => validationErrSetter(e));
	}

	$: handleValidationPromise(
		ValidateGameDir(gameDir),
		(e) => (gameDirValidationErr = e)
	);

	$: handleValidationPromise(
		ValidateInfluxHost(influxHost),
		(e) => (influxHostValidationErr = e)
	);

	$: handleValidationPromise(
		ValidateInfluxPort(influxPort.toString()),
		(e) => (influxPortValidationErr = e)
	);

	$: handleValidationPromise(
		ValidateInfluxOrg(influxOrg),
		(e) => (influxOrgValidationErr = e)
	);

	$: handleValidationPromise(
		ValidateInfluxBucket(influxBucket),
		(e) => (influxBucketValidationErr = e)
	);

	$: handleValidationPromise(
		ValidateInfluxToken(influxToken),
		(e) => (influxTokenValidationErr = e)
	);

	async function loadConfig() {
		let cfg = await GetConfig();
		gameDir = cfg.game.install_dir;
		influxHost = cfg.influx_db.host;
		influxPort = cfg.influx_db.port;
		influxOrg = cfg.influx_db.org;
		influxBucket = cfg.influx_db.bucket;
		influxToken = cfg.influx_db.token;
	}

	async function saveConfig() {
		let cfg = new config.Config({
			game: {
				install_dir: gameDir,
			},
			influx_db: {
				host: influxHost,
				port: influxPort,
				org: influxOrg,
				bucket: influxBucket,
				token: influxToken,
			},
		});
		let connDetails = await SaveAndValidateConfig(cfg);

		dispatch("connected", connDetails);
		dispatch("changed");
	}

	let loadingConfig = false;

	function onOpen() {
		loadingConfig = true;
		loadConfig()
			.catch((e) => {
				errorTitle = "Could not load config:";
				errorDetail = e;
			})
			.finally(() => (loadingConfig = false));
	}

	let loadingOverlayVisible = false;
	let validatingConfig = false;

	function onConfirm() {
		validatingConfig = true;
		loadingOverlayVisible = true;
		errorTitle = null;
		saveConfig()
			.then(() => {
				open = false;
				loadingOverlayVisible = false;
			})
			.catch((e) => {
				errorTitle = "Could not save config:";
				errorDetail = e;
			})
			.finally(() => {
				validatingConfig = false;
			});
	}

	let loadingDesc: string;
	$: {
		if (autodetectRunning) {
			loadingDesc = "Finding game folder...";
		} else if (loadingConfig) {
			loadingDesc = "Loading configuration...";
		} else if (validatingConfig) {
			loadingDesc = "Validating configuration...";
		} else {
			loadingDesc = "";
		}
	}
</script>

<Modal
	bind:open
	modalHeading="Settings"
	primaryButtonText="Save"
	primaryButtonDisabled={configInvalid ||
		validatingConfig ||
		loadingOverlayVisible}
	secondaryButtonText="Cancel"
	preventCloseOnClickOutside
	hasForm
	hasScrollingContent
	selectorPrimaryFocus="#input-influx-host"
	on:open={onOpen}
	on:click:button--primary={onConfirm}
	on:click:button--secondary={() => (open = false)}
>
	<LoadingOverlay
		bind:open={loadingOverlayVisible}
		{loadingDesc}
		{errorTitle}
		{errorDetail}
	/>
	<Form>
		<Tile light style="margin-bottom: 1rem;">
			<Grid narrow padding>
				<Row><h5>Game</h5></Row>
				<Row>
					<Column>
						<TextInput
							bind:value={gameDir}
							invalid={gameDirValidationErr !== ""}
							invalidText={gameDirValidationErr}
							labelText="Directory"
							required
						/>
					</Column>
					<Column>
						<div id="game-dir-buttons">
							<Button
								on:click={openGameDirDialog}
								icon={Folder}
								size="field"
								iconDescription="Open"
							/>
							<Button
								on:click={autodetectGameDir}
								disabled={autodetectRunning}
								kind="secondary"
								size="field">Autodetect</Button
							>
						</div>
					</Column>
				</Row>
				<Row>
					{#if autodetectError}
						<InlineNotification
							kind="warning"
							lowContrast
							on:close={(e) => {
								e.preventDefault();
								autodetectError = null;
							}}
						>
							<strong slot="title">Autodetection failed: </strong>
							<span slot="subtitle">{autodetectError}</span>
						</InlineNotification>
					{/if}
				</Row>
			</Grid>
		</Tile>
		<Tile light>
			<Grid narrow padding>
				<Row><h5>InfluxDB</h5></Row>
				<Row>
					<Column>
						<TextInput
							id="input-influx-host"
							bind:value={influxHost}
							invalid={influxHostValidationErr !== ""}
							invalidText={influxHostValidationErr}
							labelText="Host"
							helperText="IP or hostname, without http(s)"
							required
						/>
					</Column>
					<Column>
						<NumberInput
							bind:value={influxPort}
							invalid={influxPortValidationErr !== ""}
							invalidText={influxPortValidationErr}
							label="Port"
							required
						/>
					</Column>
				</Row>
				<Row>
					<Column>
						<TextInput
							bind:value={influxOrg}
							invalid={influxOrgValidationErr !== ""}
							invalidText={influxOrgValidationErr}
							labelText="Organization"
							required
						/>
					</Column>
					<Column>
						<TextInput
							bind:value={influxBucket}
							invalid={influxBucketValidationErr !== ""}
							invalidText={influxBucketValidationErr}
							labelText="Bucket"
							required
						/>
					</Column>
				</Row>
				<Row>
					<PasswordInput
						bind:value={influxToken}
						invalid={influxTokenValidationErr !== ""}
						invalidText={influxTokenValidationErr}
						labelText="Token"
						required
					/>
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
