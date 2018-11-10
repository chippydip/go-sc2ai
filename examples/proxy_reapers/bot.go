package main

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/terran"
	"github.com/chippydip/go-sc2ai/enums/zerg"
)

func (bot *proxyReapers) initLocations() {
	// My CC is on start position
	bot.myStartLocation = bot.units[terran.CommandCenter].First().Pos.ToPoint2D()
	bot.enemyStartLocation = *bot.info.GameInfo().StartRaw.StartLocations[0]
}

func (bot *proxyReapers) findBuildingsPositions() {
	homeMinerals := bot.mineralFields.Units().CloserThan(10, bot.myStartLocation)
	if homeMinerals.Len() == 0 {
		return // This should not happen
	}
	vec := SquareDirection(homeMinerals.Center().VecTo(bot.myStartLocation))
	pos := bot.myStartLocation.Add(vec.Mul(3.5))
	bot.positionsForSupplies = append(bot.positionsForSupplies, pos)
	bot.positionsForSupplies = append(bot.positionsForSupplies, Neighbours4(pos, 2)...)

	pos = bot.enemyStartLocation.Offset(bot.myStartLocation, 25)
	pos = ClosestToP2D(bot.baseLocations, pos).Offset(bot.myStartLocation, 1)

	pfb := []*api.RequestQueryBuildingPlacement{{
		AbilityId: ability.Build_Barracks,
		TargetPos: &pos}}
	for _, np := range Neighbours8(pos, 4) {
		if bot.isBuildable(np) {
			npc := np // Because you can't just pass address, you need copy of value
			pfb = append(pfb, &api.RequestQueryBuildingPlacement{
				AbilityId: ability.Build_Barracks,
				TargetPos: &npc})
		}
	}
	resp := bot.info.Query(api.RequestQuery{
		Placements:                 pfb,
		IgnoreResourceRequirements: true})

	for key, result := range resp.Placements {
		if result.Result == api.ActionResult_Success {
			bot.positionsForBarracks = append(bot.positionsForBarracks, *pfb[key].TargetPos)
		}
	}
}

func (bot *proxyReapers) getSCV() *Unit {
	return bot.units[terran.SCV].FirstFiltered(func(unit *Unit) bool { return unit.IsGathering() })
}

func (bot *proxyReapers) strategy() {
	// Buildings
	suppliesCount := bot.units.OfType(terran.SupplyDepot, terran.SupplyDepotLowered).Len()
	if bot.CanBuy(ability.Build_SupplyDepot) && bot.orders[ability.Build_SupplyDepot] == 0 &&
		suppliesCount < len(bot.positionsForSupplies) && bot.foodLeft < 6 {
		pos := bot.positionsForSupplies[suppliesCount]
		if scv := bot.getSCV(); scv != nil {
			bot.unitCommandTargetPos(scv, ability.Build_SupplyDepot, pos, false)
		}
	}
	raxPending := bot.units[terran.Barracks].Len()
	if bot.CanBuy(ability.Build_Barracks) && raxPending < 3 && len(bot.positionsForBarracks) > raxPending {
		pos := bot.positionsForBarracks[raxPending]
		scv := bot.units[terran.SCV].ByTag(bot.builder2)
		if raxPending == 0 || raxPending == 2 {
			scv = bot.units[terran.SCV].ByTag(bot.builder1)
		}
		if scv == nil {
			scv = bot.getSCV()
		}
		if scv != nil {
			bot.unitCommandTargetPos(scv, ability.Build_Barracks, pos, false)
		}
		return
	}
	if bot.CanBuy(ability.Build_Refinery) && (raxPending == 1 && bot.units[terran.Refinery].Len() == 0 ||
		raxPending == 3 && bot.units[terran.Refinery].Len() == 1) {
		// Find first geyser that is close to my base, but it doesn't have Refinery on top of it
		geyser := bot.vespeneGeysers.Units().CloserThan(10, bot.myStartLocation).FirstFiltered(func(unit *Unit) bool {
			return bot.units[terran.Refinery].CloserThan(1, unit.Pos.ToPoint2D()).Len() == 0
		})
		if scv := bot.getSCV(); scv != nil && geyser != nil {
			bot.unitCommandTargetTag(scv, ability.Build_Refinery, geyser.Tag, false)
		}
	}

	// Morph
	cc := bot.units[terran.CommandCenter].
		FirstFiltered(func(unit *Unit) bool { return unit.IsReady() && unit.IsIdle() })
	// bot.CanBuy(ability.Morph_OrbitalCommand) requires 550 minerals?
	if cc != nil && bot.orders[ability.Train_Reaper] >= 2 && bot.minerals >= 150 {
		bot.unitCommand(cc, ability.Morph_OrbitalCommand)
	}
	if supply := bot.units[terran.SupplyDepot].First(); supply != nil {
		bot.unitCommand(supply, ability.Morph_SupplyDepot_Lower)
	}

	// Cast
	cc = bot.units[terran.OrbitalCommand].FirstFiltered(func(unit *Unit) bool { return unit.Energy >= 50 })
	if cc != nil {
		if homeMineral := bot.mineralFields.Units().CloserThan(10, bot.myStartLocation).First(); homeMineral != nil {
			bot.unitCommandTargetTag(cc, ability.Effect_CalldownMULE, homeMineral.Tag, false)
		}
	}

	// Units
	ccs := bot.units.OfType(terran.CommandCenter, terran.OrbitalCommand, terran.PlanetaryFortress)
	cc = ccs.FirstFiltered(func(unit *Unit) bool { return unit.IsReady() && unit.IsIdle() })
	if cc != nil && bot.units[terran.SCV].Len() < 20 && bot.CanBuy(ability.Train_SCV) {
		bot.unitCommand(cc, ability.Train_SCV)
		return
	}
	if rax := bot.units[terran.Barracks].
		FirstFiltered(func(unit *Unit) bool { return unit.IsReady() && unit.IsIdle() }); rax != nil && bot.CanBuy(ability.Train_Reaper) {
		bot.unitCommand(rax, ability.Train_Reaper)
	}
}

