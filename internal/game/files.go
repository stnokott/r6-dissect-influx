package game

import (
	"fmt"
	"io/fs"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/stnokott/r6-dissect-influx/internal/constants"
)

type RoundsReader struct {
	gameDir string
}

func NewRoundsReader(gameDir string) (r *RoundsReader, err error) {
	_, errStat := os.Stat(gameDir)
	if errStat != nil {
		if os.IsNotExist(errStat) {
			err = fmt.Errorf(`game directory "%s" does not exist`, gameDir)
		} else {
			err = fmt.Errorf(`could not read game directory "%s"`, gameDir)
		}
		return
	}

	r = &RoundsReader{gameDir: gameDir}
	return
}

func (r *RoundsReader) WatchAsync() (<-chan RoundInfo, <-chan error) {
	chRounds := make(chan RoundInfo, 10)
	chErrors := make(chan error, 1)
	go r.watchMatches(chRounds, chErrors)
	return chRounds, chErrors
}

func (r *RoundsReader) watchMatches(chRounds chan<- RoundInfo, chErrors chan<- error) {
	defer func() {
		log.Println("ending watchMatches")
		close(chRounds)
		close(chErrors)
	}()
	matchReplayFolder := r.waitForMatchReplayDir()

	chRoundReplayFiles := make(chan string, 10)
	chFileErrors := make(chan error, 1)
	go watchRounds(matchReplayFolder, chRoundReplayFiles, chFileErrors)
	for {
		select {
		case filePath, ok := <-chRoundReplayFiles:
			if !ok {
				return
			}
			log.Println("found new round:", filePath)
			matchInfo, err := parseFile(filePath)
			if err != nil {
				chErrors <- err
			} else {
				chRounds <- matchInfo
			}
		case err, ok := <-chFileErrors:
			if !ok {
				return
			}
			if err != nil {
				chErrors <- err
			}
		}
	}
}

func (r *RoundsReader) waitForMatchReplayDir() string {
	matchReplayDir := path.Join(r.gameDir, constants.MATCH_REPLAY_FOLDER_NAME)
	log.Println("waiting for match replay folder...")
	for {
		log.Println("checking", matchReplayDir)
		stat, err := os.Stat(matchReplayDir)
		if err == nil && stat.IsDir() {
			return matchReplayDir
		}
		<-time.After(10 * time.Second)
	}
}

func watchRounds(matchReplayFolder string, chFiles chan<- string, chErrors chan<- error) {
	var errInitial error
	defer func() {
		log.Println("ending watchRounds")
		if errInitial != nil {
			chErrors <- errInitial
		}
		close(chFiles)
		close(chErrors)
	}()

	log.Println("getting initial match replay files")
	var initialReplayFiles []string
	initialReplayFiles, errInitial = getInitialMatchReplays(matchReplayFolder)
	if errInitial != nil {
		return
	}
	for _, replayFile := range initialReplayFiles {
		chFiles <- replayFile
	}

	var matchDirWatcher *fsnotify.Watcher
	if matchDirWatcher, errInitial = fsnotify.NewWatcher(); errInitial != nil {
		errInitial = fmt.Errorf("creating FS watcher failed: %w", errInitial)
		return
	}
	if errInitial = matchDirWatcher.Add(matchReplayFolder); errInitial != nil {
		errInitial = fmt.Errorf(`adding "%s" to FS watcher failed: %w`, matchReplayFolder, errInitial)
		return
	}
	log.Println("watching", matchReplayFolder, "...")

	for {
		select {
		case event, ok := <-matchDirWatcher.Events:
			if !ok {
				return
			}
			// log.Println("event:", event)
			if event.Has(fsnotify.Create) {
				// new folder for game created / new game started
				log.Println("created:", event.Name)
				st, err := os.Lstat(event.Name)
				if err != nil {
					log.Println("WARNING: could not stat", event.Name)
				} else {
					if st.IsDir() {
						if err = matchDirWatcher.Add(event.Name); err != nil {
							log.Println("WARNING: could not start watching", event.Name)
						} else {
							log.Println("now watching", event.Name)
						}
					} else if isReplayFile(event.Name) {
						chFiles <- event.Name
					} else {
						log.Println("WARNING: discarding file", event.Name)
					}
				}
			}
		case errWatcher, ok := <-matchDirWatcher.Errors:
			if !ok {
				return
			}
			chErrors <- errWatcher
		}
	}
}

// isReplayFile checks if the provided filepath refers to a match replay file.
// It assumes that the provided path points to a file, not a dir.
func isReplayFile(path string) bool {
	return strings.HasSuffix(path, constants.MATCH_REPLAY_SUFFIX)
}

func getInitialMatchReplays(matchReplayFolder string) ([]string, error) {
	var result []string
	err := filepath.WalkDir(matchReplayFolder, func(path string, d fs.DirEntry, errWalk error) error {
		if errWalk != nil {
			log.Printf("WARNING: discarding %s: %v", path, errWalk)
		} else if !d.IsDir() && isReplayFile(path) {
			result = append(result, path)
		}
		return nil
	})
	return result, err
}
