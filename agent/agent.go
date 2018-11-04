package agent

import (
	"fmt"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

type AgentFunc func(Agent)

func (f AgentFunc) RunAgent(info client.AgentInfo) {
	f(Agent{info: info})
}

type Agent struct {
	info    client.AgentInfo
	actions []*api.Action
	done    bool
	// PlayerID() api.PlayerID
	// GameInfo() *api.ResponseGameInfo
	// Data() *api.ResponseData
	// Observation() *api.ResponseObservation
	// Query(query api.RequestQuery) func() *api.ResponseQuery
	// SendActions(actions []*api.Action) func() []api.ActionResult
	// SendObserverActions(obsActions []*api.ObserverAction)
	// SendDebugCommands(commands []*api.DebugCommand)
	// LeaveGame()
	// IsInGame() bool
	// Update(stepSize int) ([]api.UpgradeID, error)
}

func (a *Agent) Info() client.AgentInfo {
	return a.info
}

func (a *Agent) IsInGame() bool {
	return !a.done && a.info.IsInGame()
}

func (a *Agent) Update(stepSize int) []api.UpgradeID {
	a.sendActions()
	upgrades, err := a.info.Update(stepSize)
	if err != nil {
		fmt.Println(err)
		a.done = true
		return nil
	}
	return upgrades
}

func (a *Agent) LeaveGame() {
	a.info.LeaveGame()
}
