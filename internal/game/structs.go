package game

import (
	"time"

	"github.com/redraskal/r6-dissect/dissect"
)

type RoundInfo struct {
	MatchID             string
	Time                time.Time
	SeasonSlug          string
	RecordingPlayerName string
	MatchType           string
	GameMode            string
	MapName             string
	Teams               [2]Team
	Site                string
	WonRound            bool
	WinCondition        dissect.WinCondition
}

type Team struct {
	Role    dissect.TeamRole
	Players []Player
}

type Player struct {
	Username string
	Operator string
}
