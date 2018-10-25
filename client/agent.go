package client

import (
	"fmt"

	"github.com/chippydip/go-sc2ai/api"
)

// Agent ...
type Agent interface {
	OnGameStart(AgentInfo)
	OnStep()
	OnGameEnd()
}

// AgentInfo ...
type AgentInfo interface {
	PlayerID() uint32
	GameInfo() *api.ResponseGameInfo
	Data() *api.ResponseData
	Observation() *api.ResponseObservation
	Query(query api.RequestQuery) *api.ResponseQuery
	SendActions(actions []*api.Action) []api.ActionResult
	SendObserverActions(obsActions []*api.ObserverAction)
	SendDebugCommands(commands []*api.DebugCommand)
	LeaveGame()
}

// PlayerID ...
func (c *Client) PlayerID() uint32 {
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
	resp, err := c.connection.query(query)()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return resp
}

// SendActions ...
func (c *Client) SendActions(actions []*api.Action) []api.ActionResult {
	resp, err := c.connection.action(api.RequestAction{
		Actions: actions,
	})()
	if err != nil {
		fmt.Println(err)
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
