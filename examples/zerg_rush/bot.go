package main

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/buff"
	"github.com/chippydip/go-sc2ai/enums/zerg"
	"github.com/chippydip/go-sc2ai/search"
)

type bot struct {
	*botutil.Bot

	myStartLocation    api.Point2D
	myNaturalLocation  api.Point2D
	enemyStartLocation api.Point2D

	camera api.Point2D
}

func runAgent(info client.AgentInfo) {
	bot := bot{Bot: botutil.NewBot(info)}

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
	// My hatchery is on start position
	bot.myStartLocation = bot.Self[zerg.Hatchery].First().Pos2D()
	bot.enemyStartLocation = *bot.GameInfo().GetStartRaw().GetStartLocations()[0]
	bot.camera = bot.myStartLocation

	// Find natural location
	expansions := search.CalculateExpansionLocations(bot.Bot, false)
	query := make([]*api.RequestQueryPathing, len(expansions))
	for i, exp := range expansions {
		pos := exp.Center()
		query[i] = &api.RequestQueryPathing{
			Start: &api.RequestQueryPathing_StartPos{
				StartPos: &bot.myStartLocation,
			},
			EndPos: &pos,
		}
	}
	resp := bot.Query(api.RequestQuery{Pathing: query})
	best, minDist := -1, float32(256)
	for i, result := range resp.GetPathing() {
		if result.Distance < minDist && result.Distance > 5 {
			best, minDist = i, result.Distance
		}
	}
	bot.myNaturalLocation = expansions[best].Center()

	// Send a friendly hello
	bot.Chat("(glhf)")
}

func (bot *bot) strategy() {
	// Do we have a pool? if not, try to build one
	pool := bot.Self[zerg.SpawningPool].First()
	if pool.IsNil() {
		pos := bot.myStartLocation.Offset(bot.enemyStartLocation, 5)
		if !bot.BuildUnitAt(zerg.Drone, ability.Build_SpawningPool, pos) {
			return // save up
		}
	}

	hatches := bot.Self.Count(zerg.Hatchery)

	// Build overlords as needed (want at least 3 spare supply per hatch)
	if bot.FoodLeft() <= 3*hatches && bot.Self.CountInProduction(zerg.Overlord) == 0 {
		if !bot.BuildUnit(zerg.Larva, ability.Train_Overlord) {
			return // save up
		}
	}

	// Any more than 14 drones will delay the first round of lings (waiting for larva)
	maxDrones := 14
	if hatches > 1 {
		maxDrones = 16 // but we can saturate later
	}

	// Build drones to our cap
	droneCount := bot.Self.CountAll(zerg.Drone)
	bot.BuildUnits(zerg.Larva, ability.Train_Drone, maxDrones-droneCount)

	// We need a pool before trying to build lings or queens
	if pool.IsNil() || pool.BuildProgress < 1 {
		return
	}

	// Spend any extra larva on zerglings
	bot.BuildUnits(zerg.Larva, ability.Train_Zergling, 100)

	// Get a queen for every hatch if we still have minerals
	bot.BuildUnits(zerg.Hatchery, ability.Train_Queen, hatches-bot.Self.CountAll(zerg.Queen))

	// Expand to natural (mostly just for the larva, but might as well put it in the right spot)
	if hatches < 2 {
		bot.BuildUnitAt(zerg.Drone, ability.Build_Hatchery, bot.myNaturalLocation)
	}
}

func (bot *bot) tactics() {
	// If a hatch needs an injection, find the closest queen with energy
	bot.Self[zerg.Hatchery].IsBuilt().NoBuff(buff.QueenSpawnLarvaTimer).Each(func(u botutil.Unit) {
		bot.Self[zerg.Queen].HasEnergy(25).ClosestTo(u.Pos2D()).OrderTarget(ability.Effect_InjectLarva, u)
	})

	lings := bot.Self[zerg.Zergling]
	if lings.Len() < 6 {
		return // wait for critical mass
	}

	// Auto-follow the action
	camera, minDist2 := bot.myStartLocation, bot.myStartLocation.Distance2(bot.enemyStartLocation)
	for _, c := range search.Cluster(lings, 16) {
		pos := c.Center()
		if dist2 := pos.Distance2(bot.enemyStartLocation); dist2 < minDist2 {
			camera, minDist2 = pos, dist2
		}
	}
	bot.camera = bot.camera.Add(bot.camera.VecTo(camera).Div(10))
	bot.MoveCamera(bot.camera)

	targets := bot.getTargets()
	if targets.Len() == 0 {
		lings.OrderPos(ability.Attack, &bot.enemyStartLocation)
		return
	}

	lings.Each(func(ling botutil.Unit) {
		target := targets.ClosestTo(ling.Pos2D())
		if ling.Pos2D().Distance2(target.Pos2D()) > 4*4 {
			// If target is far, attack it as unit, ling will run ignoring everything else
			ling.OrderTarget(ability.Attack, target)
		} else if target.UnitType == zerg.ChangelingZergling || target.UnitType == zerg.ChangelingZerglingWings {
			// Must specificially attack changelings, attack move is not enough
			ling.OrderTarget(ability.Attack, target)
		} else {
			// Attack as position, ling will choose best target around
			pos := target.Pos2D()
			ling.OrderPos(ability.Attack, &pos)
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
