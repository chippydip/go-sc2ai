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
	PlayerID() api.PlayerID
	GameInfo() *api.ResponseGameInfo
	Data() *api.ResponseData
	Observation() *api.ResponseObservation
	Query(query api.RequestQuery) *api.ResponseQuery
	SendActions(actions []*api.Action) []api.ActionResult
	SendObserverActions(obsActions []*api.ObserverAction)
	SendDebugCommands(commands []*api.DebugCommand)
	LeaveGame()

	IsInGame() bool
	Update(stepSize int) ([]api.UpgradeID, error)
}

// PlayerID ...
func (c *Client) PlayerID() api.PlayerID {
	return c.playerID
}

// GameInfo ...
func (c *Client) GameInfo() *api.ResponseGameInfo {
	return c.gameInfo
}

// Data ...
func (c *Client) Data() *api.ResponseData {
	return c.data
}

// Observation ...
func (c *Client) Observation() *api.ResponseObservation {
	return c.observation
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
