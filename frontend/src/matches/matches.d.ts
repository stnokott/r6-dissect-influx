export namespace matches {
	export class EventNames {
		NewRound: string
		RoundWatcherError: string
		RoundWatcherStopped: string
	}

	export class RoundInfo {
		// Players             []Player
		SeasonSlug: string
		RecordingPlayerName: string
		MatchID: string
		// Time                time.Time
		MatchType: string
		GameMode: string
		MapName: string
	}
}
