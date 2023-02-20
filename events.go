package main

type EventNames struct {
	NewRound             string
	RoundWatcherError    string
	RoundWatcherStopped  string
	LatestReleaseInfo    string
	LatestReleaseInfoErr string
	UpdateProgress       string
	UpdateErr            string
}

var eventNames = EventNames{
	NewRound:             "R6_NEW_ROUND",
	RoundWatcherError:    "R6_ROUND_WATCHER_ERROR",
	RoundWatcherStopped:  "R6_ROUND_WATCHER_STOPPED",
	LatestReleaseInfo:    "R6_RELEASE_INFO",
	LatestReleaseInfoErr: "R6_RELEASE_INFO_ERR",
	UpdateProgress:       "R6_UPDATE_PROGRESS",
	UpdateErr:            "R6_UPDATE_ERROR",
}

func (*App) GetEventNames() EventNames {
	return eventNames
}
