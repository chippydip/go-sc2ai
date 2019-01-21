package main

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/terran"
	"github.com/chippydip/go-sc2ai/enums/zerg"
	"github.com/chippydip/go-sc2ai/search"
)

type bot struct {
	*botutil.Bot

	myStartLocation    api.Point2D
	homeMineral        botutil.Unit
	enemyStartLocation api.Point2D
	baseLocations      []api.Point2D

	positionsForSupplies []api.Point2D
	positionsForBarracks api.Point2D
	barracksQuery        botutil.Query

	builder1 api.UnitTag
	builder2 api.UnitTag
	retreat  map[api.UnitTag]bool
}

func runAgent(info client.AgentInfo) {
	bot := bot{Bot: botutil.NewBot(info)}
	bot.LogActionErrors()

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
	bot.initLocations()
	for _, uc := range search.CalculateExpansionLocations(bot.Bot, false) {
		bot.baseLocations = append(bot.baseLocations, uc.Center())
	}
	bot.findBuildingsPositions()
	bot.retreat = map[api.UnitTag]bool{}

	// Send a friendly hello
	bot.Chat("(glhf)")
}

func (bot *bot) initLocations() {
	// My CC is on start position
	bot.myStartLocation = bot.Self[terran.CommandCenter].First().Pos2D()
	bot.enemyStartLocation = *bot.GameInfo().StartRaw.StartLocations[0]
}

func (bot *bot) findBuildingsPositions() {
	homeMinerals := bot.Neutral.Minerals().CloserThan(10, bot.myStartLocation)
	if homeMinerals.Len() == 0 {
		return // This should not happen
	}

	// Pick locations for supply depots
	pos := bot.myStartLocation.Offset(homeMinerals.Center(), -6)
	bot.positionsForSupplies = append(bot.positionsForSupplies, pos)
	bot.positionsForSupplies = append(bot.positionsForSupplies, neighbours8(pos, 2)...)

	// Determine proxy location
	pos = bot.enemyStartLocation.Offset(bot.myStartLocation, 25)
	pos = closestToPos(bot.baseLocations, pos).Offset(bot.myStartLocation, 1)
	bot.positionsForBarracks = pos

	// Build a re-usable query to check if we can build barracks
	bot.barracksQuery = botutil.NewQuery(bot)
	bot.barracksQuery.IgnoreResourceRequirements()

	bot.barracksQuery.Placement(ability.Build_Barracks, pos)
	for _, np := range neighbours8(pos, 4) {
		bot.barracksQuery.Placement(ability.Build_Barracks, np)
	}
}

func (bot *bot) getSCV() botutil.Unit {
	return bot.Self[terran.SCV].Choose(func(u botutil.Unit) bool { return u.IsGathering() }).First()
}

func (bot *bot) strategy() {
	// Update the home mineral (in case the old one mined out)
	bot.homeMineral = bot.Neutral.Minerals().CloserThan(10, bot.myStartLocation).First()

	// Build supply depots as needed
	depotCount := bot.Self.Count(terran.SupplyDepot) + bot.Self.Count(terran.SupplyDepotLowered)
	if bot.FoodLeft() < 6 && bot.Self.CountInProduction(terran.SupplyDepot) == 0 && depotCount < len(bot.positionsForSupplies) {
		pos := bot.positionsForSupplies[depotCount]
		if scv := bot.getSCV(); !scv.IsNil() {
			if !scv.BuildUnitAt(ability.Build_SupplyDepot, pos) {
				return
			}
		}
	}

	// Build barracks
	barracksCount := bot.Self.Count(terran.Barracks)
	if barracksCount < 4 {
		var scv botutil.Unit
		if barracksCount == 0 || barracksCount == 2 {
			// Get the builder for barracks 0 and 2
			scv = bot.UnitByTag(bot.builder1)
			if scv.IsNil() && bot.builder1 != 0 {
				scv = bot.getSCV()
				if !scv.IsNil() {
					bot.builder1 = scv.Tag
				}
			}
		} else {
			// Get the builder for barracks 1 and 3
			scv = bot.UnitByTag(bot.builder2)
			if scv.IsNil() && bot.builder2 != 0 {
				scv = bot.getSCV()
				if !scv.IsNil() {
					bot.builder2 = scv.Tag
				}
			}
		}
		if !scv.IsNil() {
			// Build the barracks
			if scv.Pos2D().Distance2(bot.positionsForBarracks) > 25 {
				// Move closer first to bust the fog
				scv.OrderPos(ability.Move, bot.positionsForBarracks)
			} else {
				// Query target build locations and use the first one that's available
				results := bot.barracksQuery.Execute()
				for i, result := range results.Placements() {
					if result.Result == api.ActionResult_Success {
						scv.BuildUnitAt(ability.Build_Barracks, *results.PlacementQuery(i).TargetPos)
						break
					}
				}
			}
		}
	}

	// Build a refinery for every two barracks
	refineryCount := bot.Self.Count(terran.Refinery)
	if refineryCount < (barracksCount+1)/2 {
		// Find first geyser that is close to my base, but it doesn't have Refinery on top of it
		if geyser := bot.Neutral.Vespene().CloserThan(10, bot.myStartLocation).Choose(func(u botutil.Unit) bool {
			return bot.Self[terran.Refinery].CloserThan(1, u.Pos2D()).First().IsNil()
		}).First(); !geyser.IsNil() {
			if scv := bot.getSCV(); !scv.IsNil() && !scv.BuildUnitOn(ability.Build_Refinery, geyser) {
				return
			}
		}
	}

	// Morph
	// bot.CanBuy(ability.Morph_OrbitalCommand) requires 550 minerals?
	if bot.Self.CountInProduction(terran.Reaper) >= 2 && bot.Minerals > 150 {
		if cc := bot.Self[terran.CommandCenter].Choose(func(u botutil.Unit) bool {
			return u.IsBuilt() && u.IsIdle()
		}).First(); !cc.IsNil() {
			cc.Order(ability.Morph_OrbitalCommand)
		}
	}
	if supply := bot.Self[terran.SupplyDepot].First(); !supply.IsNil() && supply.IsBuilt() {
		supply.Order(ability.Morph_SupplyDepot_Lower)
	}

	// Cast
	if cc := bot.Self[terran.OrbitalCommand].HasEnergy(50).First(); !cc.IsNil() {
		if !bot.homeMineral.IsNil() {
			cc.OrderTarget(ability.Effect_CalldownMULE, bot.homeMineral)
		}
	}

	// Units
	if bot.Self.CountAll(terran.SCV) < 18 {
		if !bot.BuildUnit(terran.OrbitalCommand, ability.Train_SCV) &&
			!bot.BuildUnit(terran.CommandCenter, ability.Train_SCV) &&
			!bot.BuildUnit(terran.PlanetaryFortress, ability.Train_SCV) {
			// do nothing
		}
	}
	bot.BuildUnits(terran.Barracks, ability.Train_Reaper, 10)
}

