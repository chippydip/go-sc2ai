package agent

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

type AgentFunc func(Agent)

func (f AgentFunc) RunAgent(info client.AgentInfo) {
	a := Agent{info: info}
	a.updateFood()
	f(a)
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

func (a *Agent) Step(stepSize int) {
	a.sendActions()
	err := a.info.Step(stepSize)
	a.updateFood()
	if err != nil {
		log.Print(err)
		a.done = true
	}
}

func (a *Agent) LeaveGame() {
	a.info.LeaveGame()
}

func (a *Agent) updateFood() {
	// TODO: skip out if player is not zerg?
	n, data := 0, a.Info().Data().Units
	for _, u := range a.GetAllUnits() {
		// Count number of units that consume half a food
		if u.Alliance == api.Alliance_Self && data[u.UnitType].FoodRequired == 0.5 {
			n++
		}
	}
	// The game rounds fractional food down, but should really round up since you
	// this makes it seem like you can build units when you actually can't.
	if n%2 != 0 {
		a.playerCommon().FoodUsed++
	}
}
