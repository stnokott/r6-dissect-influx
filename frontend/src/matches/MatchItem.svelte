<script lang="ts">
	import { Tag, Tile } from "carbon-components-svelte";
	import type { TagProps } from "carbon-components-svelte/types/Tag/Tag.svelte";
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

<Tile>
	{@const firstRound = roundInfos[0]}
	{@const playTime = new Date(firstRound.Time).toLocaleString()}
	<Tag>{playTime}</Tag>
	<Tag>{firstRound.GameMode}</Tag>
	<Tag type={matchTypeColors[firstRound.MatchType]}>{firstRound.MatchType}</Tag>
	<Tag type={gameModeColors[firstRound.GameMode]}>{firstRound.GameMode}</Tag>
	<Tag>{firstRound.MapName}</Tag>
	<Tag>{firstRound.SeasonSlug}</Tag>
</Tile>
