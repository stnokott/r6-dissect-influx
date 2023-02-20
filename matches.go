package main

import (
	"log"

	"github.com/stnokott/r6-dissect-influx/internal/game"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) StartRoundWatcher() error {
	reader, err := game.NewRoundsReader(a.config.Game.InstallDir)
	if err != nil {
		return err
	}
	chRoundInfos, chErrors := reader.WatchAsync()
	go func() {
		defer runtime.EventsEmit(a.ctx, eventNames.RoundWatcherStopped)

		for {
			select {
			case roundInfo, ok := <-chRoundInfos:
				if !ok {
					return
				}
				log.Println("got match info for ID:", roundInfo.MatchID)
				runtime.EventsEmit(a.ctx, eventNames.NewRound, roundInfo)
			case err, ok := <-chErrors:
				if !ok {
					return
				}
				if err != nil {
					runtime.EventsEmit(a.ctx, eventNames.RoundWatcherError, err)
				}
			}
		}
	}()

	return nil
}
