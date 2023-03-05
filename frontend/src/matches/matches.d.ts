export namespace matches {
	export class RoundInfo {
		Players: Array<PlayerInfo>
		SeasonSlug: string
		RecordingPlayerName: string
		MatchID: string
		Time: string
		MatchType: "CASUAL" | "UNRANKED" | "RANKED"
		GameMode: "BOMB" | "HOSTAGE"
		MapName: string
		RoundWon: boolean
	}

	export class PlayerInfo {
		Username: string
		Operator: string
		Role: PlayerRole
	}

	export enum PlayerRole {
		Defense = 4,
		Attack = 0
	}
}
