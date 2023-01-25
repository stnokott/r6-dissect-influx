//go:build ignore

package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/josephspurrier/goversioninfo"
)

func main() {
	major := flag.Int("major", -1, "Major part of version (X.0.0)")
	minor := flag.Int("minor", -1, "Minor part of version (0.X.0)")
	patch := flag.Int("patch", -1, "Patch part of version (0.0.X)")
	projectName := flag.String("project_name", "", "Name of the GitHub project")
	outpath := flag.String("o", "", "Filepath of output JSON")
	flag.Parse()

	version := goversioninfo.FileVersion{
		Major: *major,
		Minor: *minor,
		Patch: *patch,
		Build: 0,
	}

	// see github.com/josephspurrier/goversioninfo
	info := goversioninfo.VersionInfo{
		FixedFileInfo: goversioninfo.FixedFileInfo{
			FileVersion:    version,
			ProductVersion: version,
			FileFlagsMask:  "3f",
			FileFlags:      "00",
			FileOS:         "040004",
			FileType:       "01",
			FileSubType:    "00",
		},
		StringFileInfo: goversioninfo.StringFileInfo{
			Comments:        "https://github.com/stnokott/" + *projectName,
			CompanyName:     "github.com/stnokott",
			FileDescription: "Pushes R6 Siege match data to InfluxDB",
			FileVersion:     fmt.Sprintf("%d.%d.%d", *major, *minor, *patch),
			InternalName:    *projectName,
			LegalCopyright:  "github.com/stnokott",
			ProductName:     *projectName,
			ProductVersion:  fmt.Sprintf("v%d.%d.%d", *major, *minor, *patch),
		},
		VarFileInfo: goversioninfo.VarFileInfo{
			Translation: goversioninfo.Translation{
				LangID:    goversioninfo.LngUSEnglish,
				CharsetID: goversioninfo.CsUnicode,
			},
		},
		IconPath: "../../assets/icon.ico",
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		panic(err)
	}

	if err = os.WriteFile(*outpath, data, 0644); err != nil {
		panic(err)
	}
}
