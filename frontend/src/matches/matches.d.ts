export namespace matches {
	export class RoundInfo {
		MatchID: string
		Time: string
		SeasonSlug: string
		MatchType: "CASUAL" | "UNRANKED" | "RANKED"
		GameMode: "BOMB" | "HOSTAGE"
		MapName: string
		Teams: [Team, Team]
		Site: string
		Won: boolean
		WinCondition: "KILLED_OPPONENTS" | "SECURED_AREA" | "DISABLED_DEFUSER" | "DEFUSED_BOMB" | "EXTRACTED_HOSTAGE" | "TIME"
		TeamIndex: number
		PlayerName: string
	}

	export class Team {
		Role: "ATTACK" | "DEFENSE"
		Players: Array<Player>
	}

	export class Player {
		Username: string
		Operator: string
	}
}
