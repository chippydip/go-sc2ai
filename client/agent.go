package client

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/chippydip/go-sc2ai/api"
)

// Agent ...
type Agent interface {
	RunAgent(AgentInfo)
}

// AgentFunc ...
type AgentFunc func(AgentInfo)

// RunAgent ...
func (f AgentFunc) RunAgent(info AgentInfo) {
	f(info)
}

// AgentInfo ...
type AgentInfo interface {
	IsRealtime() bool

	PlayerID() api.PlayerID
	GameInfo() *api.ResponseGameInfo
	ReplayInfo() *api.ResponseReplayInfo
	Data() *api.ResponseData
	Observation() *api.ResponseObservation
	Upgrades() []api.UpgradeID
	HasUpgrade(upgrade api.UpgradeID) bool

	IsInGame() bool
	Step(stepSize int) error

	Query(query api.RequestQuery) *api.ResponseQuery
	SendActions(actions []*api.Action) []api.ActionResult
	SendObserverActions(obsActions []*api.ObserverAction)
	SendDebugCommands(commands []*api.DebugCommand)
	ClearDebugDraw()
	LeaveGame()
	SaveReplay(path string)

	OnBeforeStep(func())
	OnObservation(func())
	OnAfterStep(func())

	SetPerfInterval(steps uint32)
}

// IsRealtime returns true if the bot was launched in realtime mode.
func (c *Client) IsRealtime() bool {
	return c.realtime
}

// PlayerID ...
func (c *Client) PlayerID() api.PlayerID {
	return c.playerID
}

// GameInfo ...
func (c *Client) GameInfo() *api.ResponseGameInfo {
	return c.gameInfo
}

// ReplayInfo ...
func (c *Client) ReplayInfo() *api.ResponseReplayInfo {
	return c.replayInfo
}

// Data ...
func (c *Client) Data() *api.ResponseData {
	return c.data
}

// Observation ...
func (c *Client) Observation() *api.ResponseObservation {
	return c.observation
}

// Upgrades ...
func (c *Client) Upgrades() []api.UpgradeID {
	return c.newUpgrades
}

// HasUpgrade ...
func (c *Client) HasUpgrade(upgrade api.UpgradeID) bool {
	_, ok := c.upgrades[upgrade]
	return ok
}

// Query ...
func (c *Client) Query(query api.RequestQuery) *api.ResponseQuery {
	resp, err := c.connection.query(query)
	if err != nil {
		log.Print(err)
		return nil
	}
	return resp
}

// SendActions ...
func (c *Client) SendActions(actions []*api.Action) []api.ActionResult {
	c.actions += len(actions)

	if c.replayInfo != nil {
		return nil // ignore actions in a replay
	}

	resp, err := c.connection.action(api.RequestAction{
		Actions: actions,
	})
	if err != nil {
		log.Print(err)
		return nil
	}
	return resp.GetResult()
}

// SendObserverActions ...
func (c *Client) SendObserverActions(obsActions []*api.ObserverAction) {
	c.observerActions += len(obsActions)

	if c.replayInfo == nil {
		return // ignore observer actions in a normal game
	}

	c.connection.obsAction(api.RequestObserverAction{
		Actions: obsActions,
	})
}

// SendDebugCommands ...
func (c *Client) SendDebugCommands(commands []*api.DebugCommand) {
	c.debugCommands += len(commands)

	c.lastDraw = nil
	for _, cmd := range commands {
		if _, ok := cmd.Command.(*api.DebugCommand_Draw); ok {
			c.lastDraw = append(c.lastDraw, cmd)
		}
	}
	if c.debugDraw == nil && len(c.lastDraw) > 0 {
		c.debugDraw = deferCleanup(func() { c.ClearDebugDraw() })
	}

	c.connection.debug(api.RequestDebug{
		Debug: commands,
	})
}

func deferCleanup(cleanup func()) chan struct{} {
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
	cancel := make(chan struct{})

	go func() {
		select {
		case <-sig:
			cleanup()
			os.Exit(1)
		case <-cancel:
		}
	}()
	return cancel
}

// ClearDebugDraw sends an empty draw command if previous draw commands have been sent
func (c *Client) ClearDebugDraw() {
	if c.debugDraw == nil {
		return
	}
	close(c.debugDraw)
	c.debugDraw = nil
	c.lastDraw = nil

	c.SendDebugCommands([]*api.DebugCommand{
		&api.DebugCommand{
			Command: &api.DebugCommand_Draw{
				Draw: &api.DebugDraw{},
			},
		},
	})
}

// LeaveGame ...
func (c *Client) LeaveGame() {
	c.connection.leaveGame(api.RequestLeaveGame{})
}

// SaveReplay ...
func (c *Client) SaveReplay(path string) {
	responseSaveReplay, err := c.connection.saveReplay(api.RequestSaveReplay{})
	if err != nil {
		log.Print(err)
		return
	}
	err = os.WriteFile(path, responseSaveReplay.GetData(), 0644)
	if err != nil {
		log.Print(err)
	}
}

// OnBeforeStep ...
func (c *Client) OnBeforeStep(callback func()) {
	if callback != nil {
		c.beforeStep = append(c.beforeStep, callback)
	}
}

// OnObservation is called after every observation. This is generally equivalent
// to OnAfterStep except in realtime mode when multiple observations may occur
// to reach the desired step size. This can be used to observe transient data
// that only appears in a single observation (actions, events, chat, etc).
func (c *Client) OnObservation(callback func()) {
	if callback != nil {
		c.subStep = append(c.subStep, callback)
	}
}

// OnAfterStep ...
func (c *Client) OnAfterStep(callback func()) {
	if callback != nil {
		c.afterStep = append(c.afterStep, callback)
	}
}

// SetPerfInterval determines how often perfornace data will be updated. Values
// less than or equal to 0 will disable display (defalts to zero).
func (c *Client) SetPerfInterval(steps uint32) {
	c.perfInterval = steps
}
