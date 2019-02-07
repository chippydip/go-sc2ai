package botutil

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

// Player ...
type Player struct {
	api.PlayerCommon
	api.PlayerInfo

	OpponentID   api.PlayerID
	OpponentRace api.Race
}

// NewPlayer ...
func NewPlayer(info client.AgentInfo) *Player {
	p := &Player{}
	for _, pi := range info.GameInfo().GetPlayerInfo() {
		if pi.GetPlayerId() == info.PlayerID() {
			p.PlayerInfo = *pi
		} else {
			p.OpponentID = pi.GetPlayerId()
			p.OpponentRace = pi.GetRaceRequested()
		}
	}
	update := func() {
		if pc := info.Observation().GetObservation().GetPlayerCommon(); pc != nil {
			p.PlayerCommon = *pc
		}

		if p.OpponentRace == api.Race_Random {
			for _, u := range info.Observation().GetObservation().GetRawData().GetUnits() {
				if u.GetOwner() == p.OpponentID {
					data := info.Data().GetUnits()[u.GetUnitType()]
					p.OpponentRace = data.GetRace()
					log.Printf("Detected OpponentRace: %v", p.OpponentRace)
					break
				}
			}
		}
	}
	update()
	info.OnAfterStep(update)
	return p
}

// FoodLeft returns the amount under (positive) or over (negative) the current food cap.
func (p *Player) FoodLeft() int {
	return int(p.FoodCap) - int(p.FoodUsed)
}

// Cost represents the full cost of an ability.
type Cost struct {
	Minerals, Vespene, Food uint32
}

// Mul multiplies the cost by the given count and returns a new Cost.
func (c Cost) Mul(count uint32) Cost {
	return Cost{
		Minerals: c.Minerals * count,
		Vespene:  c.Vespene * count,
		Food:     c.Food * count,
	}
}

// CanAfford determines if the player can currently afford the given cost.
func (p *Player) CanAfford(cost Cost) bool {
	return p.Minerals >= cost.Minerals && p.Vespene >= cost.Vespene &&
		(cost.Food == 0 || p.FoodCap >= p.FoodUsed+cost.Food)

}

// Spend tentatively marks the given resources as unavailable.
func (p *Player) Spend(cost Cost) {
	p.Minerals -= cost.Minerals
	p.Vespene -= cost.Vespene
	p.FoodUsed += cost.Food
}

// UpgradeCost returns the cost to research the given upgrade.
func (b *Bot) UpgradeCost(upgrade api.UpgradeID) Cost {
	data := b.Data().GetUpgrades()[upgrade]
	return Cost{
		Minerals: data.MineralCost,
		Vespene:  data.VespeneCost,
		Food:     0,
	}
}
