package constants

import (
	"net/url"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/stnokott/r6-dissect-influx/internal/utils"
)

var (
	Version     = "v0.0.1"
	Commit      = "unknown"
	ProjectName = "r6-dissect-influx"
	SemVer      = version.Must(version.NewSemver(Version))
	GithubURL   = utils.MustArg(url.Parse, "https://github.com/stnokott/"+ProjectName)

	UpdateCheckInterval = 5 * time.Minute
)

const (
	APP_ID       string = "org.stnokott.r6.dissect-influx"
	WINDOW_TITLE string = "R6 Match Parser"

	MATCH_REPLAY_FOLDER_NAME string = "MatchReplay"
	MATCH_REPLAY_SUFFIX      string = ".rec"
)
