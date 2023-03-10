package update

import (
	"errors"
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
		err = errors.Join(err, resp.Body.Close())
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
		err = errors.Join(err, outFile.Close())
	}()
	_, err = io.Copy(outFile, resp.Body)
	filepath = outFile.Name()
	return
}
