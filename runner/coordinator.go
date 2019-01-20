package runner

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

var clients []*client.Client
var started = false
var lastPort = 0

// SetParticipants ...
func SetParticipants(participants ...client.PlayerSetup) {
	gameSettings.playerSetup, clients, started, lastPort, processSettings.processInfo = nil, nil, false, 0, nil
	for _, p := range participants {
		if p.Agent != nil {
			clients = append(clients, &client.Client{Agent: p.Agent})
		}
		gameSettings.playerSetup = append(gameSettings.playerSetup, p.PlayerSetup)
	}
}

// LaunchStarcraft ...
func LaunchStarcraft() {
	if _, err := os.Stat(processSettings.processPath); err != nil {
		log.Print("Executable path can't be found, try running the StarCraft II executable first.")
		if len(processSettings.processPath) > 0 {
			log.Printf("%v does not exist on your filesystem.", processSettings.processPath)
		}
		os.Exit(1)
	}

	if len(clients) == 0 {
		log.Panic("No agents set")
	}

	portStart := 0
	if len(processSettings.processInfo) != len(clients) {
		portStart = launchProcesses()
	}

	SetupPorts(len(clients), portStart, true)
	started = true
	lastPort = portStart
}

func reLaunchStarcraft() {
	for _, pi := range processSettings.processInfo {
		if proc, err := os.FindProcess(pi.PID); err == nil && proc != nil {
			proc.Kill()
		}
	}
	processSettings.processInfo = nil

	LaunchStarcraft()
}

// StartGame ...
func StartGame(mapPath string) {
	if !CreateGame(mapPath) {
		log.Fatal("Failed to create game.")
	}
	JoinGame()
}

