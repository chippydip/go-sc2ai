package botutil

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/unit"
)

// Builder provides operations to make building/morphing/training/warping units easier.
type Builder struct {
	player *Player
	units  *UnitContext
}

// NewBuilder creates a new Builder and registers it to fix FoodUsed rounding for zerg.
func NewBuilder(info client.AgentInfo, player *Player, units *UnitContext) *Builder {
	b := &Builder{player, units}

	// This is only really an issue for zerg
	if player.RaceActual != api.Race_Zerg {
		return b
	}

	update := func() {
		// Count number of units that consume half a food
		n := b.units.Self.CountIf(func(u Unit) bool {
			return u.FoodRequired == 0.5
		})

		// The game rounds fractional food down, but should really round up since
		// this makes it seem like you can build units when you actually can't.
		if n%2 != 0 {
			b.player.FoodUsed++
		}
	}
	update()
	info.OnAfterStep(update)

	return b
}

// BuildUnit commands an available producer to use the train ability to build/morph/train/warp a unit.
// If the food, mineral, and vespene requirements are not met or no producer was found it does nothing and returns false.
func (b *Builder) BuildUnit(producer api.UnitTypeID, train api.AbilityID) bool {
	return b.BuildUnits(producer, train, 1) > 0
}

// BuildUnits commands available producers to use the train ability to build/morph/train/warp count units.
// Returns the number of units actually ordered based on producer, food, mineral, and vespene availability.
func (b *Builder) BuildUnits(producer api.UnitTypeID, train api.AbilityID, count int) int {
	if count <= 0 {
		return 0
	}

	cost := b.getProductionCost(producer, train)

	// Find all available producers
	origCount := count
	b.units.Self[producer].EachUntil(func(u Unit) bool {
		// Check if we can afford one more
		if !b.canAfford(cost) {
			return false
		}

		if u.BuildProgress < 1 || (len(u.Orders) > 0 && u.IsStructure()) {
			return false
		}

		// Produce the unit and adjust available resources
		u.Order(train)
		b.spend(cost)
		count--
		return count == 0
	})

	return origCount - count
}

// TODO: BuildUnitsWithAddon

// BuildUnitAt commands an available producer to use the train ability to build/morph/train/warp a unit at the given location.
// If the food, mineral, and vespene requirements are not met or no producer was found it does nothing and returns false.
func (b *Builder) BuildUnitAt(producer api.UnitTypeID, train api.AbilityID, pos api.Point2D) bool {
	// Check if we can afford one
	cost := b.getProductionCost(producer, train)
	if !b.canAfford(cost) {
		return false
	}

	// Find the closest available producer
	u := b.getNearestBuilder(producer, pos)
	if u.IsNil() {
		return false
	}

	// Produce the unit and adjust available resources
	u.OrderPos(train, &pos)
	b.spend(cost)
	return true
}

// BuildUnitOn commands an available producer to use the train ability to build/morph/train/warp a unit on the given target.
// If the food, mineral, and vespene requirements are not met or no producer was found it does nothing and returns false.
func (b *Builder) BuildUnitOn(producer api.UnitTypeID, train api.AbilityID, target Unit) bool {
	// Check if we can afford one
	cost := b.getProductionCost(producer, train)
	if !b.canAfford(cost) {
		return false
	}

	// Find the closest available producer
	u := b.getNearestBuilder(producer, target.Pos2D())
	if u.IsNil() {
		return false
	}

	// Produce the unit and adjust available resources
	u.OrderTarget(train, target)
	b.spend(cost)
	return true
}

func (b *Builder) getNearestBuilder(producer api.UnitTypeID, pos api.Point2D) Unit {
	builders := b.units.Self[producer].Drop(func(u Unit) bool {
		return u.BuildProgress < 1 || (len(u.Orders) > 0 && u.IsStructure())
	})
	return builders.ClosestTo(pos)
}

func (b *Builder) getProductionCost(producerType api.UnitTypeID, train api.AbilityID) unitCost {
	// Get the unit that will be built/trained
	targetType := ability.Produces(train)
	if targetType == unit.Invalid {
		log.Panicf("%v does not produce a unit", train)
	}

	producer := b.units.data[producerType]
	target := b.units.data[targetType]

	// Net food requirement:
	//  morphing units should have a net cost of 0
	//  interceptors should be net negative (clamp to zero)
	//  normal production will just be the target requirement
	//  zerglings are producded two at a time
	food, multiplier := target.FoodRequired-producer.FoodRequired, uint32(1)
	if food < 0 {
		food = 0
	} else if food == 0.5 {
		multiplier = 2
	}

	// Double-check that we have an integer food cost now (do we need to handle anything other than zerglings?)
	foodMult := food * float32(multiplier)
	if float32(uint32(foodMult)) != foodMult {
		log.Panicf("unexpected FoodRequirement: %v -> %v x%v for %v", producer.FoodRequired, target.FoodRequired, multiplier, targetType)
	}

	// Return the per-build cost for this unit
	return unitCost{
		uint32(foodMult),
		target.MineralCost * multiplier,
		target.VespeneCost * multiplier,
	}
}

type unitCost struct {
	food, minerals, vespene uint32
}

func (b *Builder) canAfford(cost unitCost) bool {
	return (cost.food == 0 || b.player.FoodCap >= b.player.FoodUsed+cost.food) &&
		b.player.Minerals >= cost.minerals && b.player.Vespene >= cost.vespene
}

func (b *Builder) spend(cost unitCost) {
	b.player.FoodUsed += cost.food
	b.player.Minerals -= cost.minerals
	b.player.Vespene -= cost.vespene
}

// Convenience methods for giving orders directly to units:

// BuildUnitAt ...
func (u Unit) BuildUnitAt(train api.AbilityID, pos api.Point2D) bool {
	if u.IsNil() {
		return false
	}
	b := u.ctx.bot

	// Check if we can afford one
	cost := b.getProductionCost(u.UnitType, train)
	if !b.canAfford(cost) {
		return false
	}

	// Produce the unit and adjust available resources
	u.OrderPos(train, &pos)
	b.spend(cost)
	return true

}

// BuildUnitOn ...
func (u Unit) BuildUnitOn(train api.AbilityID, target Unit) bool {
	if u.IsNil() {
		return false
	}
	b := u.ctx.bot

	// Check if we can afford one
	cost := b.getProductionCost(u.UnitType, train)
	if !b.canAfford(cost) {
		return false
	}

	// Produce the unit and adjust available resources
	u.OrderTarget(train, target)
	b.spend(cost)
	return true

}
