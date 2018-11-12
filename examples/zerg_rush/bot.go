package main

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/protoss"
	"github.com/chippydip/go-sc2ai/enums/zerg"
	"github.com/chippydip/go-sc2ai/search"
)

type bot struct {
	client.AgentInfo

	*botutil.Player
	Units *botutil.Units
	*botutil.Actions
	*botutil.Builder

	myStartLocation    api.Point2D
	myNaturalLocation  api.Point2D
	enemyStartLocation api.Point2D

	camera api.Point2D
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
	// My hatchery is on start position
	bot.myStartLocation = bot.Units.First(botutil.IsSelfType(zerg.Hatchery)).Pos.ToPoint2D()
	bot.enemyStartLocation = *bot.GameInfo().GetStartRaw().GetStartLocations()[0]
	bot.camera = bot.myStartLocation

	// Find natural location
	expansions := search.CalculateExpansionLocations(bot, false)
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
	pool := bot.Units.First(botutil.IsSelfType(zerg.SpawningPool))
	if pool.IsNil() {
		pos := bot.myStartLocation.Offset(bot.enemyStartLocation, 5)
		if !bot.BuildUnitAt(zerg.Drone, ability.Build_SpawningPool, pos) {
			return // save up
		}
	}

	hatches := bot.Units.CountSelfType(zerg.Hatchery)

	// Build overlords as needed (want at least 3 spare supply per hatch)
	if bot.FoodLeft() <= 3*hatches && bot.Units.CountSelfTypeInProduction(zerg.Overlord) == 0 {
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
	droneCount := bot.Units.CountSelfTypeAll(zerg.Drone)
	bot.BuildUnits(zerg.Larva, ability.Train_Drone, maxDrones-droneCount)

	// We need a pool before trying to build lings or queens
	if pool.IsNil() || pool.BuildProgress < 1 {
		return
	}

	// Spend any extra larva on zerglings
	bot.BuildUnits(zerg.Larva, ability.Train_Zergling, 100)

	// Get a queen for every hatch if we still have minerals
	bot.BuildUnits(zerg.Hatchery, ability.Train_Queen, hatches-bot.Units.CountSelfType(zerg.Queen))

	// Expand to natural (mostly just for the larva, but might as well put it in the right spot)
	if hatches < 2 {
		bot.BuildUnitAt(zerg.Drone, ability.Build_Hatchery, bot.myNaturalLocation)
	}
}

func (bot *bot) tactics() {
	// If a hatch needs an injection, find the closest queen with energy
	hatch := bot.Units.First(func(u botutil.Unit) bool {
		return botutil.IsSelfType(zerg.Hatchery)(u) && len(u.BuffIds) == 0 && u.BuildProgress == 1
	})
	if !hatch.IsNil() {
		queen := bot.Units.ClosestWithFilter(hatch.Pos.ToPoint2D(), func(u botutil.Unit) bool {
			return u.IsSelf() && u.IsType(zerg.Queen) && u.Energy >= 25
		})
		if !queen.IsNil() {
			bot.UnitCommandOnTarget(queen, ability.Effect_InjectLarva, hatch)
		}
	}

	lings := bot.Units.Choose(botutil.IsSelfType(zerg.Zergling))
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

	targets, _ := bot.getTargets()
	if targets.Len() == 0 {
		bot.UnitsCommandAtPos(&lings, ability.Attack, &bot.enemyStartLocation)
		return
	}

	lings.Each(func(ling botutil.Unit) {
		target := targets.Closest(ling.Pos.ToPoint2D())
		if ling.Pos.ToPoint2D().Distance2(target.Pos.ToPoint2D()) > 4*4 {
			// If target is far, attack it as unit, ling will run ignoring everything else
			bot.UnitCommandOnTarget(ling, ability.Attack, target)
		} else if target.UnitType == zerg.ChangelingZergling || target.UnitType == zerg.ChangelingZerglingWings {
			// Must specificially attack changelings, attack move is not enough
			bot.UnitCommandOnTarget(ling, ability.Attack, target)
		} else {
			// Attack as position, ling will choose best target around
			pos := target.Pos.ToPoint2D()
			bot.UnitCommandAtPos(ling, ability.Attack, &pos)
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
