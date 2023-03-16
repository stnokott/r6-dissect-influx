package update

import (
	"archive/zip"
	"errors"
	"fmt"
	"io"
	"strings"

	"github.com/minio/selfupdate"
)

type updateProgress struct {
	Task string
	Err  error
}

func (r *Release) DownloadAndApply() <-chan updateProgress {
	ch := make(chan updateProgress, 3)
	go r.downloadAndApplySync(ch)
	return ch
}

func (r *Release) downloadAndApplySync(chProgress chan<- updateProgress) {
	var err error
	defer func() {
		if err != nil {
			chProgress <- updateProgress{Err: err}
		}
		close(chProgress)
	}()

	// find binary asset
	var asset asset
	asset, err = r.getBinaryAsset()
	if err != nil {
		return
	}

	chProgress <- updateProgress{Task: "Downloading release..."}
	var filepath string
	filepath, err = downloadAsset(asset)
	if err != nil {
		return
	}

	chProgress <- updateProgress{Task: "Unzipping release..."}
	var reader *zip.ReadCloser
	reader, err = zip.OpenReader(filepath)
	if err != nil {
		return
	}
	defer func() {
		err = errors.Join(err, reader.Close())
	}()

	if err = applyUpdate(reader, asset, chProgress); err != nil {
		return
	}
	err = reader.Close()
}

func applyUpdate(reader *zip.ReadCloser, asset asset, chProgress chan<- updateProgress) (err error) {
	var in io.ReadCloser
	for _, file := range reader.File {
		if strings.HasSuffix(file.Name, ".exe") {
			in, err = file.Open()
			defer func() {
				err = errors.Join(err, in.Close())
			}()
			chProgress <- updateProgress{Task: "Applying update..."}
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
