package game

import (
	"fmt"
	"os"
	"path/filepath"

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
	}

	roundWon, err := hasWonRound(r, f)
	if err != nil {
		err = fmt.Errorf("error determining if round was won: %w", err)
		return
	}

	players := make([]Player, len(r.Header.Players))
	for i, player := range r.Header.Players {
		players[i] = Player{
			Username: player.Username,
			Operator: player.RoleName,
			Role:     Role(player.Alliance),
		}
	}
	info = RoundInfo{
		SeasonSlug:          r.Header.GameVersion,
		RecordingPlayerName: r.Header.RecordingPlayer().Username,
		MatchID:             r.Header.MatchID,
		Time:                r.Header.Timestamp,
		MatchType:           r.Header.MatchType.String(),
		GameMode:            r.Header.GameMode.String(),
		MapName:             r.Header.Map.String(),
		Players:             players,
		RoundWon:            roundWon,
	}
	return
}

func hasWonRound(r *dissect.DissectReader, replayFilePath string) (bool, error) {
	observerTeam := getObserverTeam(r)

	// check if this was the first round
	if r.Header.Teams[0].Score+r.Header.Teams[1].Score == 1 {
		// we have won if the observing player's team has 1 point
		return observerTeam.Score == 1, nil
	}

	// at this point, we know that this was not the first round.
	// we need to find out the previous round's score to determine who won this round.
	matchDir := filepath.Dir(replayFilePath)
	// FIXME: wait for update of github.com/redraskal/r6-dissect that includes info about round end directly in *dissect.DissectReader
	// since reading a whole match is a lot of overhead.
	previousRound, err := getRoundByIndex(matchDir, r.Header.RoundNumber)
	if err != nil {
		return false, err
	}

	prevObserverTeam := getObserverTeam(previousRound)
	return observerTeam.Score > prevObserverTeam.Score, nil
}

func getRoundByIndex(matchDir string, roundIndex int) (*dissect.DissectReader, error) {
	matchReader, err := dissect.NewMatchReader(matchDir)
	if err != nil {
		return nil, err
	}
	if roundIndex >= matchReader.NumRounds() {
		return nil, fmt.Errorf("round index %d out of bounds (have %d rounds)", roundIndex, matchReader.NumRounds())
	}
	return matchReader.RoundAt(roundIndex), nil
}

func getObserverTeam(r *dissect.DissectReader) dissect.Team {
	teams := r.Header.Teams
	if teams[0].Name == "YOUR TEAM" {
		return teams[0]
	} else {
		return teams[1]
	}
}
