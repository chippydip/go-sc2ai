package runner

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/chippydip/go-sc2ai/api"
)

var (
	replayDir            = ""
	replayFiles          = []string(nil)
	replayFilter         = (func(info *api.ResponseReplayInfo) bool)(nil)
	replayObservedPlayer = api.PlayerID(1)
	replayCurrentFile    = ""
)

// SetReplayPath sets a directory of replay files or a single replay to load.
func SetReplayPath(path string) error {
	replayFiles = nil
	if p, err := filepath.Abs(path); err != nil {
		log.Printf("Failed to get absolute path: %v", err)
	} else {
		path = p
	}

	if isReplayFile(filepath.Ext(path)) {
		replayFiles = []string{path}
		return nil
	}

	replayDir = path

	// Gather and append all files from the directory.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() && isReplayFile(filepath.Ext(file.Name())) {
			replayFiles = append(replayFiles, filepath.Join(path, file.Name()))
		}
	}
	return nil
}

func isReplayFile(path string) bool {
	return strings.ToLower(path) == ".sc2replay"
}

// SetReplayPlayerID specifies which player the agent would like to observe. This can be called
// from the filter function specified in SetReplayFilter to determine the player dynamically.
func SetReplayPlayerID(player api.PlayerID) {
	replayObservedPlayer = player
}

// SetReplayFilter provides a filter which determines if a replay should be run. This allows such
// things as MMR filtering within a large group of replays. The filter function can also call
// SetReplayPlayerID before returning true to alter the player who will be observed for the replay.
func SetReplayFilter(filter func(info *api.ResponseReplayInfo) bool) {
	replayFilter = filter
}

// CurrentReplayPath provides access to the replay filename and full path of the current replay (if any).
func CurrentReplayPath() string {
	return replayCurrentFile
}

func runReplays(config *gameConfig) bool {
	if len(replayFiles) == 0 {
		return false
	}

	for _, file := range replayFiles {
		replayCurrentFile = file
		if startReplay(config, file) {
			run(config.clients)
		}
		replayCurrentFile = ""
	}
	return true
}

func startReplay(config *gameConfig, path string) bool {
	// TODO: Parse the replay header ourselves to determine the correct BaseBuild and DataVersion
	// since RequestReplayInfo seems to fail if the versions don't match (may be a new sc2 bug)
	// See https://github.com/GraylinKim/sc2reader/wiki/.sc2replay
	// and https://github.com/GraylinKim/sc2reader/wiki/Serialized-Data

	// Get info about the replay
	info, err := config.clients[0].RequestReplayInfo(path)
	if err != nil {
		log.Printf("Unable to get replay info: %v", err)
		return false
	}

	// Allow the bot user to skip certain replays after looking at the info
	if replayFilter != nil && !replayFilter(info) {
		log.Printf("Skipping replay: %v", path)
		return false
	}

	// Check if we need to re-launch the game
	current := config.clients[0].Proto()
	if info.GetBaseBuild() != current.GetBaseBuild() || info.GetDataVersion() != current.GetDataVersion() {
		log.Printf("Version mis-match, relaunching client")
		SetGameVersion(info.GetBaseBuild(), info.GetDataVersion())

		config.reLaunchStarcraft()

		current = config.clients[0].Proto()
		if info.GetBaseBuild() != current.GetBaseBuild() {
			log.Fatalf("Failed to launch correct base build: %v %v", current.GetBaseBuild(), info.GetBaseBuild())
		}
		if info.GetDataVersion() != current.GetDataVersion() {
			log.Fatalf("Failed to launch correct data version: %v %v", current.GetDataVersion(), info.GetDataVersion())
		}
	}

	log.Printf("Launching replay: %v", path)
	err = config.clients[0].RequestStartReplay(api.RequestStartReplay{
		Replay: &api.RequestStartReplay_ReplayPath{
			ReplayPath: path,
		},
		ObservedPlayerId: replayObservedPlayer,
		Options:          processInterfaceOptions,
		Realtime:         processRealtime,
	})
	if err != nil {
		log.Fatalf("Unable to start replay: %v", err)
	}

	return true
}
