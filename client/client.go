package client

import (
	"fmt"
	"log"
	"time"

	"github.com/chippydip/go-sc2ai/api"
)

// Client ...
type Client struct {
	connection
	Agent    Agent
	realtime bool

	playerID    api.PlayerID
	gameInfo    *api.ResponseGameInfo
	data        *api.ResponseData
	observation *api.ResponseObservation
	upgrades    map[api.UpgradeID]struct{}
	newUpgrades []api.UpgradeID

	beforeStep []func()
	afterStep  []func()
}

// Connect ...
func (c *Client) Connect(address string, port int, timeout time.Duration) error {
	attempts := int(timeout.Seconds() + 1.5)
	if attempts < 1 {
		attempts = 1
	}

	connected := false
	for i := 0; i < attempts; i++ {
		if err := c.connection.Connect(address, port); err == nil {
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

	log.Printf("Connected to %v:%v", address, port)
	return nil
}

// TryConnect ...
func (c *Client) TryConnect(address string, port int) error {
	if err := c.connection.Connect(address, port); err != nil {
		return err
	}

	log.Printf("Connected to %v:%v", address, port)
	return nil
}

// RemoteSaveMap(data []byte, remotePath string) error

// CreateGame ...
func (c *Client) CreateGame(mapPath string, players []*api.PlayerSetup, realtime bool) error {
	r, err := c.connection.createGame(api.RequestCreateGame{
		Map: &api.RequestCreateGame_LocalMap{
			LocalMap: &api.LocalMap{
				MapPath: mapPath,
			},
		},
		PlayerSetup: players,
		Realtime:    realtime,
	})
	if err != nil {
		return err
	}
	c.realtime = realtime

	if r.Error != api.ResponseCreateGame_nil {
		return fmt.Errorf("%v: %v", r.Error, r.GetErrorDetails())
	}

	return nil
}

// RequestJoinGame ...
func (c *Client) RequestJoinGame(setup *api.PlayerSetup, options *api.InterfaceOptions, ports Ports) error {
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
	r, err := c.connection.joinGame(req)
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
	_, err := c.connection.leaveGame(api.RequestLeaveGame{})
	return err
}

// Init ...
func (c *Client) Init() error {
	var infoErr, dataErr, obsErr error

	// Fire off all three requests
	c.gameInfo, infoErr = c.connection.gameInfo(api.RequestGameInfo{})
	c.data, dataErr = c.connection.data(api.RequestData{
		AbilityId:  true,
		UnitTypeId: true,
		UpgradeId:  true,
		BuffId:     true,
		EffectId:   true,
	})
	c.observation, obsErr = c.connection.observation(api.RequestObservation{})
	c.upgrades = map[api.UpgradeID]struct{}{}

	return firstOrNil(infoErr, dataErr, obsErr)
}

// Step ...
func (c *Client) Step(stepSize int) error {
	var err error

	// Call before callbacks
	for _, cb := range c.beforeStep {
		cb()
	}

	// Step the simulation forward if this isn't in realtime mode
	if !c.realtime && stepSize > 0 {
		if _, err := c.connection.step(api.RequestStep{
			Count: uint32(stepSize),
		}); err != nil {
			return err
		}
	}

	// Get an updated observation
	step := c.observation.GetObservation().GetGameLoop() + uint32(stepSize)
	for {
		if c.observation, err = c.connection.observation(api.RequestObservation{}); err != nil {
			return err
		}
		if c.observation.GetObservation().GetGameLoop() >= step {
			break
		}
	}

	// Check for new upgrades
	c.newUpgrades = nil
	for _, upgrade := range c.observation.GetObservation().GetRawData().GetPlayer().UpgradeIds {
		if _, ok := c.upgrades[upgrade]; !ok {
			c.newUpgrades = append(c.newUpgrades, upgrade)
			c.upgrades[upgrade] = struct{}{}
		}
	}

	if len(c.newUpgrades) > 0 {
		// Re-fetch unit data since some of it is upgrade-dependent
		// TODO: also (re-)fetch unit -> ability mapping?
		var data *api.ResponseData
		data, err = c.connection.data(api.RequestData{
			UnitTypeId: true,
		})
		c.data.Units = data.GetUnits()
	}

	// Call after callbacks
	for _, cb := range c.afterStep {
		cb()
	}
	return err
}

// SaveReplay(path string) error

// Print() error

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
	return c.connection.observation(api.RequestObservation{})
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
