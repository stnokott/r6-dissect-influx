// Auto-generated file by goversioninfo. Do not edit.
package constants

import (
	"encoding/json"

	"github.com/josephspurrier/goversioninfo"
)

func unmarshalGoVersionInfo(b []byte) goversioninfo.VersionInfo {
	vi := goversioninfo.VersionInfo{}
	json.Unmarshal(b, &vi)
	return vi
}

var versionInfo = unmarshalGoVersionInfo([]byte(`{
	"FixedFileInfo":{
		"FileVersion": {
			"Major": 0,
			"Minor": 2,
			"Patch": 2,
			"Build": 0
		},
		"ProductVersion": {
			"Major": 0,
			"Minor": 2,
			"Patch": 2,
			"Build": 0
		},
		"FileFlagsMask": "3f",
		"FileFlags": "",
		"FileOS": "040004",
		"FileType": "01",
		"FileSubType": "00"
	},
	"StringFileInfo":{
		"Comments": "https://github.com/stnokott/r6-dissect-influx",
		"CompanyName": "github.com/stnokott",
		"FileDescription": "Pushes R6 Siege match data to InfluxDB",
		"FileVersion": "0.2.2",
		"InternalName": "r6-dissect-influx",
		"LegalCopyright": "github.com/stnokott",
		"LegalTrademarks": "",
		"OriginalFilename": "",
		"PrivateBuild": "",
		"ProductName": "r6-dissect-influx",
		"ProductVersion": "v0.2.2",
		"SpecialBuild": ""
	},
	"VarFileInfo":{
		"Translation": {
			"LangID": 1033,
			"CharsetID": 1200
		}
	}
}`))
