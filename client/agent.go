package client

import (
	"log"

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
	LeaveGame()

	OnBeforeStep(func())
	OnSubStep(func())
	OnAfterStep(func())
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
	c.connection.obsAction(api.RequestObserverAction{
		Actions: obsActions,
	})
}

// SendDebugCommands ...
func (c *Client) SendDebugCommands(commands []*api.DebugCommand) {
	c.connection.debug(api.RequestDebug{
		Debug: commands,
	})
}

// LeaveGame ...
func (c *Client) LeaveGame() {
	c.connection.leaveGame(api.RequestLeaveGame{})
}

// OnBeforeStep ...
func (c *Client) OnBeforeStep(callback func()) {
	if callback != nil {
		c.beforeStep = append(c.beforeStep, callback)
	}
}

// OnSubStep is called after every observation. This is generally equivalent
// to OnAfterStep except in realtime mode when multiple observations may occur
// to reach the desired step size. This can be used to observe transient data
// that only appears in a single observation (actions, events, chat, etc).
func (c *Client) OnSubStep(callback func()) {
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
