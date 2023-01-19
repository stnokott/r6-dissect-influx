package constants

var (
	Version     = "dev"
	Commit      = "unknown"
	CompileTime = "unknown"
	ProjectName = "unknown"
)

const (
	APP_ID string = "org.stnokott.r6.dissect-influx"

	INFLUX_BATCH_SIZE = 500

	DEFAULT_MATCH_REPLAY_FOLDER_NAME string = "MatchReplay"

	INFO_CONNECTED = "Successfully connected"

	STATUS_DISCONNECTED = "Disconnected"
	STATUS_CONNECTED    = "Connected"
)
