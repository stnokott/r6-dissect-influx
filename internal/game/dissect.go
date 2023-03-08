package game

import (
	"os"

	"github.com/redraskal/r6-dissect/dissect"
)

func parseFile(f string) (info RoundInfo, err error) {
	var file *os.File
	file, err = os.Open(f)
	if err != nil {
		return
	}
	defer func() {
		errClose := file.Close()
		if errClose != nil && err == nil {
			err = errClose
		}
	}()

	var r *dissect.DissectReader
	r, err = dissect.NewReader(file)
	if err != nil {
		return
	}
	if err = r.Read(); !dissect.Ok(err) {
		return
	} else {
		// reset to nil since error is not deemed problematic
		err = nil
	}

	winningTeamIndex := getWinningTeamIndex(r)
	winningTeam := r.Header.Teams[winningTeamIndex]
	observingPlayer := r.Header.RecordingPlayer()
	info = RoundInfo{
		MatchID:             r.Header.MatchID,
		Time:                r.Header.Timestamp,
		SeasonSlug:          r.Header.GameVersion,
		RecordingPlayerName: r.Header.RecordingPlayer().Username,
		MatchType:           r.Header.MatchType.String(),
		GameMode:            r.Header.GameMode.String(),
		MapName:             r.Header.Map.String(),
		Teams:               makeTeams(r),
		Site:                r.Header.Site,
		WonRound:            observingPlayer.TeamIndex == winningTeamIndex,
		WinCondition:        winningTeam.WinCondition,
	}
	return
}

func makeTeams(r *dissect.DissectReader) [2]Team {
	// initialize teams slice
	var teams [2]Team
	for i := 0; i < 2; i++ {
		teams[i] = Team{Role: r.Header.Teams[0].Role, Players: make([]Player, 0)}
	}

	// fill teams with players
	for _, player := range r.Header.Players {
		teams[player.TeamIndex].Players = append(teams[player.TeamIndex].Players, Player{
			Username: player.Username,
			Operator: player.RoleName,
		})
	}

	return teams
}

func getWinningTeamIndex(r *dissect.DissectReader) int {
	if r.Header.Teams[0].Won {
		return 0
	} else {
		return 1
	}
}
