package main

// EventNames contains the names for events emitted by the application during runtime.
type EventNames struct {
	NewRound             string
	RoundWatcherStarted  string
	RoundWatcherError    string
	RoundWatcherStopped  string
	LatestReleaseInfo    string
	LatestReleaseInfoErr string
	UpdateProgress       string
	UpdateErr            string
}

var eventNames = EventNames{
	NewRound:             "R6_NEW_ROUND",
	RoundWatcherStarted:  "R6_ROUNDS_WATCHER_STARTED",
	RoundWatcherError:    "R6_ROUND_WATCHER_ERROR",
	RoundWatcherStopped:  "R6_ROUND_WATCHER_STOPPED",
	LatestReleaseInfo:    "R6_RELEASE_INFO",
	LatestReleaseInfoErr: "R6_RELEASE_INFO_ERR",
	UpdateProgress:       "R6_UPDATE_PROGRESS",
	UpdateErr:            "R6_UPDATE_ERROR",
}

// GetEventNames returns the event names emitted by the application during runtime.
func (*App) GetEventNames() EventNames {
	return eventNames
}
