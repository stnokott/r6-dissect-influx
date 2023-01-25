package update

import (
	"io"
	"net/http"
	"os"
	"path"
	"strconv"

	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

func downloadAsset(a asset) (filepath string, err error) {
	assetURL := githubAPIReleasesBase + "/assets/" + strconv.Itoa(a.ID)

	var req *http.Request
	req, err = http.NewRequest("GET", assetURL, nil)
	if err != nil {
		return
	}
	req.Header.Add("Accept", "application/octet-stream")

	var resp *http.Response
	resp, err = http.DefaultClient.Do(req)
	if err != nil {
		return
	}
	defer func() {
		errInner := resp.Body.Close()
		if errInner != nil && err == nil {
			err = errInner
		}
	}()

	tmpFolder := path.Join(os.TempDir(), constants.ProjectName)
	err = os.Mkdir(tmpFolder, 0666)
	if err != nil && !os.IsExist(err) {
		return
	}

	var outFile *os.File
	outFile, err = os.Create(path.Join(tmpFolder, a.Filename))
	if err != nil {
		return
	}
	defer func() {
		errInner := outFile.Close()
		if errInner != nil && err == nil {
			err = errInner
		}
	}()
	_, err = io.Copy(outFile, resp.Body)
	filepath = outFile.Name()
	return
}
