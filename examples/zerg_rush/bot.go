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

	myStartLocation    api.Point2D
	myNaturalLocation  api.Point2D
	enemyStartLocation api.Point2D

	camera api.Point2D
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
	// My hatchery is on start position
	bot.myStartLocation = bot.GetUnit(filter.IsSelfType(zerg.Hatchery)).Pos.ToPoint2D()
	bot.enemyStartLocation = *bot.Info().GameInfo().GetStartRaw().GetStartLocations()[0]
	bot.camera = bot.myStartLocation

	// Find natural location
	expansions := search.CalculateExpansionLocations(bot.Info(), false)
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
	resp := bot.Info().Query(api.RequestQuery{Pathing: query})
	best, minDist := -1, float32(256)
	for i, result := range resp.GetPathing() {
		if result.Distance < minDist && result.Distance > 5 {
			best, minDist = i, result.Distance
		}
	}
	bot.myNaturalLocation = expansions[best].Center()

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

	hatches := bot.CountUnits(zerg.Hatchery)

	// Build overlords as needed (want at least 3 spare supply per hatch)
	foodLeft := bot.FoodCap() - bot.FoodUsed()
	if foodLeft <= 3*hatches && bot.CountUnitsInProduction(zerg.Egg, ability.Train_Overlord) == 0 {
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
	droneCount := bot.CountUnitsAndProduction(zerg.Egg, ability.Train_Drone)
	bot.BuildUnits(zerg.Larva, ability.Train_Drone, maxDrones-droneCount)

	// We need a pool before trying to build lings or queens
	if pool == nil || pool.BuildProgress < 1 {
		return
	}

	// Spend any extra larva on zerglings
	bot.BuildUnits(zerg.Larva, ability.Train_Zergling, 100)

	// Get a queen for every hatch if we still have minerals
	bot.BuildUnits(zerg.Hatchery, ability.Train_Queen, hatches-bot.CountUnits(zerg.Queen))

	// Expand to natural (mostly just for the larva, but might as well put it in the right spot)
	if hatches < 2 {
		bot.BuildUnitAt(zerg.Drone, ability.Build_Hatchery, bot.myNaturalLocation)
	}
}

func (bot *bot) tactics() {
	// If a hatch needs an injection, find the closest queen with energy
	hatch := bot.GetUnit(func(u *api.Unit) bool {
		return filter.IsSelfType(zerg.Hatchery)(u) && len(u.BuffIds) == 0 && u.BuildProgress == 1
	})
	if hatch != nil {
		queen := bot.GetClosestUnit(hatch.Pos.ToPoint2D(), func(u *api.Unit) bool {
			return filter.IsSelfType(zerg.Queen)(u) && u.Energy >= 25
		})
		if queen != nil {
			bot.UnitCommandAtTarget(queen.Tag, ability.Effect_InjectLarva, hatch.Tag)
		}
	}

	lings := bot.GetUnits(filter.IsSelfType(zerg.Zergling))
	if len(lings) < 6 {
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
	if len(targets) == 0 {
		bot.UnitsCommandAtPos(lings.Tags(), ability.Attack, bot.enemyStartLocation)
		return
	}

	for _, ling := range lings {
		target := search.ClosestUnit(ling.Pos.ToPoint2D(), targets...)
		if ling.Pos.ToPoint2D().Distance2(target.Pos.ToPoint2D()) > 4*4 {
			// If target is far, attack it as unit, ling will run ignoring everything else
			bot.UnitCommandAtTarget(ling.Tag, ability.Attack, target.Tag)
		} else if target.UnitType == zerg.ChangelingZergling || target.UnitType == zerg.ChangelingZerglingWings {
			// Must specificially attack changelings, attack move is not enough
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
			u.UnitType != zerg.Larva && u.UnitType != zerg.Egg && u.UnitType != protoss.AdeptPhaseShift
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
