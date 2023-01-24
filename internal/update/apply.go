package update

import (
	"archive/zip"
	"fmt"
	"io"
	"strings"

	"github.com/minio/selfupdate"
)

func (r *Release) DownloadAndApply() (err error) {
	// find binary asset
	var asset asset
	asset, err = r.getBinaryAsset()
	if err != nil {
		return
	}

	var filepath string
	filepath, err = downloadAsset(asset)
	if err != nil {
		return
	}

	var reader *zip.ReadCloser
	reader, err = zip.OpenReader(filepath)
	if err != nil {
		return
	}
	defer func() {
		errInner := reader.Close()
		if errInner != nil && err == nil {
			err = errInner
		}
	}()

	var in io.ReadCloser
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, ".exe") {
			in, err = file.Open()
			defer func() {
				errInner := in.Close()
				if errInner != nil && err == nil {
					err = errInner
				}
			}()
			err = selfupdate.Apply(in, selfupdate.Options{})
			if err != nil {
				if errRollback := selfupdate.RollbackError(err); errRollback != nil {
					err = fmt.Errorf("failed to rollback after unsuccessful update: %w -> %w", err, errRollback)
				}
			}
			return
		}
	}
	err = fmt.Errorf("did not find qualifying EXE file in release asset %s", asset.Filename)
	return
}
