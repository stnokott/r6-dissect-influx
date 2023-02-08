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
		ToastNotification,
	} from "carbon-components-svelte";
	import Folder from "carbon-icons-svelte/lib/Folder.svelte";
	import { onMount } from "svelte";
	import {
		GetConfig,
		SaveConfig,
		AutodetectGameDir,
		ValidateGameDir,
		ValidateInfluxHost,
		ValidateInfluxPort,
		ValidateInfluxOrg,
		ValidateInfluxBucket,
		ValidateInfluxToken,
	} from "../../wailsjs/go/config/API";
	import { config } from "../../wailsjs/go/models";

	export let open = false;

	let errorTitle: string;
	let errorDetail: string;

	let gameDir: string = "";
	let influxHost: string = "";
	let influxPort: number = 8086;
	let influxOrg: string = "";
	let influxBucket: string = "";
	let influxToken: string = "";

	let gameDirValidationErr: string;
	let influxHostValidationErr: string;
	let influxPortValidationErr: string;
	let influxOrgValidationErr: string;
	let influxBucketValidationErr: string;
	let influxTokenValidationErr: string;

	$: formInvalid =
		gameDirValidationErr !== null ||
		influxHostValidationErr !== null ||
		influxPortValidationErr !== null ||
		influxOrgValidationErr !== null ||
		influxBucketValidationErr !== null ||
		influxTokenValidationErr !== null;

	let autodetectRunning = false;
	let autodetectError: string = null;

	async function autodetect(): Promise<void> {
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
		p.then(() => validationErrSetter(null)).catch((e) =>
			validationErrSetter(e)
		);
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

	async function writeConfig() {
		let cfg = new config.ConfigJSON({
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
		await SaveConfig(cfg);
	}

	function onOpen() {
		loadConfig().catch((e) => {
			errorTitle = "Could not load config";
			errorDetail = e;
		});
	}

	let writingConfig = false;

	function onConfirm() {
		writingConfig = true;
		writeConfig()
			.then(() => (open = false))
			.catch((e) => {
				errorTitle = "Could not write config";
				errorDetail = e;
			})
			.finally(() => (writingConfig = false));
	}
</script>

<Modal
	bind:open
	modalHeading="Settings"
	primaryButtonText="Save"
	primaryButtonDisabled={formInvalid || writingConfig}
	secondaryButtonText="Cancel"
	preventCloseOnClickOutside
	hasForm
	hasScrollingContent
	on:open={onOpen}
	on:click:button--primary={onConfirm}
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
							invalid={gameDirValidationErr !== null}
							invalidText={gameDirValidationErr}
							labelText="Directory"
							required
						/>
					</Column>
					<Column>
						<div id="game-dir-buttons">
							<Button icon={Folder} size="field" iconDescription="Open" />
							<Button
								on:click={autodetect}
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
							bind:value={influxHost}
							invalid={influxHostValidationErr !== null}
							invalidText={influxHostValidationErr}
							labelText="Host"
							helperText="IP or hostname, without http(s)"
							required
						/>
					</Column>
					<Column>
						<NumberInput
							bind:value={influxPort}
							invalid={influxPortValidationErr !== null}
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
							invalid={influxOrgValidationErr !== null}
							invalidText={influxOrgValidationErr}
							labelText="Organization"
							required
						/>
					</Column>
					<Column>
						<TextInput
							bind:value={influxBucket}
							invalid={influxBucketValidationErr !== null}
							invalidText={influxBucketValidationErr}
							labelText="Bucket"
							required
						/>
					</Column>
				</Row>
				<Row>
					<PasswordInput
						bind:value={influxToken}
						invalid={influxTokenValidationErr !== null}
						invalidText={influxTokenValidationErr}
						labelText="Token"
						required
					/>
				</Row>
			</Grid>
		</Tile>
	</Form>
	{#if errorTitle}
		<ToastNotification>
			kind="error" title="Error" subtitle={errorTitle}
			caption={errorDetail}
			timeout={5000}
			on:close={(e) => {
				e.preventDefault();
				errorTitle = null;
			}}
		</ToastNotification>
	{/if}
</Modal>

<style>
	#game-dir-buttons {
		/* accounting for input label */
		margin-top: calc(1rem + 0.5rem);
	}
</style>
