package game

import (
	"time"

	"github.com/redraskal/r6-dissect/dissect"
)

type RoundInfo struct {
	MatchID    string
	Time       time.Time
	SeasonSlug string
	MatchType  string
	GameMode   string
	MapName    string
	Teams      [2]Team
	Site       string
	// the following attributes relate to the recording player's team
	Won          bool
	WinCondition dissect.WinCondition
	TeamIndex    int
	PlayerName   string
}

type Team struct {
	Role    dissect.TeamRole
	Players []Player
}

func makeTeams(r *dissect.DissectReader) [2]Team {
	// initialize teams slice
	var teams [2]Team
	for i := 0; i < 2; i++ {
		teams[i] = Team{Role: r.Header.Teams[0].Role, Players: make([]Player, 0)}
	}

	// fill teams with players
	for _, player := range r.Header.Players {
		teams[player.TeamIndex].Players = append(teams[player.TeamIndex].Players, newPlayer(player))
	}

	return teams
}

type Player struct {
	Username string
	Operator string
}

func newPlayer(p dissect.Player) Player {
	return Player{
		Username: p.Username,
		Operator: p.RoleName,
	}
}
