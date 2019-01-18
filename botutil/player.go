package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

// Player ...
type Player struct {
	api.PlayerCommon
	api.PlayerInfo

	OpponentRace api.Race
}

// NewPlayer ...
func NewPlayer(info client.AgentInfo) *Player {
	p := &Player{}
	for _, pi := range info.GameInfo().GetPlayerInfo() {
		if pi.GetPlayerId() == info.PlayerID() {
			p.PlayerInfo = *pi
		} else {
			if p.OpponentRace == api.Race_NoRace {
				p.OpponentRace = pi.GetRaceRequested()
			} else if p.OpponentRace != pi.GetRaceRequested() {
				p.OpponentRace = api.Race_Random
			}
		}
	}
	update := func() {
		if pc := info.Observation().GetObservation().GetPlayerCommon(); pc != nil {
			p.PlayerCommon = *pc
		}
	}
	update()
	info.OnAfterStep(update)
	return p
}

// FoodLeft ...
func (p *Player) FoodLeft() int {
	return int(p.FoodCap) - int(p.FoodUsed)
}