// CreateGame ...
func CreateGame(mapPath string) bool {
	if mapPath == "" {
		mapPath = gameSettings.mapName
	}
	if !started {
		log.Panic("Game not started")
	}

	// Create with the first client
	err := clients[0].CreateGame(gameSettings.mapName, gameSettings.playerSetup, processSettings.realtime)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

// JoinGame ...
func JoinGame() bool {
	// TODO: Make this parallel and get rid of the WaitJoinGame method
	for i, client := range clients {
		if err := client.RequestJoinGame(gameSettings.playerSetup[i], interfaceOptions, gameSettings.ports); err != nil {
			log.Fatalf("Unable to join game: %v", err)
		}
	}

	// Check if any errors occurred during game start
	// errors := false
	// for _, client := range clients {
	// 	errs := client.GetClientErrors()
	// 	if len(errs) > 0 {
	// 		client.Agent.OnError(errs, agent.Control().GetProtocolErrors())
	// 		errors = true
	// 	}

	// 	//agent.Control().UseGeneralizedAbility(useGeneralizedAbilityID)
	// }
	// if errors {
	// 	return false
	// }

	// Run all clients on game start
	// for _, agent := range agents {
	// 	agent.Control().GetObservation()
	// }
	// for _, agent := range agents {
	// 	agent.OnGameFullStart()
	// }
	// for _, agent := range agents {
	// 	agent.Control().OnGameStart()
	// 	agent.OnGameStart()
	// }
	// for _, agent := range agents {
	// 	agent.Control().IssueEvents(agent.Actions().Commands())
	// }
	return true
}

func launchProcesses() int {
	processSettings.processInfo = make([]client.ProcessInfo, len(clients))

	// Start an sc2 process for each bot
	var wg sync.WaitGroup
	for i, c := range clients {
		wg.Add(1)
		go func(i int, c *client.Client) {
			defer wg.Done()

			launchAndAttach(c, i)

		}(i, c)
	}
	wg.Wait()

	return processSettings.portStart + len(clients) - 1
}

func launchAndAttach(c *client.Client, clientIndex int) {
	timeout := time.Duration(processSettings.timeoutMS) * time.Millisecond

	pi := client.ProcessInfo{}
	pi.Port = processSettings.portStart + len(processSettings.processInfo) - 1

	// See if we can connect to an old instance real quick before launching
	if err := c.TryConnect(processSettings.netAddress, pi.Port); err != nil {
		args := []string{
			"-listen", processSettings.netAddress,
			"-port", strconv.Itoa(pi.Port),
			// DirectX will fail if multiple games try to launch in fullscreen mode. Force them into windowed mode.
			"-displayMode", "0",
		}

		if len(processSettings.dataVersion) > 0 {
			args = append(args, "-dataVersion", processSettings.dataVersion)
		}
		args = append(args, processSettings.extraCommandLines...)

		// TODO: window size and position

		path := processSettings.processPath
		if processSettings.baseBuild != 0 {
			// Get the exe name and then back out to the Versions directory
			dir, exe := filepath.Split(path)
			dir = filepath.Dir(path) // remove trailing slash
			for dir != "." && filepath.Base(dir) != "Versions" {
				dir = filepath.Dir(dir)
			}

			// Get the path of the correct version and make sure the exe exists
			path = filepath.Join(dir, fmt.Sprintf("Base%v", processSettings.baseBuild), exe)
			if _, err := os.Stat(path); err != nil {
				log.Fatalf("Base version not found: %v", err)
			}
		}

		pi.Path = path
		pi.PID = startProcess(path, args)
		if pi.PID == 0 {
			log.Print("Unable to start sc2 executable with path: ", path)
		} else {
			log.Printf("Launched SC2 (%v), PID: %v", path, pi.PID)
		}

		// Attach
		if err := c.Connect(processSettings.netAddress, pi.Port, timeout); err != nil {
			log.Panic("Failed to connect")
		}
	}

	c.SetProcessInfo(pi)
	processSettings.processInfo[clientIndex] = pi
}

func startProcess(path string, args []string) int {
	cmd := exec.Command(path, args...)

	// Set the working directory on windows
	if runtime.GOOS == "windows" {
		path, exe := filepath.Split(path)
		path = filepath.Dir(path) // remove trailing slash
		for path != "." && filepath.Base(path) != "StarCraft II" {
			path = filepath.Dir(path)
		}
		if strings.Contains(exe, "_x64") {
			path = filepath.Join(path, "Support64")
		} else {
			path = filepath.Join(path, "Support")
		}
		cmd.Dir = path
	}

	if err := cmd.Start(); err != nil {
		log.Print(err)
		return 0
	}

	return cmd.Process.Pid
}

func attachClients() {
	// Since connect is blocking do it after the processes are launched.
	timeout := time.Duration(processSettings.timeoutMS) * time.Millisecond
	for i, client := range clients {
		pi := processSettings.processInfo[i]

		if err := client.Connect(processSettings.netAddress, pi.Port, timeout); err != nil {
			log.Panic("Failed to connect")
		}
	}
}

// Connect ...
func Connect(port int) {
	pi := client.ProcessInfo{Path: processSettings.netAddress, PID: 0, Port: port}

	// Set process info for each bot
	for range clients {
		processSettings.processInfo = append(processSettings.processInfo, pi)
	}

	attachClients()

	// Assume starcraft has started after succesfully attaching to a server
	started = true
}

// SetupPorts ...
func SetupPorts(numAgents int, startPort int, checkSingle bool) {
	humans := numAgents
	if checkSingle {
		humans = 0
		for _, p := range gameSettings.playerSetup {
			if p.Type == api.PlayerType_Participant {
				humans++
			}
		}
	}

	if humans > 1 {
		var ports = gameSettings.ports
		ports.SharedPort = int32(startPort + 1)
		ports.ServerPorts = &api.PortSet{
			GamePort: int32(startPort + 2),
			BasePort: int32(startPort + 3),
		}

		for i := 0; i < numAgents; i++ {
			var base = int32(startPort + 4 + i*2)
			ports.ClientPorts = append(ports.ClientPorts, &api.PortSet{GamePort: base, BasePort: base + 1})
		}
		gameSettings.ports = ports
	}
}

// Run ...
func Run() {
	wg := sync.WaitGroup{}
	wg.Add(len(clients))

	for _, c := range clients {
		go func(client *client.Client) {
			defer wg.Done()

			runAgent(client)
			cleanup(client)
		}(c)
	}

	wg.Wait()
}

func runAgent(c *client.Client) {
	defer func() {
		if p := recover(); p != nil {
			client.ReportPanic(p)
		}

		// If the bot crashed before losing, keep the game running (force the opponent to earn the win)
		for c.IsInGame() {
			if err := c.Step(224); err != nil { // 10 seconds per update
				log.Print(err)
				break
			}
		}
	}()

	err := c.Init() // get GameInfo, Data, and Observation
	if err == nil && c.IsInGame() {
		c.Agent.RunAgent(c)
	}
}

func cleanup(c *client.Client) {
	if gamePort == 0 {
		// Leave the game (but only in non-ladder games)
		c.RequestLeaveGame()
	}

	// Print the winner
	for _, player := range c.Observation().GetPlayerResult() {
		if player.GetPlayerId() == c.PlayerID() {
			log.Print(player.GetResult())
		}
	}
}

// StartReplay ...
func StartReplay(path string) {
	// Get info about the replay
	info, err := clients[0].RequestReplayInfo(path)
	if err != nil {
		log.Fatalf("Unable to get replay info: %v", err)
	}

	// Check if we need to re-launch the game
	current := clients[0].Proto()
	if info.GetBaseBuild() != current.GetBaseBuild() || info.GetDataVersion() != current.GetDataVersion() {
		log.Printf("Version mis-match, relaunching client")
		processSettings.baseBuild = info.GetBaseBuild()
		processSettings.dataVersion = info.GetDataVersion()

		reLaunchStarcraft()

		current = clients[0].Proto()
		if info.GetBaseBuild() != current.GetBaseBuild() {
			log.Fatalf("Failed to launch correct base build: %v %v", current.GetBaseBuild(), info.GetBaseBuild())
		}
		if info.GetDataVersion() != current.GetDataVersion() {
			log.Fatalf("Failed to launch correct data version: %v %v", current.GetDataVersion(), info.GetDataVersion())
		}
	}

	log.Printf("Launching replay: %v", path)
	err = clients[0].RequestStartReplay(api.RequestStartReplay{
		Replay: &api.RequestStartReplay_ReplayPath{
			ReplayPath: path,
		},
		ObservedPlayerId: replaySettings.player,
		Options:          interfaceOptions,
		Realtime:         processSettings.realtime,
	})
	if err != nil {
		log.Fatalf("Unable to start replay: %v", err)
	}
}

// SetReplayPath ...
func SetReplayPath(path string) error {
	replaySettings.files = nil

	if filepath.Ext(path) == ".SC2Replay" {
		replaySettings.files = []string{path}
		return nil
	}

	replaySettings.dir = path

	// Gather and append all files from the directory.
	files, err := ioutil.ReadDir(path)
	if err != nil {
		return err
	}
	for _, file := range files {
		if !file.IsDir() && filepath.Ext(file.Name()) == ".SC2Replay" {
			replaySettings.files = append(replaySettings.files, filepath.Join(path, file.Name()))
		}
	}
	return nil
}

// SetReplayPlayerID ...
func SetReplayPlayerID(player api.PlayerID) {
	replaySettings.player = player
}

// LoadReplayList ...
func LoadReplayList(path string) error {
	return errors.New("NYI")
}

// SaveReplayList ...
func SaveReplayList(path string) error {
	return errors.New("NYI")
}
