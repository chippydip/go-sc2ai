package client

import (
	"fmt"
	"log"
	"strings"
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
	replayInfo  *api.ResponseReplayInfo
	data        *api.ResponseData
	observation *api.ResponseObservation
	upgrades    map[api.UpgradeID]struct{}
	newUpgrades []api.UpgradeID

	beforeStep []func()
	subStep    []func()
	afterStep  []func()

	debugDraw chan struct{}

	perfInterval uint32
	lastDraw     []*api.DebugCommand

	perfStart       time.Time
	perfStartFrame  uint32
	beforeStepTime  time.Duration
	stepTime        time.Duration
	observationTime time.Duration
	afterStepTime   time.Duration

	actions          int
	maxActions       int
	actionsCompleted int
	observerActions  int
	debugCommands    int
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

// RequestReplayInfo ...
func (c *Client) RequestReplayInfo(path string) (*api.ResponseReplayInfo, error) {
	r, err := c.connection.replayInfo(api.RequestReplayInfo{
		Replay: &api.RequestReplayInfo_ReplayPath{
			ReplayPath: path,
		},
		DownloadData: true,
	})
	if err != nil {
		return nil, err
	}
	if r.Error != api.ResponseReplayInfo_nil {
		return nil, fmt.Errorf("%v: %v", r.Error.String(), r.GetErrorDetails())
	}
	return r, nil
}

// Proto ...
func (c *Client) Proto() api.ResponsePing {
	return c.connection.ResponsePing
}

// RequestStartReplay ...
func (c *Client) RequestStartReplay(request api.RequestStartReplay) error {
	c.replayInfo = nil

	r, err := c.connection.startReplay(request)
	if err != nil {
		return err
	}
	if r.Error != api.ResponseStartReplay_nil {
		return fmt.Errorf("%v: %v", r.Error.String(), r.GetErrorDetails())
	}

	c.replayInfo, err = c.RequestReplayInfo(request.GetReplayPath())
	if err != nil {
		log.Print(err)
	}
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

	c.perfStart = time.Now()
	c.perfStartFrame = c.observation.GetObservation().GetGameLoop()

	// This info isn't provided for replays, so try to normalize things
	if c.replayInfo != nil {
		c.gameInfo.MapName = c.replayInfo.MapName
		c.gameInfo.LocalMapPath = c.replayInfo.LocalMapPath
		c.gameInfo.PlayerInfo = make([]*api.PlayerInfo, len(c.replayInfo.PlayerInfo))
		for i, pie := range c.replayInfo.PlayerInfo {
			c.gameInfo.PlayerInfo[i] = pie.PlayerInfo
		}
	}

	return firstOrNil(infoErr, dataErr, obsErr)
}

// Step ...
func (c *Client) Step(stepSize int) error {
	var err error

	// Call before callbacks
	t := time.Now()
	for _, cb := range c.beforeStep {
		cb()
	}
	c.beforeStepTime += time.Since(t)

	// Step the simulation forward if this isn't in realtime mode
	t = time.Now()
	if !c.realtime && stepSize > 0 {
		if _, err := c.connection.step(api.RequestStep{
			Count: uint32(stepSize),
		}); err != nil {
			return err
		}
	}
	c.stepTime += time.Since(t)

	// Get an updated observation
	t = time.Now()
	step := c.observation.GetObservation().GetGameLoop() + uint32(stepSize)
	for {
		if c.observation, err = c.connection.observation(api.RequestObservation{GameLoop: step}); err != nil {
			return err
		}

		actionsCompleted := len(c.observation.GetActions())
		c.actionsCompleted += actionsCompleted
		if actionsCompleted > c.maxActions {
			c.maxActions = actionsCompleted
		}

		if !c.IsInGame() {
			// Clear draw commands in case the game is left running
			c.ClearDebugDraw()
			return nil
		}

		// Call sub-step callbacks
		for _, cb := range c.subStep {
			cb()
		}

		if c.observation.GetObservation().GetGameLoop() >= step {
			break
		}
	}
	c.observationTime += time.Since(t)

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
	t = time.Now()
	for _, cb := range c.afterStep {
		cb()
	}
	c.afterStepTime += time.Since(t)

	// Performance reporting (update every perfInterval game frames)
	if c.perfInterval > 0 && c.observation.GetObservation().GetGameLoop()%c.perfInterval == 0 {
		c.reportPerf()
	}
	return err
}

func (c *Client) reportPerf() {
	perfStart, perfStartFrame := time.Now(), c.observation.GetObservation().GetGameLoop()
	total, frames := perfStart.Sub(c.perfStart), time.Duration(perfStartFrame-c.perfStartFrame)

	text := "" +
		fmt.Sprintf("frames:      %v\n", int(frames)) +
		fmt.Sprintf("frameTime:   %v\n", total/frames) +
		"\n" +
		fmt.Sprintf("beforeStep:  %v\n", c.beforeStepTime/frames) +
		fmt.Sprintf("step:        %v\n", c.stepTime/frames) +
		fmt.Sprintf("observation: %v\n", c.observationTime/frames) +
		fmt.Sprintf("afterStep:   %v\n", c.afterStepTime/frames) +
		"\n" +
		fmt.Sprintf("actions:     %v/%v\n", c.actionsCompleted, c.actions) +
		fmt.Sprintf("maxActions:  %v\n", c.maxActions) +
		fmt.Sprintf("obsActions:  %v\n", c.observerActions) +
		fmt.Sprintf("debugCmds:   %v\n", c.debugCommands) +
		""
	text = strings.Replace(text, "µ", "u", -1)

	// Reset perf counters
	c.perfStart = perfStart
	c.perfStartFrame = perfStartFrame
	c.beforeStepTime = 0
	c.stepTime = 0
	c.observationTime = 0
	c.afterStepTime = 0
	c.actions = 0
	c.maxActions = 0
	c.actionsCompleted = 0
	c.observerActions = 0
	c.debugCommands = 0

	c.SendDebugCommands(append(c.lastDraw, &api.DebugCommand{
		Command: &api.DebugCommand_Draw{
			Draw: &api.DebugDraw{
				Text: []*api.DebugText{
					&api.DebugText{
						Color:      &api.Color{R: 255, G: 255, B: 255},
						Text:       text,
						VirtualPos: &api.Point{X: 0, Y: 0, Z: 0},
					},
				},
			},
		},
	}))
	c.lastDraw = c.lastDraw[:len(c.lastDraw)-1]
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