func (bot *proxyReapers) tactics() {
	step := bot.info.Observation().Observation.GameLoop
	// If there is idle scv, order it to gather minerals
	if scv := bot.units[terran.SCV].FirstFiltered(func(unit *Unit) bool { return unit.IsIdle() }); scv != nil {
		if homeMineral := bot.mineralFields.Units().CloserThan(10, bot.myStartLocation).First(); homeMineral != nil {
			bot.unitCommandTargetTag(scv, ability.Harvest_Gather_SCV, homeMineral.Tag, false)
		}
	}
	// Don't issue orders too often, or game won't be able to react
	if step%6 == 0 {
		// If there is ready unsaturated refinery and an scv gathering, send it there
		refinery := bot.units[terran.Refinery].
			FirstFiltered(func(unit *Unit) bool { return unit.IsReady() && unit.AssignedHarvesters < 3 })
		if refinery != nil {
			if scv := bot.getSCV(); scv != nil {
				bot.unitCommandTargetTag(scv, ability.Harvest_Gather_SCV, refinery.Tag, false)
			}
		}
	}

	if step == 224 { // 10 sec
		scv := bot.getSCV()
		pos := bot.positionsForBarracks[0]
		bot.unitCommandTargetPos(scv, ability.Move, pos, false)
		bot.builder1 = scv.Tag
	}
	if step == 672 { // 30 sec
		scv := bot.getSCV()
		pos := bot.positionsForBarracks[1]
		bot.unitCommandTargetPos(scv, ability.Move, pos, false)
		bot.builder2 = scv.Tag
	}

	bot.okTargets = nil
	bot.goodTargets = nil
	for _, units := range bot.enemyUnits {
		for _, unit := range units {
			if !unit.IsFlying && unit.IsNot(zerg.Larva, zerg.Egg) {
				bot.okTargets.Add(unit)
				if !unit.IsStructure() {
					bot.goodTargets.Add(unit)
				}
			}
		}
	}

	reapers := bot.units[terran.Reaper]
	if len(bot.okTargets) == 0 {
		bot.unitsCommandTargetPos(reapers, ability.Attack, bot.enemyStartLocation)
	} else {
		for _, reaper := range reapers {
			// retreat
			if bot.retreat[reaper.Tag] && reaper.Health > 50 {
				delete(bot.retreat, reaper.Tag)
			}
			if reaper.Health < 21 || bot.retreat[reaper.Tag] {
				bot.retreat[reaper.Tag] = true
				bot.unitCommandTargetPos(reaper, ability.Move, bot.positionsForBarracks[0], false)
				continue
			}

			// Keep range
			// Weapon is recharging
			if reaper.WeaponCooldown > 1 {
				// There is an enemy
				if closestEnemy := bot.goodTargets.ClosestTo(reaper.Pos.ToPoint2D()); closestEnemy != nil {
					// And it is closer than shooting distance - 0.5
					if reaper.InRange(closestEnemy, -0.5) {
						// Retreat a little
						bot.unitCommandTargetPos(reaper, ability.Move, bot.positionsForBarracks[0], false)
						continue
					}
				}
			}

			// Attack
			if len(bot.goodTargets) > 0 {
				target := bot.goodTargets.ClosestTo(reaper.Pos.ToPoint2D())
				// Snapshots couldn't be targeted using tags
				if reaper.Pos.ToPoint2D().Distance(target.Pos.ToPoint2D()) > 4 &&
					target.DisplayType != api.DisplayType_Snapshot {
					// If target is far, attack it as unit, ling will run ignoring everything else
					bot.unitCommandTargetTag(reaper, ability.Attack, target.Tag, false)
				} else {
					// Attack as position, ling will choose best target around
					bot.unitCommandTargetPos(reaper, ability.Attack, target.Pos.ToPoint2D(), false)
				}
			} else {
				target := bot.okTargets.ClosestTo(reaper.Pos.ToPoint2D())
				bot.unitCommandTargetPos(reaper, ability.Attack, target.Pos.ToPoint2D(), false)
			}
		}
	}
}
