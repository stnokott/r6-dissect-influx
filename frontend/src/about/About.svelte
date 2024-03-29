<script lang="ts">
	import {
		Button,
		ExpandableTile,
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
	import type { AppInfo, EventNames, ReleaseInfo } from "../app";
	import { onMount } from "svelte";

	export let open = false;

	// stupid wrapper to add typing because Wails doesn't generate typing for GetAppInfo
	function appInfo(): Promise<AppInfo> {
		return GetAppInfo();
	}

	let latestReleaseInfo: ReleaseInfo | null = null;
	let latestReleaseInfoErr = "";
	let updateCheckCooldownFunc: NodeJS.Timeout | null = null;
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
	let updateTask: string;
	let updateErr = "";

	function startUpdate(release: ReleaseInfo | null) {
		if (release !== null && release.IsNewer) {
			StartUpdate()
				.then(() => {
					updateOverlayVisible = true;
					updateTask = "Preparing...";
					updateErr = "";
				})
				.catch((e) => (updateErr = e));
		}
	}

	function onUpdateProgress(description: string) {
		updateTask = description;
	}

	function onUpdateErr(err: string) {
		updateTask = "";
		updateErr = err;
	}

	function onLatestReleaseInfo(r: ReleaseInfo) {
		latestReleaseInfo = r;
		latestReleaseInfoErr = "";
		updateErr = "";
	}

	function onLatestReleaseInfoErr(e: string) {
		latestReleaseInfoErr = e;
	}

	onMount(() => {
		GetEventNames().then((e: EventNames) => {
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
		errorTitle={updateErr === "" ? null : "Update failed"}
		errorDetail={updateErr}
	/>
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
					{:else if latestReleaseInfoErr === ""}
						<Tag type={latestReleaseInfo.IsNewer ? "green" : "gray"}>
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
			</div>

			<div id="update-buttons-container">
				{#if latestReleaseInfo && latestReleaseInfo.IsNewer}
					<Button
						icon={CloudDownload}
						size="small"
						disabled={updateOverlayVisible}
						on:click={() => startUpdate(latestReleaseInfo)}
					>
						Apply Update
					</Button>
				{/if}
				{#if updateCheckCooldownFunc === null}
					<Button
						icon={WatsonHealthRotate_360}
						size="small"
						kind="tertiary"
						disabled={updateOverlayVisible}
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
		grid-template-columns: repeat(2, fit-content(100%));
		grid-template-rows: repeat(2, 1fr);
		grid-template-areas:
			"lbl-current-version tag-current-version"
			"lbl-latest-version  tag-latest-version "
			"buttons-container   buttons-container  ";
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

	#update-buttons-container {
		margin: 1rem 0;
	}
</style>
