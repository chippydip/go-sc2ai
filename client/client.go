package client

import (
	"fmt"
	"time"

	"github.com/chippydip/go-sc2ai/api"
)

// Client ...
type Client struct {
	connection
	Agent Agent

	playerID    api.PlayerID
	gameInfo    *api.ResponseGameInfo
	data        *api.ResponseData
	observation *api.ResponseObservation
	upgrades    map[api.UpgradeID]struct{}
}

// Connect ...
func (c *Client) Connect(address string, port int, timeout time.Duration) error {
	attempts := (int(timeout.Seconds()*1000) + 1500) / 1000
	if attempts < 1 {
		attempts = 1
	}

	connected := false
	for i := 0; i < attempts; i++ {
		if err := c.connection.Connect(address, port, timeout); err == nil {
			connected = true
			break
		}
		time.Sleep(time.Second)

		if i == 0 {
			fmt.Print("Waiting for connection")
		} else {
			fmt.Print(".")
		}
	}
	fmt.Println()

	if !connected {
		return fmt.Errorf("Unable to connect to game")
	}

	fmt.Printf("Connected to %v:%v\n", address, port)
	return nil
}

// RemoteSaveMap(data []byte, remotePath string) error

// CreateGame ...
func (c *Client) CreateGame(mapPath string, players []PlayerSetup, realtime bool) error {
	playerSetup := make([]*api.PlayerSetup, len(players))
	for i, p := range players {
		playerType, race, difficulty := p.PlayerType, p.Race, p.Difficulty
		playerSetup[i] = &api.PlayerSetup{
			Type:       playerType,
			Race:       race,
			Difficulty: difficulty,
		}
	}
	r, err := c.connection.createGame(api.RequestCreateGame{
		Map: &api.RequestCreateGame_LocalMap{
			LocalMap: &api.LocalMap{
				MapPath: mapPath,
			},
		},
		PlayerSetup: playerSetup,
		Realtime:    realtime,
	})()
	if err != nil {
		return err
	}

	if r.Error != api.ResponseCreateGame_nil {
		return fmt.Errorf("%v: %v", r.Error.String(), r.GetErrorDetails())
	}

	return nil
}

// RequestJoinGame ...
func (c *Client) RequestJoinGame(setup PlayerSetup, options *api.InterfaceOptions, ports Ports) error {
	req := api.RequestJoinGame{
		Participation: &api.RequestJoinGame_Race{
			Race: setup.Race,
		},
		Options: options,
	}
	if ports.isValid() {
		req.SharedPort = ports.SharedPort
		req.ServerPorts = ports.ServerPorts
		req.ClientPorts = ports.ClientPorts
	}
	r, err := c.connection.joinGame(req)()
	if err != nil {
		return err
	}

	if r.Error != api.ResponseJoinGame_nil {
		return fmt.Errorf("%v: %v", r.Error.String(), r.GetErrorDetails())
	}

	c.playerID = r.GetPlayerId()
	return nil
}

// RequestLeaveGame ...
func (c *Client) RequestLeaveGame() error {
	_, err := c.connection.leaveGame(api.RequestLeaveGame{})()
	return err
}

// Init ...
func (c *Client) Init() error {
	// Fire off all three requests
	infoReq := c.connection.gameInfo(api.RequestGameInfo{})
	dataReq := c.connection.data(api.RequestData{
		AbilityId:  true,
		UnitTypeId: true,
		UpgradeId:  true,
		BuffId:     true,
		EffectId:   true,
	})
	obsReq := c.connection.observation(api.RequestObservation{})

	// Wait for completion
	var infoErr, dataErr, obsErr error
	c.gameInfo, infoErr = infoReq()
	c.data, dataErr = dataReq()
	c.observation, obsErr = obsReq()
	return firstOrNil(infoErr, dataErr, obsErr)
}

// Update ...
func (c *Client) Update(stepSize int) ([]api.UpgradeID, error) {
	// Step the simulation forward if this isn't in realtime mode
	stepReq := func() error { return nil }
	if stepSize > 0 {
		f := c.connection.step(api.RequestStep{
			Count: uint32(stepSize),
		})
		stepReq = func() error { _, err := f(); return err }
	}

	// Get an updated observation
	obsReq := c.connection.observation(api.RequestObservation{})

	// Wait for completion
	var stepErr, obsErr error
	stepErr = stepReq()
	c.observation, obsErr = obsReq()
	if err := firstOrNil(stepErr, obsErr); err != nil {
		return nil, err
	}

	// Check for new upgrades
	var newUpgrades []api.UpgradeID
	for _, upgrade := range c.observation.GetObservation().GetRawData().GetPlayer().UpgradeIds {
		if _, ok := c.upgrades[upgrade]; !ok {
			newUpgrades = append(newUpgrades, upgrade)
			c.upgrades[upgrade] = struct{}{}
		}
	}
	if len(newUpgrades) == 0 {
		return nil, nil
	}

	// Re-fetch unit data since some of it is upgrade-dependent
	// TODO: also (re-)fetch unit -> ability mapping?
	data, err := c.connection.data(api.RequestData{
		UnitTypeId: true,
	})()
	c.data.Units = data.GetUnits()
	return newUpgrades, err
}

// SaveReplay(path string) error

// Pint() error

// // General
// WaitForResponse() (*GameResponse, error)

// SetProcessInfo ...
func (c *Client) SetProcessInfo(pi ProcessInfo) {
	// TODO
}

// GetProcessInfo() ProcessInfo

// GetAppState() AppState

// GetLastStatus() api.Status

// IsInGame ...
func (c *Client) IsInGame() bool {
	return c.connection.Status == api.Status_in_game || c.connection.Status == api.Status_in_replay
}

// IsFinishedGame() bool
// HasResponsePending() bool

// GetObservation ...
func (c *Client) GetObservation() (*api.ResponseObservation, error) {
	return c.connection.observation(api.RequestObservation{})()
}

// PollResponse() bool
// ConsumeResponse() bool

// IssueEvents(commands []Tag) bool
// OnGameStart()

// // Diagnostic
// DumpProtoUsage()

// Error(err ClientError, errors []string)
// ErrorIf(condition bool, err ClientError, errors []string)

// func (c *Client) GetClientErrors() []ClientError {
// 	return nil
// }

// func (c *Client) GetProtocolErrors() []string {
// 	return nil
// }

// ClearClientErrors()
// ClearProtocolErrors()

// UseGeneralizedAbility(value bool)

// // Save/Load
// Save()
// Load()

func firstOrNil(errs ...error) error {
	for _, e := range errs {
		if e != nil {
			return e
		}
	}
	return nil
}
