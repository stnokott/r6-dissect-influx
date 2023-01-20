package game

import (
	"os"

	"github.com/redraskal/r6-dissect/reader"
)

func parseFile(f string) (info RoundInfo, err error) {
	var r *os.File
	r, err = os.Open(f)
	if err != nil {
		return
	}
	defer func() {
		errClose := r.Close()
		if errClose != nil && err == nil {
			err = errClose
		}
	}()

	var c reader.DissectReader
	c, err = reader.NewReader(r)
	if err != nil {
		return
	}

	players := make([]Player, len(c.Header.Players))
	for i, player := range c.Header.Players {
		players[i] = Player{
			Username: player.Username,
			Operator: player.RoleName,
			Role:     Role(player.Alliance),
		}
	}
	info = RoundInfo{
		SeasonSlug:          c.Header.GameVersion,
		RecordingPlayerName: c.Header.RecordingPlayer().Username,
		MatchID:             c.Header.MatchID,
		Time:                c.Header.Timestamp,
		MatchType:           c.Header.MatchType.String(),
		GameMode:            c.Header.GameMode.String(),
		MapName:             c.Header.Map.String(),
		Players:             players,
	}
	return
}
