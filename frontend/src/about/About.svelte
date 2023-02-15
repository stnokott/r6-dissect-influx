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

	let updateInProgress = false;
	let updateTask: string;
	let updateErr: string;

	function startUpdate(release: app.ReleaseInfo) {
		if (release.IsNewer) {
			// TODO: catch
			StartUpdate()
				.then(() => (updateInProgress = true))
				.catch((e) => console.log(e));
		}
	}

	function onUpdateProgress(p: app.UpdateProgress) {
		console.log("progress");
		console.log(p);
		if (!p.Complete) {
			updateTask = p.Task;
		} else {
			// TODO: popup
			updateInProgress = false;
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
		bind:open={updateInProgress}
		loadingDesc={updateTask}
		errorTitle="Update failed"
		errorDetail={updateErr}
	/>
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
