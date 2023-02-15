package main

type EventNames struct {
	NewRound            string
	RoundWatcherError   string
	RoundWatcherStopped string
}

var eventNames = EventNames{
	NewRound:            "R6_NEW_ROUND",
	RoundWatcherError:   "R6_ROUND_WATCHER_ERROR",
	RoundWatcherStopped: "R6_ROUND_WATCHER_STOPPED",
}

func (_ *App) GetEventNames() EventNames {
	return eventNames
}
