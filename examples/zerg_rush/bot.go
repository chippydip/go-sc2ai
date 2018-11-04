package main

import (
	"github.com/chippydip/go-sc2ai/agent"
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/zerg"
	"github.com/chippydip/go-sc2ai/filter"
	"github.com/chippydip/go-sc2ai/search"
)

type bot struct {
	agent.Agent

	myStartLocation    api.Point2D
	enemyStartLocation api.Point2D
}

func runAgent(a agent.Agent) {
	bot := bot{Agent: a}
	bot.init()

	for bot.IsInGame() {
		bot.strategy()
		bot.tactics()

		bot.Update(1)
	}
}

func (bot *bot) init() {
	// My hatchery is on start position
	bot.myStartLocation = bot.GetUnit(filter.IsSelfType(zerg.Hatchery)).Pos.ToPoint2D()
	bot.enemyStartLocation = *bot.Info().GameInfo().GetStartRaw().GetStartLocations()[0]

	// Send a friendly hello
	bot.ChatAll("(glhf)")
}

func (bot *bot) strategy() {
	// Do we have a pool? if not, try to build one
	pool := bot.GetUnit(filter.IsSelfType(zerg.SpawningPool))
	if pool == nil {
		pos := bot.myStartLocation.Offset(bot.enemyStartLocation, 5)
		if !bot.BuildUnitAt(zerg.Drone, ability.Build_SpawningPool, pos) {
			return // save up
		}
	}

	// Build overlords as needed
	foodLeft := bot.FoodCap() - bot.FoodUsed()
	if foodLeft < 2 && bot.CountUnitsInProduction(zerg.Egg, ability.Train_Overlord) == 0 {
		bot.BuildUnit(zerg.Larva, ability.Train_Overlord)
	}

	// Build drones up to 14
	droneCount := bot.CountUnitsAndProduction(zerg.Egg, ability.Train_Drone)
	if droneCount < 14 {
		bot.BuildUnits(zerg.Larva, ability.Train_Drone, 14-droneCount)
	}

	// Spend any extra larva on zerglings
	bot.BuildUnits(zerg.Larva, ability.Train_Zergling, 100)

	// If we run out of larva and still have minerals, train a Queen
	if bot.CountUnits(zerg.Queen) == 0 {
		bot.BuildUnit(zerg.Hatchery, ability.Train_Queen)
	}
}

func (bot *bot) tactics() {
	// If queen can inject, do it
	queen := bot.GetUnit(filter.IsSelfType(zerg.Queen))
	if queen != nil && queen.Energy >= 25 {
		hatch := bot.GetUnit(filter.IsSelfType(zerg.Hatchery))
		bot.UnitCommandAtTarget(queen.Tag, ability.Effect_InjectLarva, hatch.Tag)
	}

	lings := bot.GetUnits(filter.IsSelfType(zerg.Zergling))
	if len(lings) < 6 {
		return // wait for critical mass
	}

	targets, _ := bot.getTargets()
	if len(targets) == 0 {
		bot.UnitsCommandAtPos(lings.Tags(), ability.Attack, bot.enemyStartLocation)
		return
	}

	for _, ling := range lings {
		target := search.ClosestUnit(ling.Pos.ToPoint2D(), targets...)
		if ling.Pos.ToPoint2D().Distance2(target.Pos.ToPoint2D()) > 4*4 {
			// If target is far, attack it as unit, ling will run ignoring everything else
			bot.UnitCommandAtTarget(ling.Tag, ability.Attack, target.Tag)
		} else {
			// Attack as position, ling will choose best target around
			bot.UnitCommandAtPos(ling.Tag, ability.Attack, target.Pos.ToPoint2D())
		}
	}
}

// Get the current target list, prioritizing good targets over ok targets
func (bot *bot) getTargets() ([]*api.Unit, bool) {
	// OK targets are anything that's not flying or a zerg larva/egg
	ok := bot.GetUnits(func(u *api.Unit) bool {
		return u.Alliance == api.Alliance_Enemy && !u.IsFlying &&
			u.UnitType != zerg.Larva && u.UnitType != zerg.Egg
	})

	// Good targets are OK targets that aren't structures
	good := ok.Filter(func(u *api.Unit) bool {
		return !bot.UnitHasAttribute(u, api.Attribute_Structure)
	})

	if len(good) > 0 {
		return good, true
	}
	return ok, false
}
