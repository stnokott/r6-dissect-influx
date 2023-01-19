package constants

var (
	Version = "dev"
	Commit  = "unknown"
)

const (
	APP_ID       string = "org.stnokott.r6.dissect-influx"
	WINDOW_TITLE string = "R6 Match Parser"

	INFLUX_BATCH_SIZE = 500

	DEFAULT_MATCH_REPLAY_FOLDER_NAME string = "MatchReplay"

	INFO_CONNECTED = "Successfully connected"

	STATUS_DISCONNECTED = "Disconnected"
	STATUS_CONNECTED    = "Connected"
)
