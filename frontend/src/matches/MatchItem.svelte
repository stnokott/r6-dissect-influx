<script lang="ts">
	import {
		ImageLoader,
		SkeletonPlaceholder,
		Tag,
		Tile,
	} from "carbon-components-svelte";
	import type { TagProps } from "carbon-components-svelte/types/Tag/Tag.svelte";
	import MatchProgressIndicator from "../components/MatchProgressIndicator.svelte";
	import MatchProgressStep from "../components/MatchProgressStep.svelte";
	import type { matches } from "./matches";

	export let roundInfos: Array<matches.RoundInfo>;

	type TagTypes<T extends string> = {
		[key in T]: TagProps["type"];
	};

	const matchTypeColors: TagTypes<matches.RoundInfo["MatchType"]> = {
		CASUAL: "blue",
		UNRANKED: "purple",
		RANKED: "magenta",
	};

	const gameModeColors: TagTypes<matches.RoundInfo["GameMode"]> = {
		BOMB: "blue",
		HOSTAGE: "cyan",
	};
</script>

<Tile style="position: relative; min-height: 120px;">
	{@const firstRound = roundInfos[0]}
	{@const playTime = new Date(firstRound.Time).toLocaleString()}
	<div style="z-index: 2; position: relative">
		<MatchProgressIndicator style="margin-bottom: 5px">
			{#each roundInfos as roundInfo, i}
				<MatchProgressStep won={roundInfo.WonRound} />
			{/each}
		</MatchProgressIndicator>

		<Tag>{playTime}</Tag>
		<Tag>{firstRound.GameMode}</Tag>
		<Tag type={matchTypeColors[firstRound.MatchType]}
			>{firstRound.MatchType}</Tag
		>
		<Tag type={gameModeColors[firstRound.GameMode]}>{firstRound.GameMode}</Tag>
		<Tag>{firstRound.SeasonSlug}</Tag>
	</div>

	<div style="z-index: 2;" id="map-name-container">
		<Tag type="outline">{firstRound.MapName}</Tag>
	</div>

	<div id="map-image-container">
		<div id="map-image">
			<ImageLoader
				src="/images/maps/{firstRound.MapName}.jpg"
				alt={firstRound.MapName}
				style="transform: translateY(-33.333%);"
			>
				<svelte:fragment slot="loading"
					><SkeletonPlaceholder
						style="width: 480px; height: 270px;"
					/></svelte:fragment
				>
				<svelte:fragment slot="error"
					><div id="map-image-placeholder">
						<bold>MAP UNKNOWN</bold>
					</div></svelte:fragment
				>
			</ImageLoader>
			<div id="map-image-fade" />
		</div>
	</div>
</Tile>

<style lang="scss">
	@use "@carbon/themes/scss/themes" as *;
	@use "@carbon/themes" with (
		$theme: $g90
	);

	#map-name-container {
		position: absolute;
		right: 0;
		bottom: 0;
	}

	#map-image-container {
		z-index: 1;
		position: absolute;
		top: 0;
		right: 0;
		height: 100%;

		overflow: hidden;
		opacity: 0.7;
	}

	#map-image {
		position: relative;
	}

	/* 
		FIXME: I'm not able to figure out how to have a fade-effect on the image and center it vertically at the same time
		without sacrificing automatic sizing of the image.
		This is the weird hack. IT IS NOT CENTERED VERTICALLY
	*/
	#map-image > * {
		transform: translateY(-33.333%);
	}

	#map-image-placeholder {
		width: 480px;
		height: 270px;
		background: linear-gradient(
			to right,
			themes.$layer-01 0%,
			themes.$layer-03 100%
		);
		text-align: center;
		line-height: 270px;
	}

	#map-image-fade {
		position: absolute;
		left: 0;
		top: 0;

		height: 100%;
		width: 100%;

		background: linear-gradient(
			to right,
			themes.$layer-01 0%,
			transparent 100%
		);
	}
</style>