func (bot *bot) tactics() {
	// If there is idle scv, order it to gather minerals
	if !bot.homeMineral.IsNil() {
		idleSCVs := bot.Self[terran.SCV].Choose(func(u botutil.Unit) bool { return u.IsIdle() })
		bot.UnitsOrderTarget(idleSCVs, ability.Harvest_Gather, bot.homeMineral)
	}

	// Don't issue orders too often, or game won't be able to react
	if bot.GameLoop%6 == 0 {
		// If there is ready unsaturated refinery and an scv gathering, send it there
		if refinery := bot.Self[terran.Refinery].Choose(func(u botutil.Unit) bool {
			return u.IsBuilt() && u.AssignedHarvesters < 3
		}).First(); !refinery.IsNil() {
			if scv := bot.getSCV(); !scv.IsNil() {
				scv.OrderTarget(ability.Harvest_Gather, refinery)
			}
		}
	}

	if bot.GameLoop == 224 { // 10 sec
		if scv := bot.getSCV(); !scv.IsNil() {
			scv.OrderPos(ability.Move, bot.positionsForBarracks)
			bot.builder1 = scv.Tag
		}
	}
	if bot.GameLoop == 672 { // 30 sec
		if scv := bot.getSCV(); !scv.IsNil() {
			scv.OrderPos(ability.Move, bot.positionsForBarracks)
			bot.builder2 = scv.Tag
		}
	}

	// Attack!
	reapers := bot.Self[terran.Reaper]
	if reapers.Len() == 0 {
		return
	}

	targets := bot.getTargets()
	if targets.Len() == 0 {
		bot.UnitsOrderPos(reapers, ability.Attack, bot.enemyStartLocation)
		return
	}

	reapers.Each(func(reaper botutil.Unit) {
		// retreat
		if bot.retreat[reaper.Tag] && reaper.Health > 50 {
			delete(bot.retreat, reaper.Tag)
		}
		if reaper.Health < 21 || bot.retreat[reaper.Tag] {
			bot.retreat[reaper.Tag] = true
			reaper.OrderPos(ability.Move, bot.positionsForBarracks)
			return
		}

		target := targets.ClosestTo(reaper.Pos2D())

		// Keep range
		// Weapon is recharging
		if reaper.WeaponCooldown > 1 {
			// Enemy is closer than shooting distance - 0.5
			if reaper.InRange(target, -0.5) {
				// Retreat a little
				reaper.OrderPos(ability.Move, bot.positionsForBarracks)
				return
			}
		}

		// Attack
		if reaper.Pos2D().Distance2(target.Pos2D()) > 4*4 {
			// If target is far, attack it as unit, ling will run ignoring everything else
			reaper.OrderTarget(ability.Attack, target)
		} else if target.UnitType == zerg.ChangelingMarine || target.UnitType == zerg.ChangelingMarineShield {
			// Must specificially attack changelings, attack move is not enough
			reaper.OrderTarget(ability.Attack, target)
		} else {
			// Attack as position, ling will choose best target around
			reaper.OrderPos(ability.Attack, target.Pos2D())
		}
	})
}

func (bot *bot) getTargets() botutil.Units {
	// Prioritize things that can fight back
	if targets := bot.Enemy.Ground().CanAttack().All(); targets.Len() > 0 {
		return targets
	}

	// Otherwise just kill all the buildings
	return bot.Enemy.Ground().Structures().All()
}
