package update

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

var (
	githubAPIReleasesBase = "https://api.github.com/repos/stnokott/" + constants.ProjectName + "/releases"
	latestReleaseURL      = githubAPIReleasesBase + "/latest"
)

func GetLatestRelease() (result *Release, err error) {
	var resp *http.Response
	resp, err = http.Get(latestReleaseURL)
	if err != nil {
		return
	} else if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusNotFound {
			err = errors.New("no release found")
			return
		}
		var details string
		body, errInner := io.ReadAll(resp.Body)
		if errInner != nil {
			details = "unknown"
		} else {
			_ = resp.Body.Close()
			details = string(body)
		}
		err = fmt.Errorf("GET %s: non-ok HTTP status code %d: %v", latestReleaseURL, resp.StatusCode, details)
		return
	}

	defer func() {
		err = errors.Join(err, resp.Body.Close())
	}()

	result = new(Release)
	err = json.NewDecoder(resp.Body).Decode(result)
	return
}
