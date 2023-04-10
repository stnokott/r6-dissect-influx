package main

import (
	"context"
	"errors"
	"log"

	"github.com/stnokott/r6-dissect-influx/internal/db"
	"github.com/stnokott/r6-dissect-influx/internal/game"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// StartRoundWatcher starts a background task that monitors the game's replay directory for new files.
// Once a new round if found and parsed, an event with its data will be emitted.
func (a *App) StartRoundWatcher() error {
	watcher, err := game.NewRoundsWatcher(a.config.Game.InstallDir)
	if err != nil {
		return err
	}
	var ctxRoundWatcher context.Context
	ctxRoundWatcher, a.roundsWatcherStop = context.WithCancel(context.Background())
	chRoundInfos, chErrors := watcher.Start(ctxRoundWatcher)
	go func() {
		defer runtime.EventsEmit(a.ctx, eventNames.RoundWatcherStopped)

		for {
			select {
			case <-ctxRoundWatcher.Done():
				log.Println("roundsWatcher in App cancelled")
				return
			case roundInfo, ok := <-chRoundInfos:
				if !ok {
					return
				}
				log.Println("got match info for ID:", roundInfo.MatchID)
				runtime.EventsEmit(a.ctx, eventNames.NewRound, roundInfo)
				db.AddRound(roundInfo)
			case err, ok := <-chErrors:
				if !ok {
					return
				}
				if err != nil {
					runtime.EventsEmit(a.ctx, eventNames.RoundWatcherError, err.Error())
				}
			}
		}
	}()

	runtime.EventsEmit(a.ctx, eventNames.RoundWatcherStarted)
	return nil
}

// StopRoundWatcher stops the background task started by StartRoundWatcher.
// It will return an error if no such task is running.
func (a *App) StopRoundWatcher() error {
	if a.roundsWatcherStop == nil {
		return errors.New("no round watcher running")
	}
	a.roundsWatcherStop()
	return nil
}
