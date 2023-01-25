package constants

import (
	"net/url"

	"github.com/hashicorp/go-version"
	"github.com/stnokott/r6-dissect-influx/internal/utils"
)

var (
	Version     = "v0.0.1"
	Commit      = "unknown"
	ProjectName = "r6-dissect-influx"
	SemVer      = version.Must(version.NewSemver(Version))
	GithubURL   = utils.Must(url.Parse, "https://github.com/stnokott/"+ProjectName)
)

const (
	APP_ID       string = "org.stnokott.r6.dissect-influx"
	WINDOW_TITLE string = "R6 Match Parser"

	INFLUX_BATCH_SIZE = 500

	MATCH_REPLAY_FOLDER_NAME string = "MatchReplay"
	MATCH_REPLAY_SUFFIX      string = ".rec"

	INFO_CONNECTED = "Successfully connected"

	STATUS_DISCONNECTED = "Disconnected"
	STATUS_CONNECTED    = "Connected"
)
