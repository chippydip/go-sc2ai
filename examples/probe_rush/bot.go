package main

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/protoss"
	"github.com/chippydip/go-sc2ai/enums/zerg"
)

type bot struct {
	client.AgentInfo

	*botutil.Player
	Units *botutil.Units
	*botutil.Actions
	*botutil.Builder

	enemyStartLocation api.Point2D
	homeMineral        botutil.Unit

	done bool
}

func runAgent(info client.AgentInfo) {
	bot := bot{AgentInfo: info}

	bot.Player = botutil.NewPlayer(info)
	bot.Actions = botutil.NewActions(info)
	bot.Units = botutil.NewUnits(info)
	bot.Builder = botutil.NewBuilder(info, bot.Player, bot.Units, bot.Actions)

	bot.init()
	for bot.IsInGame() {
		bot.strategy()
		bot.tactics()

		if err := bot.Step(1); err != nil {
			log.Print(err)
			break
		}
	}
}

func (bot *bot) init() {
	// Get the default attack target
	bot.enemyStartLocation = *bot.GameInfo().StartRaw.StartLocations[0]

	// Pick a mineral patch to retreat to (mineral-walk)
	nexusPos := bot.Units.First(botutil.IsSelfType(protoss.Nexus)).Pos.ToPoint2D()
	bot.homeMineral = bot.Units.ClosestWithFilter(nexusPos, botutil.IsMineral)

	// Send a friendly hello
	bot.Chat("(glhf)")
}

func (bot *bot) strategy() {
	// Build probe if we can
	bot.BuildUnit(protoss.Nexus, ability.Train_Probe)

	// Chronoboost self if building something
	nexus := bot.Units.First(botutil.IsSelfType(protoss.Nexus))
	if len(nexus.Orders) > 0 && nexus.Energy >= 50 {
		bot.UnitCommandOnTarget(nexus, ability.Effect_ChronoBoostEnergyCost, nexus)
	}
}

func (bot *bot) tactics() {
	// Make sure we still have some probes left
	probes := bot.Units.Choose(botutil.IsSelfType(protoss.Probe))
	if probes.Len() == 0 {
		if !bot.done {
			bot.Chat("(gg)") // we lose
			bot.done = true
		}
		return
	}

	// Get enemy units to target
	targets, _ := bot.getTargets()

	if targets.Len() == 0 {
		// Attack enemy base position
		bot.UnitsCommandAtPos(&probes, ability.Attack, &bot.enemyStartLocation)
		return
	}

	probes.Each(func(probe botutil.Unit) {
		if probe.Shield == 0 {
			// Mineral walk retreat until shields start to recharge
			bot.UnitCommandOnTarget(probe, ability.Harvest_Gather, bot.homeMineral)
		} else {
			// Attack the location of the closest unit
			target := targets.Closest(probe.Pos.ToPoint2D())
			pos := target.Pos.ToPoint2D()
			bot.UnitCommandAtPos(probe, ability.Attack, &pos)
		}
	})
}

// Get the current target list, prioritizing good targets over ok targets
func (bot *bot) getTargets() (botutil.Units, bool) {
	// OK targets are anything that's not flying or a zerg larva/egg
	ok := bot.Units.Choose(func(u botutil.Unit) bool {
		return u.IsEnemy() && !u.IsFlying && !u.IsAnyType(zerg.Larva, zerg.Egg, protoss.AdeptPhaseShift)
	})

	// Good targets are OK targets that aren't structures
	good := ok.Drop(botutil.IsStructure)

	if good.Len() > 0 {
		return good, true
	}
	return ok, false
}
