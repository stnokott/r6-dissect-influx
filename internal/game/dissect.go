package game

import (
	"errors"
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
		err = errors.Join(err, file.Close())
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
	recordingPlayer := r.Header.RecordingPlayer()
	info = RoundInfo{
		MatchID:      r.Header.MatchID,
		Time:         r.Header.Timestamp,
		SeasonSlug:   r.Header.GameVersion,
		MatchType:    r.Header.MatchType.String(),
		GameMode:     r.Header.GameMode.String(),
		MapName:      r.Header.Map.String(),
		Teams:        makeTeams(r),
		Site:         r.Header.Site,
		Won:          recordingPlayer.TeamIndex == winningTeamIndex,
		WinCondition: winningTeam.WinCondition,
		TeamIndex:    recordingPlayer.TeamIndex,
		PlayerName:   recordingPlayer.Username,
	}
	return
}

func getWinningTeamIndex(r *dissect.DissectReader) int {
	if r.Header.Teams[0].Won {
		return 0
	} else {
		return 1
	}
}
