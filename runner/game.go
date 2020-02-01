package runner

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

type gameConfig struct {
	netAddress  string
	processInfo []client.ProcessInfo
	playerSetup []*api.PlayerSetup
	ports       client.Ports

	clients  []*client.Client
	started  bool
	lastPort int
}

func newGameConfig(participants ...client.PlayerSetup) *gameConfig {
	config := &gameConfig{
		"127.0.0.1",
		nil,
		nil,
		client.Ports{},
		nil,
		false,
		0,
	}

	for _, p := range participants {
		if p.Agent != nil {
			config.clients = append(config.clients, &client.Client{Agent: p.Agent})
		}
		config.playerSetup = append(config.playerSetup, p.PlayerSetup)
	}
	return config
}

func (config *gameConfig) startGame(mapPath string) {
	if !config.createGame(mapPath) {
		log.Fatal("Failed to create game.")
	}
	config.joinGame()
}

func (config *gameConfig) createGame(mapPath string) bool {
	if !config.started {
		log.Panic("Game not started")
	}

	// Create with the first client
	err := config.clients[0].CreateGame(mapPath, config.playerSetup, processRealtime)
	if err != nil {
		log.Print(err)
		return false
	}
	return true
}

func (config *gameConfig) joinGame() bool {
	// TODO: Make this parallel and get rid of the WaitJoinGame method
	for i, client := range config.clients {
		if err := client.RequestJoinGame(config.playerSetup[i], processInterfaceOptions, config.ports); err != nil {
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

func (config *gameConfig) connect(port int) {
	pi := client.ProcessInfo{Path: "", PID: 0, Port: port}

	// Set process info for each bot
	for range config.clients {
		config.processInfo = append(config.processInfo, pi)
	}

	// Since connect is blocking do it after the processes are launched.
	for i, client := range config.clients {
		pi := config.processInfo[i]

		if err := client.Connect(config.netAddress, pi.Port, processConnectTimeout); err != nil {
			log.Panic("Failed to connect")
		}
	}

	// Assume starcraft has started after succesfully attaching to a server
	config.started = true
}

func (config *gameConfig) setupPorts(numAgents int, startPort int, checkSingle bool) {
	humans := numAgents
	if checkSingle {
		humans = 0
		for _, p := range config.playerSetup {
			if p.Type == api.PlayerType_Participant {
				humans++
			}
		}
	}

	if humans > 1 {
		var ports = config.ports
		ports.SharedPort = int32(startPort + 1)
		ports.ServerPorts = &api.PortSet{
			GamePort: int32(startPort + 2),
			BasePort: int32(startPort + 3),
		}

		for i := 0; i < numAgents; i++ {
			var base = int32(startPort + 4 + i*2)
			ports.ClientPorts = append(ports.ClientPorts, &api.PortSet{GamePort: base, BasePort: base + 1})
		}
		config.ports = ports
	}
}
