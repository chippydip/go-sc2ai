package agent

import (
	"fmt"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/unit"
	"github.com/chippydip/go-sc2ai/filter"
	"github.com/chippydip/go-sc2ai/search"
)

// GetAllUnits returns all units from the latest observation.
func (a *Agent) GetAllUnits() filter.Units {
	return a.info.Observation().GetObservation().GetRawData().GetUnits()
}

// GetUnit returns the first unit matching the given filter from the latest observation.
func (a *Agent) GetUnit(filter func(*api.Unit) bool) *api.Unit {
	for _, unit := range a.GetAllUnits() {
		if filter(unit) {
			return unit
		}
	}
	return nil
}

// GetClosestUnit returns the closest unit matching the given filter from the latest observation.
func (a *Agent) GetClosestUnit(pos api.Point2D, filter func(*api.Unit) bool) *api.Unit {
	return search.ClosestUnitWithFilter(pos, filter, a.GetAllUnits()...)
}

// GetUnits returns all units matching the given filter from the latest observation.
func (a *Agent) GetUnits(filter func(*api.Unit) bool) filter.Units {
	var list []*api.Unit
	for _, unit := range a.GetAllUnits() {
		if filter(unit) {
			list = append(list, unit)
		}
	}
	return list
}

// CountUnits ...
func (a *Agent) CountUnits(unitType api.UnitTypeID) int {
	n := 0
	for _, u := range a.GetAllUnits() {
		if u.Alliance == api.Alliance_Self && u.UnitType == unitType {
			n++
		}
	}
	return n
}

// CountUnitsInProduction ...
func (a *Agent) CountUnitsInProduction(unitType api.UnitTypeID, abilityType api.AbilityID) int {
	n := 0
	for _, u := range a.GetAllUnits() {
		if u.Alliance == api.Alliance_Self && u.UnitType == unitType {
			for _, order := range u.Orders {
				if order.AbilityId == abilityType {
					n++
				}
			}
		}
	}
	return n
}

// CountUnitsAndProduction ...
func (a *Agent) CountUnitsAndProduction(unitType api.UnitTypeID, abilityType api.AbilityID) int {
	return a.CountUnits(ability.Produces(abilityType)) + a.CountUnitsInProduction(unitType, abilityType)
}

// UnitHasAttribute ...
func (a *Agent) UnitHasAttribute(u *api.Unit, attr api.Attribute) bool {
	for _, a := range a.info.Data().Units[u.UnitType].Attributes {
		if a == attr {
			return true
		}
	}
	return false
}

// BuildUnit commands an available producer to use the train ability to build/morph/train/warp a unit.
// If the food, mineral, and vespene requirements are not met or no producer was found it does nothing and returns false.
func (a *Agent) BuildUnit(producer api.UnitTypeID, train api.AbilityID) bool {
	return a.BuildUnits(producer, train, 1) > 0
}

// BuildUnits commands available producers to use the train ability to build/morph/train/warp count units.
// Returns the number of units actually ordered based on producer, food, mineral, and vespene availability.
func (a *Agent) BuildUnits(producer api.UnitTypeID, train api.AbilityID, count int) int {
	cost := a.getProductionCost(producer, train)

	// Loop until done
	origCount := count
	for i := 0; count > 0; {
		// Check if we can afford one more
		if !a.canAfford(cost) {
			break
		}

		u := a.getNextProducer(&i, producer)
		if u == nil {
			break
		}

		// Produce the unit and adjust available resources
		a.UnitCommand(u.Tag, train)

		a.spend(cost)
		count--
	}

	return origCount - count
}

// TODO: BuildUnitsWithAddon

// BuildUnitAt ...
func (a *Agent) BuildUnitAt(producer api.UnitTypeID, train api.AbilityID, pos api.Point2D) bool {
	// Check if we can afford one
	cost := a.getProductionCost(producer, train)
	if !a.canAfford(cost) {
		return false
	}

	var i int
	u := a.getNextProducer(&i, producer)
	if u == nil {
		return false
	}

	// Produce the unit and adjust available resources
	a.UnitCommandAtPos(u.Tag, train, pos)
	a.spend(cost)
	return true
}

func (a *Agent) getProductionCost(producerType api.UnitTypeID, train api.AbilityID) unitCost {
	// Get the unit that will be built/trained
	targetType := ability.Produces(train)
	if targetType == unit.Invalid {
		panic(fmt.Sprintf("%v does not produce a unit", train))
	}

	producer := a.info.Data().Units[producerType]
	target := a.info.Data().Units[targetType]

	// Net food requirement:
	//  morphing units should have a net cost of 0
	//  interceptors should be net negative (clamp to zero)
	//  normal production will just be the target requirement
	//  zerglings are producded two at a time
	food := target.FoodRequired - producer.FoodRequired
	multiplier := 1
	if food < 0 {
		food = 0
	} else if food == 0.5 {
		multiplier = 2
	}

	// Double-check that we have an integer food cost now (do we need to handle anything other than zerglings?)
	foodMult := food * float32(multiplier)
	if float32(int(foodMult)) != foodMult {
		panic(fmt.Sprintf("unexpected FoodRequirement: %v -> %v x%v for %v", producer.FoodRequired, target.FoodRequired, multiplier, targetType))
	}

	// Return the per-build cost for this unit
	return unitCost{
		int(foodMult),
		int(target.MineralCost) * multiplier,
		int(target.VespeneCost) * multiplier,
	}
}

type unitCost struct {
	food, minerals, vespene int
}

func (a *Agent) canAfford(cost unitCost) bool {
	return (cost.food == 0 || a.FoodCap() >= a.FoodUsed()+cost.food) &&
		a.Minerals() >= cost.minerals && a.Vespene() >= cost.vespene
}

func (a *Agent) spend(cost unitCost) {
	a.playerCommon().FoodUsed += uint32(cost.food)
	a.playerCommon().Minerals -= uint32(cost.minerals)
	a.playerCommon().Vespene -= uint32(cost.vespene)
}

func (a *Agent) getNextProducer(i *int, producer api.UnitTypeID) *api.Unit {
	// Find the next available producer
	units := a.GetAllUnits()
	for ; *i < len(units); *i++ {
		u := units[*i]
		if u.Alliance != api.Alliance_Self {
			continue
		}
		// TODO: Take reactors into accountfor len(u.Orders)? u.AddOnTag -> Unit -> isReactor
		if u.UnitType == producer && u.BuildProgress == 1 && (len(u.Orders) == 0 || filter.IsWorker(u)) {
			*i++
			return u
		}
	}
	return nil
}
