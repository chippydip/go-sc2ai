package botutil_test

import "github.com/chippydip/go-sc2ai/api"

type mockAgentInfo struct{}

func (a *mockAgentInfo) IsRealtime() bool {
	return false
}

func (a *mockAgentInfo) PlayerID() api.PlayerID {
	panic("Not Implemented")
}
func (a *mockAgentInfo) GameInfo() *api.ResponseGameInfo {
	panic("Not Implemented")
}
func (a *mockAgentInfo) ReplayInfo() *api.ResponseReplayInfo {
	panic("Not Implemented")
}
func (a *mockAgentInfo) Data() *api.ResponseData {
	panic("Not Implemented")
}
func (a *mockAgentInfo) Observation() *api.ResponseObservation {
	panic("Not Implemented")
}
func (a *mockAgentInfo) Upgrades() []api.UpgradeID {
	return nil
}
func (a *mockAgentInfo) HasUpgrade(upgrade api.UpgradeID) bool {
	return false
}

func (a *mockAgentInfo) IsInGame() bool {
	return true
}
func (a *mockAgentInfo) Step(stepSize int) error {
	return nil
}

func (a *mockAgentInfo) Query(query api.RequestQuery) *api.ResponseQuery {
	panic("Not Implemented")
}
func (a *mockAgentInfo) SendActions(actions []*api.Action) []api.ActionResult {
	panic("Not Implemented")
}
func (a *mockAgentInfo) SendObserverActions(obsActions []*api.ObserverAction) {
}
func (a *mockAgentInfo) SendDebugCommands(commands []*api.DebugCommand) {
}
func (a *mockAgentInfo) LeaveGame() {
}

func (a *mockAgentInfo) OnBeforeStep(func()) {
}
func (a *mockAgentInfo) OnSubStep(func()) {
}
func (a *mockAgentInfo) OnAfterStep(func()) {
}
