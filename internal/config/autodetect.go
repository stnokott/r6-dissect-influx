// go:build windows

package config

import (
	"errors"
	"fmt"
	"log"

	"golang.org/x/sys/windows/registry"
)

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

	if err = gameDirectoryValidator(gameDir); err != nil {
		err = fmt.Errorf(`game directory "%s" found, but: %w`, gameDir, err)
	} else {
		result = gameDir
	}

	return
}
