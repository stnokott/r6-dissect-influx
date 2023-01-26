package game

import (
	"time"
)

type RoundInfo struct {
	Players             []Player
	SeasonSlug          string
	RecordingPlayerName string
	MatchID             string
	Time                time.Time
	MatchType           string
	GameMode            string
	MapName             string
}

type Role int

const (
	ROLE_DEFENCE Role = 4
	ROLE_ATTACK  Role = 0
)

type Player struct {
	Username string
	Operator string
	Role     Role
}
