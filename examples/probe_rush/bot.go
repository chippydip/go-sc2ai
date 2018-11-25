package main

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/protoss"
)

type bot struct {
	botutil.Bot

	enemyStartLocation api.Point2D
	homeMineral        botutil.Unit

	done bool
}

func runAgent(info client.AgentInfo) {
	bot := bot{Bot: botutil.NewBot(info)}

	bot.init()
	for bot.IsInGame() && !bot.done {
		bot.update()

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
	nexusPos := bot.Self[protoss.Nexus].First().Pos2D()
	bot.homeMineral = bot.Neutral.Minerals().ClosestTo(nexusPos)

	// Send a friendly hello
	bot.Chat("(glhf)")
}

func (bot *bot) update() {
	// Build probe if we can
	bot.BuildUnit(protoss.Nexus, ability.Train_Probe)

	// Chronoboost self if building something
	nexus := bot.Self[protoss.Nexus].First()
	if len(nexus.Orders) > 0 && nexus.Energy >= 50 {
		nexus.OrderTarget(ability.Effect_ChronoBoostEnergyCost, nexus)
	}

	// Make sure we still have some probes left
	probes := bot.Self[protoss.Probe]
	if probes.Len() == 0 {
		bot.Chat("(gg)") // we lose
		bot.done = true
		return
	}

	// Get enemy units to target
	targets := bot.getTargets()

	if targets.Len() == 0 {
		// Attack enemy base position
		probes.OrderPos(ability.Attack, &bot.enemyStartLocation)
		return
	}

	probes.Each(func(probe botutil.Unit) {
		if probe.Shield == 0 {
			// Mineral walk retreat until shields start to recharge
			probe.OrderTarget(ability.Harvest_Gather, bot.homeMineral)
		} else {
			// Attack the location of the closest unit
			target := targets.ClosestTo(probe.Pos2D())
			pos := target.Pos2D()
			probe.OrderPos(ability.Attack, &pos)
		}
	})
}

// Get the current target list, prioritizing good targets over ok targets
func (bot *bot) getTargets() botutil.Units {
	// Prioritize things that can fight back
	if targets := bot.Enemy.Ground().CanAttack().All(); targets.Len() > 0 {
		return targets
	}

	// Otherwise just kill all the buildings
	return bot.Enemy.Ground().Structures().All()
}
