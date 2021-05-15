package runner

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/chippydip/go-sc2ai/client"
)

var (
	launchBaseBuild        = uint32(0)
	launchDataVersion      = ""
	launchPortStart        = 8168
	launchExtraCommandArgs = []string(nil)
)

// SetGameVersion specifies a specific base game and data version to use when launching.
func SetGameVersion(baseBuild uint32, dataVersion string) {
	launchBaseBuild = baseBuild
	launchDataVersion = dataVersion
}

func (config *gameConfig) reLaunchStarcraft() {
	config.killAll()
	config.launchStarcraft()
}

func (config *gameConfig) launchStarcraft() {
	if len(config.clients) == 0 {
		log.Panic("No agents set")
	}

	portStart := 0
	if len(config.processInfo) != len(config.clients) {
		config.killAll()
		config.processInfo = config.launchProcesses(config.clients)
		portStart = launchPortStart + len(config.processInfo) - 1
	}

	config.setupPorts(len(config.clients), portStart, true)
	config.started = true
	config.lastPort = portStart
}

func (config *gameConfig) killAll() {
	for _, pi := range config.processInfo {
		if proc, err := os.FindProcess(pi.PID); err == nil && proc != nil {
			proc.Kill()
		}
	}
	config.processInfo = nil
}

func (config *gameConfig) launchProcesses(clients []*client.Client) []client.ProcessInfo {
	// Make sure we have a valid executable path
	path := processPathForBuild(launchBaseBuild)
	if _, err := os.Stat(path); err != nil {
		log.Print("Executable path can't be found, try running the StarCraft II executable first.")
		if len(path) > 0 {
			log.Printf("%v does not exist on your filesystem.", path)
		}
	}

	info := make([]client.ProcessInfo, len(clients))

	// Start an sc2 process for each bot
	var wg sync.WaitGroup
	for i, c := range clients {
		wg.Add(1)
		go func(i int, c *client.Client) {
			defer wg.Done()

			info[i] = config.launchAndAttach(path, c)

		}(i, c)
	}
	wg.Wait()

	return info
}

func (config *gameConfig) launchAndAttach(path string, c *client.Client) client.ProcessInfo {
	pi := client.ProcessInfo{}
	pi.Port = launchPortStart + len(config.processInfo) - 1

	// See if we can connect to an old instance real quick before launching
	if err := c.TryConnect(config.netAddress, pi.Port); err != nil {
		args := []string{
			"-listen", config.netAddress,
			"-port", strconv.Itoa(pi.Port),
			// DirectX will fail if multiple games try to launch in fullscreen mode. Force them into windowed mode.
			"-displayMode", "0",
		}

		if len(launchDataVersion) > 0 {
			args = append(args, "-dataVersion", launchDataVersion)
		}
		args = append(args, launchExtraCommandArgs...)

		// TODO: window size and position

		pi.Path = path
		pi.PID = startProcess(pi.Path, args)
		if pi.PID == 0 {
			log.Print("Unable to start sc2 executable with path: ", pi.Path)
		} else {
			log.Printf("Launched SC2 (%v), PID: %v", pi.Path, pi.PID)
		}

		// Attach
		if err := c.Connect(config.netAddress, pi.Port, processConnectTimeout); err != nil {
			log.Panic("Failed to connect")
		}
	}

	c.SetProcessInfo(pi)
	return pi
}

func startProcess(path string, args []string) int {
	cmd := exec.Command(path, args...)

	// Set the working directory on windows
	if runtime.GOOS == "windows" {
		_, exe := filepath.Split(path)
		dir := sc2Path(path)
		if strings.Contains(exe, "_x64") {
			dir = filepath.Join(dir, "Support64")
		} else {
			dir = filepath.Join(dir, "Support")
		}
		cmd.Dir = dir
	}

	if err := cmd.Start(); err != nil {
		log.Print(err)
		return 0
	}

	return cmd.Process.Pid
}
