package config

import (
	"errors"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"

	"github.com/stnokott/r6-dissect-influx/internal/constants"
	"golang.org/x/sys/windows/registry"
)

// TODO: only configure *game* directory, since no MatchReplay folder exists on fresh install / new season

const gameFolderRegistryKey string = `SOFTWARE\WOW6432Node\Ubisoft\Launcher\Installs\635`

func matchReplayFolderFromRegistry() (result string, err error) {
	var key registry.Key
	key, err = registry.OpenKey(registry.LOCAL_MACHINE, gameFolderRegistryKey, registry.QUERY_VALUE)
	if err != nil {
		err = errors.New("game directory not found")
		return
	}
	defer func() {
		errInner := key.Close()
		if errInner != nil {
			// ignore error as it is not relevant for this function
			log.Printf("WARNING: could not close registry key: %v", errInner)
		}
	}()
	var gameDir string
	gameDir, _, err = key.GetStringValue("InstallDir")
	if err != nil {
		err = errors.New("game directory not found")
		return
	}

	_, err = os.Stat(gameDir)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf(`game directory "%s" found in Registry, but does not exist`, gameDir)
		} else {
			err = fmt.Errorf(`game directory "%s" found in Registry, but could not read folder: %w`, gameDir, err)
		}
		return
	}

	result = filepath.Join(gameDir, constants.DEFAULT_MATCH_REPLAY_FOLDER_NAME)
	var folderInfo fs.FileInfo
	folderInfo, err = os.Stat(result)
	if err != nil {
		if os.IsNotExist(err) {
			err = fmt.Errorf(`folder "%s" in game directory "%s" not found`, constants.DEFAULT_MATCH_REPLAY_FOLDER_NAME, gameDir)
		} else {
			err = fmt.Errorf(`could not read folder "%s"`, result)
		}
		return
	} else if !folderInfo.IsDir() {
		err = fmt.Errorf(`"%s" is not a folder`, result)
	}

	return
}
