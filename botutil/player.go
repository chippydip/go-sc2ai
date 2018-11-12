package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

type Player api.PlayerCommon

func NewPlayer(info client.AgentInfo) *Player {
	p := &Player{}
	update := func() {
		if pc := info.Observation().GetObservation().GetPlayerCommon(); pc != nil {
			*p = Player(*pc)
		}
	}
	update()
	info.OnAfterStep(update)
	return p
}

func (p *Player) FoodLeft() int {
	return int(p.FoodCap) - int(p.FoodUsed)
}
