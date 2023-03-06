package game

import (
	"context"
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

type RoundsWatcher struct {
	gameDir string
}

func NewRoundsWatcher(gameDir string) (r *RoundsWatcher, err error) {
	_, errStat := os.Stat(gameDir)
	if errStat != nil {
		if os.IsNotExist(errStat) {
			err = fmt.Errorf(`game directory "%s" does not exist`, gameDir)
		} else {
			err = fmt.Errorf(`could not read game directory "%s"`, gameDir)
		}
		return
	}

	r = &RoundsWatcher{gameDir: gameDir}
	return
}

func (r *RoundsWatcher) Start(ctx context.Context) (<-chan RoundInfo, <-chan error) {
	chRounds := make(chan RoundInfo, 10)
	chErrors := make(chan error, 1)
	go r.watchMatches(ctx, chRounds, chErrors)
	return chRounds, chErrors
}

// watchMatches parses round replay files as they are created.
// It has an underlying file watcher goroutine that watches for changes in the game's replay directory.
func (r *RoundsWatcher) watchMatches(ctx context.Context, chRounds chan<- RoundInfo, chErrors chan<- error) {
	defer func() {
		log.Println("ending watchMatches")
		close(chRounds)
		close(chErrors)
	}()

	matchReplayFolder, ok := r.waitForMatchReplayDir(ctx)
	if !ok {
		return
	}

	chRoundReplayFiles := make(chan string, 10)
	chFileErrors := make(chan error, 1)

	go watchRounds(ctx, matchReplayFolder, chRoundReplayFiles, chFileErrors)
	for {
		select {
		case <-ctx.Done():
			log.Println("watchMatches context cancelled")
			return
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

// waitForMatchReplayDir waits for the replay dir in the game folder to exist.
// It returns the replay folder once found and a boolean indicating if successful.
// The boolean will only be false when the context has been cancelled.
func (r *RoundsWatcher) waitForMatchReplayDir(ctx context.Context) (string, bool) {
	log.Println("waiting for match replay folder...")
	matchReplayDir := path.Join(r.gameDir, constants.MATCH_REPLAY_FOLDER_NAME)

	var chInterval <-chan time.Time
	for {
		log.Println("checking", matchReplayDir)
		stat, err := os.Stat(matchReplayDir)
		if err == nil && stat.IsDir() {
			return matchReplayDir, true
		}
		chInterval = time.After(10 * time.Second)
		select {
		case <-ctx.Done():
			return "", false
		case <-chInterval:
			continue
		}
	}
}

func watchRounds(ctx context.Context, matchReplayFolder string, chFiles chan<- string, chErrors chan<- error) {
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

	listenFileChanges(ctx, matchDirWatcher, chFiles, chErrors)
}

func listenFileChanges(ctx context.Context, matchDirWatcher *fsnotify.Watcher, chFiles chan<- string, chErrors chan<- error) {
	for {
		select {
		case <-ctx.Done():
			log.Println("watchRounds context cancelled")
			return
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
