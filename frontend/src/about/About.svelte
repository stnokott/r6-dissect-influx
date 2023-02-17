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
	import { CloudDownload, Launch } from "carbon-icons-svelte";
	import { marked } from "marked";
	import LoadingOverlay from "../components/LoadingOverlay.svelte";
	import {
		GetAppInfo,
		GetEventNames,
		GetLatestRelease,
		StartUpdate,
	} from "../../wailsjs/go/main/App";
	import type { app } from "../index";
	import { onMount } from "svelte";

	export let open = false;

	const promReleaseInfo: Promise<app.ReleaseInfo> = GetLatestRelease();

	// stupid wrapper to add typing because Wails doesn't generate typing for GetAppInfo
	function appInfo(): Promise<app.AppInfo> {
		return GetAppInfo();
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
			// TODO: catch
			StartUpdate()
				.then(() => {
					updateOverlayVisible = true;
					updateTask = "Preparing...";
					updateErr = null;
				})
				.catch((e) => console.log(e));
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
		console.log(`err: ${err}`);
		updateTask = null;
		updateErr = err;
	}

	onMount(() => {
		GetEventNames().then((e: app.EventNames) => {
			window.runtime.EventsOn(e.UpdateProgress, onUpdateProgress);
			window.runtime.EventsOn(e.UpdateErr, onUpdateErr);
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
			<SkeletonText />
			<SkeletonText heading />
			<SkeletonText />
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
				<span>Current version:</span>
				<Tag>{info.Version} - {info.Commit}</Tag>
				<span>Latest version:</span>
				{#await promReleaseInfo}
					<Tag skeleton />
				{:then release}
					<Tag
						type={release.IsNewer ? "green" : "gray"}
						interactive={release.IsNewer}
						icon={release.IsNewer ? CloudDownload : undefined}
						on:click={() => startUpdate(release)}
					>
						{release.Version} - {release.Commitish}
					</Tag>
				{:catch err}
					<TooltipDefinition direction="top" align="start" tooltipText={err}>
						<Tag type="red" interactive>Error</Tag>
					</TooltipDefinition>
				{/await}
			</div>
			<!--UPDATE-->
			{#await promReleaseInfo}
				<SkeletonPlaceholder style="width: 100%" />
			{:then release}
				<ExpandableTile light>
					<div slot="above">
						<span>Changelog for <strong>{release.Version}</strong></span>
					</div>
					<div slot="below">
						{#await marked.parse( release.Changelog, { gfm: true, async: true, headerIds: false } )}
							<SkeletonPlaceholder style="width: 100%; height: 10vh" />
						{:then md}
							{@html md}
						{/await}
						{@const publishedString = new Date(
							release.PublishedAt
						).toLocaleString()}
						<Tag type="outline">Published {publishedString}</Tag>
					</div>
				</ExpandableTile>
			{/await}
		{/await}
	</Tile>
</Modal>

<style>
	#versions-container {
		margin-top: 0.5rem;

		display: grid;
		grid-template-columns: repeat(2, fit-content(100%));
		align-items: center;
		column-gap: 0.5rem;
	}
</style>
