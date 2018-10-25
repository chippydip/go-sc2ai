package runner

import (
	"fmt"
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

// StartGame ...
func StartGame(mapPath string) {
	if !CreateGame(mapPath) {
		fmt.Println("Failed to create game.")
		os.Exit(1)
	}
	JoinGame()
}

// CreateGame ...
func CreateGame(mapPath string) bool {
	if mapPath == "" {
		mapPath = gameSettings.mapName
	}
	if !started {
		panic("Game not started")
	}

	// Create with the first client
	err := clients[0].CreateGame(gameSettings.mapName, gameSettings.playerSetup, processSettings.realtime)
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

// JoinGame ...
func JoinGame() bool {
	// TODO: Make this parallel and get rid of the WaitJoinGame method
	for i, client := range clients {
		if err := client.RequestJoinGame(gameSettings.playerSetup[i], interfaceOptions, gameSettings.ports); err != nil {
			fmt.Printf("Unable to join game: %v", err)
			os.Exit(1)
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

// SetParticipants ...
func SetParticipants(participants ...client.PlayerSetup) {
	gameSettings.playerSetup = nil
	clients = nil
	for _, p := range participants {
		if p.Agent != nil {
			clients = append(clients, &client.Client{Agent: p.Agent})
		}
		gameSettings.playerSetup = append(gameSettings.playerSetup, p)
	}
}

// LaunchStarcraft ...
func LaunchStarcraft() {
	if _, err := os.Stat(processSettings.processPath); err != nil {
		fmt.Println("Executable path can't be found, try running the StarCraft II executable first.")
		if len(processSettings.processPath) > 0 {
			fmt.Printf("%v does not exist on your filesystem.\n", processSettings.processPath)
		}
		os.Exit(1)
	}

	if len(clients) == 0 {
		panic("No agents set")
	}

	portStart := 0
	if len(processSettings.processInfo) != len(clients) {
		portStart = launchProcesses()
	}

	SetupPorts(len(clients), portStart, true)
	started = true
	lastPort = portStart
}

func launchProcesses() int {
	lastPort := 0
	// Start an sc2 process for each bot
	for i, client := range clients {
		lastPort = launchProcess(client, i)
	}

	attachClients()

	return lastPort
}

func launchProcess(c *client.Client, clientIndex int) int {
	pi := client.ProcessInfo{}
	pi.Port = processSettings.portStart + len(processSettings.processInfo) - 1

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

	pi.Path = processSettings.processPath
	pi.PID = startProcess(processSettings.processPath, args)
	if pi.PID == 0 {
		fmt.Printf("Unable to start sc2 executable with path: %v\n", processSettings.processPath)
	} else {
		fmt.Printf("Lanched SC2 (%v), PID: %v\n", processSettings.processPath, pi.PID)
	}

	c.SetProcessInfo(pi)
	processSettings.processInfo = append(processSettings.processInfo, pi)
	return pi.Port
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
		fmt.Println(err)
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
			panic("Failed to connect")
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

// Run ...
func Run() {
	wg := sync.WaitGroup{}
	wg.Add(len(clients))

	for _, c := range clients {
		go func(client *client.Client) {
			defer func() {
				if p := recover(); p != nil {
					fmt.Printf("Panic: %v\n", p)
				}
				wg.Done()
			}()

			runAgent(client)
		}(c)
	}

	wg.Wait()
}

func runAgent(client *client.Client) {
	stepSize := processSettings.stepSize
	if processSettings.realtime {
		stepSize = 0
	}

	err := client.Init() // get GameInfo, Data, and Observation
	agent := client.Agent

	if err == nil && client.IsInGame() {
		agent.OnGameStart(client)
		for {
			err := client.Update(stepSize)

			if err != nil || !client.IsInGame() {
				break
			}

			agent.OnStep()
		}
	}

	agent.OnGameEnd()
	client.RequestLeaveGame()

	for _, player := range client.Observation().GetPlayerResult() {
		if player.GetPlayerId() == client.PlayerID() {
			fmt.Println(player.GetResult())
		}
	}

	return
}

// SetupPorts ...
func SetupPorts(numAgents int, startPort int, checkSingle bool) {
	humans := numAgents
	if checkSingle {
		humans = 0
		for _, p := range gameSettings.playerSetup {
			if p.PlayerType == api.PlayerType_Participant {
				humans++
			}
		}
	}

	if humans > 1 {
		var ports = gameSettings.ports
		ports.SharedPort = int32(startPort + 1)
		ports.ServerPorts.GamePort = int32(startPort + 2)
		ports.ServerPorts.BasePort = int32(startPort + 3)

		for i := 0; i < numAgents; i++ {
			var base = int32(startPort + 4 + i*2)
			ports.ClientPorts = append(ports.ClientPorts, &api.PortSet{GamePort: base, BasePort: base + 1})
		}
		gameSettings.ports = ports
	}
}
