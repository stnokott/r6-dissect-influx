//go:generate goversioninfo

package main

import (
	"log"

	"fyne.io/fyne/v2/app"

	"github.com/rs/zerolog"
	"github.com/stnokott/r6-dissect-influx/internal/config"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"github.com/stnokott/r6-dissect-influx/internal/game"
	"github.com/stnokott/r6-dissect-influx/internal/root"
)

func main() {
	// necessary for r6-dissect
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	a := app.NewWithID(constants.APP_ID)
	// a.SetIcon(constants.APP_ICON)

	config.Init(a)

	reader, err := game.NewRoundsReader(config.Current.GameFolder)
	if err != nil {
		panic(err)
	}
	chRoundInfos, chErrors := reader.WatchAsync()
	go func() {
		for {
			select {
			case roundInfo, ok := <-chRoundInfos:
				if !ok {
					log.Println("match data channel closed")
					return
				}
				log.Println("got match info for ID:", roundInfo.MatchID)
			case err, ok := <-chErrors:
				if !ok {
					log.Println("match errors channel closed")
					return
				}
				if err != nil {
					log.Println("got error from match data channel:", err)
				}
			}
		}
	}()

	root.NewView(a).ShowAndRun()
}
