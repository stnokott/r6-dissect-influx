package update

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/hashicorp/go-version"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

type Release struct {
	URL         string           `json:"html_url"`
	SemVer      *version.Version `json:"tag_name"`
	PublishedAt time.Time        `json:"published_at"`
	Assets      []asset          `json:"assets"`
}

type asset struct {
	ID          int    `json:"id"`
	Filename    string `json:"name"`
	ContentType string `json:"content_type"`
}

const expectedReleaseContentType string = "application/x-zip-compressed"

func (r *Release) getBinaryAsset() (a asset, err error) {
	regexReleaseName, err := regexp.Compile(fmt.Sprintf(
		`^%s_%s_windows_amd64\.zip`,
		constants.ProjectName,
		strings.ReplaceAll(r.SemVer.Core().String(), ".", `\.`),
	))
	if err != nil {
		err = fmt.Errorf("could not compile release name regex: %w", err)
		return
	}
	for _, asset := range r.Assets {
		if regexReleaseName.MatchString(asset.Filename) && asset.ContentType == expectedReleaseContentType {
			a = asset
			return
		}
	}
	err = fmt.Errorf("no qualifying asset found among %d assets in release %s (%s)", len(r.Assets), r.SemVer.Original(), r.URL)
	return
}

func (r *Release) IsNewer() bool {
	return r.SemVer.GreaterThan(constants.SemVer)
}
