package main

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/ability"
)

type Unit struct {
	api.Unit
	hps  float64
	Hits float64
}

type Cost struct {
	Minerals int
	Vespene  int
	Food     int
	Time     int
}

type Weapon struct {
	ground, air       *api.Weapon
	groundDps, airDps float64
}

var types []*api.UnitTypeData
var attributes map[api.UnitTypeID]map[api.Attribute]bool
var unitCost map[api.UnitTypeID]Cost   // todo: is it needed?
var abilityCost map[api.AbilityID]Cost // todo: add upgrades
var weapons map[api.UnitTypeID]Weapon

func InitUnits(typeData []*api.UnitTypeData) {
	types = typeData
	attributes = map[api.UnitTypeID]map[api.Attribute]bool{}
	unitCost = map[api.UnitTypeID]Cost{}
	abilityCost = map[api.AbilityID]Cost{}
	weapons = map[api.UnitTypeID]Weapon{}

	for _, td := range types {
		for _, attribute := range td.Attributes {
			attributes[td.UnitId] = map[api.Attribute]bool{}
			attributes[td.UnitId][attribute] = true
		}
		cost := Cost{
			Minerals: int(td.MineralCost),
			Vespene:  int(td.VespeneCost),
			Food:     int(td.FoodRequired - td.FoodProvided),
			Time:     int(td.BuildTime),
		}
		unitCost[td.UnitId] = cost
		abilityCost[td.AbilityId] = cost
		w := Weapon{}
		for _, weapon := range td.Weapons {
			if weapon.Type == api.Weapon_Ground || weapon.Type == api.Weapon_Any {
				w.ground = weapon
				w.groundDps = float64(weapon.Damage * float32(weapon.Attacks) / weapon.Speed)
			}
			if weapon.Type == api.Weapon_Air || weapon.Type == api.Weapon_Any {
				w.air = weapon
				w.airDps = float64(weapon.Damage * float32(weapon.Attacks) / weapon.Speed)
			}
			weapons[td.UnitId] = w
		}
	}
}

func (bot *proxyReapers) CanBuy(ability api.AbilityID) bool {
	cost := abilityCost[ability]
	return bot.minerals >= cost.Minerals && bot.vespene >= cost.Vespene && (cost.Food <= 0 || bot.foodLeft >= cost.Food)
}

func NewUnit(unit *api.Unit) *Unit {
	return &Unit{
		Unit: *unit,
		Hits: float64(unit.Health + unit.Energy)}
}

func (u Unit) Is(ids ...api.UnitTypeID) bool {
	for _, id := range ids {
		if u.UnitType == id {
			return true
		}
	}
	return false
}

func (u Unit) IsNot(ids ...api.UnitTypeID) bool {
	return !u.Is(ids...)
}

func (u Unit) IsIdle() bool {
	return len(u.Orders) == 0
}

func (u Unit) IsGathering() bool {
	return len(u.Orders) > 0 && u.Orders[0].AbilityId == ability.Harvest_Gather_SCV
}

func (u Unit) IsReady() bool {
	return u.BuildProgress == 1
}

func (u Unit) IsStructure() bool {
	return attributes[u.UnitType][api.Attribute_Structure]
}

func (u Unit) GroundDPS() float64 {
	return weapons[u.UnitType].groundDps
}

func (u Unit) AirDPS() float64 {
	return weapons[u.UnitType].airDps
}

func (u Unit) GroundRange() float64 {
	if weapon := weapons[u.UnitType].ground; weapon != nil {
		return float64(weapon.Range)
	}
	return -1
}

func (u Unit) AirRange() float64 {
	if weapon := weapons[u.UnitType].air; weapon != nil {
		return float64(weapon.Range)
	}
	return -1
}

func (u Unit) InRange(target *Unit, gap float64) bool {
	unitRange := -100.0
	if u.GroundDPS() > 0 && !target.IsFlying {
		unitRange = u.GroundRange()
	}
	// Air range is always larger than ground
	if u.AirDPS() > 0 && target.IsFlying {
		unitRange = u.AirRange()
	}
	dist := u.Pos.ToPoint2D().Distance(target.Pos.ToPoint2D())
	return dist-gap <= float64(u.Radius+target.Radius)+unitRange
}
