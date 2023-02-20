<script lang="ts">
	import {
		Button,
		ExpandableTile,
		InlineNotification,
		Link,
		Modal,
		SkeletonPlaceholder,
		SkeletonText,
		Tag,
		Tile,
		TooltipDefinition,
	} from "carbon-components-svelte";
	import {
		CloudDownload,
		Launch,
		WatsonHealthRotate_360,
	} from "carbon-icons-svelte";
	import { marked } from "marked";
	import LoadingOverlay from "../components/LoadingOverlay.svelte";
	import {
		GetAppInfo,
		GetEventNames,
		RequestLatestReleaseInfo,
		StartUpdate,
	} from "../../wailsjs/go/main/App";
	import type { app } from "../index";
	import { onMount } from "svelte";

	export let open = false;

	// stupid wrapper to add typing because Wails doesn't generate typing for GetAppInfo
	function appInfo(): Promise<app.AppInfo> {
		return GetAppInfo();
	}

	let latestReleaseInfo: app.ReleaseInfo = null;
	let latestReleaseInfoErr: string = null;
	let updateCheckCooldownFunc = null;
	const updateCheckCooldownMs = 60 * 1000;

	function checkForUpdate() {
		if (updateCheckCooldownFunc) {
			return;
		}
		updateCheckCooldownFunc = setTimeout(() => {
			updateCheckCooldownFunc = null;
		}, updateCheckCooldownMs);
		RequestLatestReleaseInfo();
	}

	let updateOverlayVisible = false;
	let updateComplete = false;
	let updateTask: string;
	let updateErr: string = null;

	function startUpdate(release?: app.ReleaseInfo) {
		if (updateComplete) {
			// update already complete, waiting for restart
			updateOverlayVisible = true;
			return;
		}
		if (release!.IsNewer) {
			StartUpdate()
				.then(() => {
					updateOverlayVisible = true;
					updateTask = "Preparing...";
					updateErr = null;
				})
				.catch((e) => (updateErr = e));
		}
	}

	function onUpdateProgress(p: app.UpdateProgress) {
		if (!p.Complete) {
			updateTask = p.Task;
		} else {
			updateComplete = true;
			updateTask = null;
		}
	}

	function onUpdateErr(err: string) {
		updateTask = null;
		updateErr = err;
	}

	function onLatestReleaseInfo(r: app.ReleaseInfo) {
		latestReleaseInfo = r;
		latestReleaseInfoErr = null;
	}

	function onLatestReleaseInfoErr(e: string) {
		latestReleaseInfoErr = e;
	}

	onMount(() => {
		GetEventNames().then((e: app.EventNames) => {
			window.runtime.EventsOn(e.UpdateProgress, onUpdateProgress);
			window.runtime.EventsOn(e.UpdateErr, onUpdateErr);
			window.runtime.EventsOn(e.LatestReleaseInfo, onLatestReleaseInfo);
			window.runtime.EventsOn(e.LatestReleaseInfoErr, onLatestReleaseInfoErr);
		});
	});
</script>

<Modal
	{open}
	passiveModal
	modalHeading="Application information"
	hasScrollingContent
	on:close={() => (open = false)}
>
	<LoadingOverlay
		bind:open={updateOverlayVisible}
		loadingDesc={updateTask}
		errorTitle={updateErr === null ? null : "Update failed"}
		errorDetail={updateErr}
		done={updateComplete}
	>
		<InlineNotification
			kind="success"
			title="Update downloaded and applied"
			subtitle="Please restart the application"
			on:close={(e) => {
				e.preventDefault();
				updateOverlayVisible = false;
				open = false;
			}}
		/>
	</LoadingOverlay>
	<Tile>
		{#await appInfo()}
			<SkeletonText heading />
			<SkeletonPlaceholder style="width: 100%; height: 20vh;" />
		{:then info}
			<!--HEADER-->
			<h2>{info.ProjectName}</h2>
			<Link
				icon={Launch}
				on:click={() => window.runtime.BrowserOpenURL(info.GithubURL)}
				>Github Repository</Link
			>
			<!--VERSION-->
			<div id="versions-container">
				<span id="lbl-current-version">Current version:</span>
				<span id="tag-current-version"
					><Tag>{info.Version} - {info.Commit}</Tag></span
				>

				<span id="lbl-latest-version">Latest version:</span>
				<span id="tag-latest-version">
					{#if latestReleaseInfo === null}
						<Tag skeleton />
					{:else if latestReleaseInfoErr === null}
						<Tag
							type={latestReleaseInfo.IsNewer ? "green" : "gray"}
							interactive={latestReleaseInfo.IsNewer}
							icon={latestReleaseInfo.IsNewer ? CloudDownload : undefined}
							on:click={() => startUpdate(latestReleaseInfo)}
						>
							{latestReleaseInfo.Version} - {latestReleaseInfo.Commitish}
						</Tag>
					{:else}
						<TooltipDefinition
							direction="top"
							align="start"
							tooltipText={latestReleaseInfoErr}
						>
							<Tag type="red" interactive>Error</Tag>
						</TooltipDefinition>
					{/if}
				</span>
				<div id="btn-check-updates">
					{#if updateCheckCooldownFunc === null}
						<Button
							icon={WatsonHealthRotate_360}
							size="field"
							on:click={checkForUpdate}>Check for Updates</Button
						>
					{:else}
						<TooltipDefinition tooltipText="On cooldown" direction="top">
							<Button icon={WatsonHealthRotate_360} size="field" disabled
								>Check for Updates</Button
							>
						</TooltipDefinition>
					{/if}
				</div>
			</div>
			{#if updateErr && !updateOverlayVisible}
				<InlineNotification title="Error" subtitle={updateErr} />
			{/if}

			<!--UPDATE-->
			{#if latestReleaseInfo === null}
				<SkeletonPlaceholder style="width: 100%; height: 8vh;" />
			{:else}
				<ExpandableTile light>
					<div slot="above">
						<span
							>Changelog for <strong>{latestReleaseInfo.Version}</strong></span
						>
					</div>
					<div slot="below">
						{#await marked.parse( latestReleaseInfo.Changelog, { gfm: true, async: true, headerIds: false } )}
							<SkeletonPlaceholder style="width: 100%; height: 10vh" />
						{:then md}
							{@html md}
						{/await}
						{@const publishedString = new Date(
							latestReleaseInfo.PublishedAt
						).toLocaleString()}
						<Tag type="outline">Published {publishedString}</Tag>
					</div>
				</ExpandableTile>
			{/if}
		{/await}
	</Tile>
</Modal>

<style>
	#versions-container {
		margin: 0.5rem 0;

		display: grid;
		grid-template-columns: repeat(2, fit-content(100%)) auto;
		grid-template-rows: repeat(2, 1fr);
		grid-template-areas:
			"lbl-current-version tag-current-version ."
			"lbl-latest-version  tag-latest-version  btn-check-updates";
		align-items: center;
		column-gap: 0.5rem;
	}

	#lbl-current-version {
		grid-area: lbl-current-version;
	}

	#lbl-latest-version {
		grid-area: lbl-latest-version;
	}

	#tag-current-version {
		grid-area: tag-current-version;
	}

	#tag-latest-version {
		grid-area: tag-latest-version;
	}

	#btn-check-updates {
		grid-area: btn-check-updates;
		justify-self: end;
	}
</style>
