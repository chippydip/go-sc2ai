package main

import (
	"github.com/chippydip/go-sc2ai/agent"
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/protoss"
	"github.com/chippydip/go-sc2ai/enums/zerg"
	"github.com/chippydip/go-sc2ai/filter"
	"github.com/chippydip/go-sc2ai/search"
)

type bot struct {
	agent.Agent

	enemyStartLocation api.Point2D
	homeMineral        *api.Unit

	done bool
}

func runAgent(a agent.Agent) {
	bot := bot{Agent: a}
	bot.init()

	for bot.IsInGame() {
		bot.strategy()
		bot.tactics()

		bot.Step(1)
	}
}

func (bot *bot) init() {
	// Get the default attack target
	bot.enemyStartLocation = *bot.Info().GameInfo().StartRaw.StartLocations[0]

	// Pick a mineral patch to retreat to (mineral-walk)
	nexusPos := bot.GetUnit(filter.IsSelfType(protoss.Nexus)).Pos.ToPoint2D()
	bot.homeMineral = bot.GetClosestUnit(nexusPos, filter.IsMineral)

	// Send a friendly hello
	bot.ChatAll("(glhf)")
}

func (bot *bot) strategy() {
	// Build probe if we can
	bot.BuildUnit(protoss.Nexus, ability.Train_Probe)

	// Chronoboost self if building something
	nexus := bot.GetUnit(filter.IsSelfType(protoss.Nexus))
	if len(nexus.Orders) > 0 && nexus.Energy >= 50 {
		bot.UnitCommandAtTarget(nexus.Tag, ability.Effect_ChronoBoostEnergyCost, nexus.Tag)
	}
}

func (bot *bot) tactics() {
	// Make sure we still have some probes left
	probes := bot.GetUnits(filter.IsSelfType(protoss.Probe))
	if len(probes) == 0 {
		if !bot.done {
			bot.ChatAll("(gg)") // we lose
			bot.done = true
		}
		return
	}

	// Get enemy units to target
	targets, _ := bot.getTargets()

	if len(targets) == 0 {
		// Attack enemy base position
		bot.UnitsCommandAtPos(probes.Tags(), ability.Attack, bot.enemyStartLocation)
		return
	}

	for _, probe := range probes {
		if probe.Shield == 0 {
			// Mineral walk retreat until shields start to recharge
			bot.UnitCommandAtTarget(probe.Tag, ability.Harvest_Gather, bot.homeMineral.Tag)
		} else {
			// Attack the location of the closest unit
			target := search.ClosestUnit(probe.Pos.ToPoint2D(), targets...)
			bot.UnitCommandAtPos(probe.Tag, ability.Attack, target.Pos.ToPoint2D())
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
