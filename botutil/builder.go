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
	used   map[api.UnitTag]bool
}

// NewBuilder creates a new Builder and registers it to fix FoodUsed rounding for zerg.
func NewBuilder(info client.AgentInfo, player *Player, units *UnitContext) *Builder {
	b := &Builder{player, units, map[api.UnitTag]bool{}}

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

		for k := range b.used {
			delete(b.used, k)
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

	cost := b.ProductionCost(producer, train)

	// Find all available producers
	origCount := count
	b.units.Self[producer].EachUntil(func(u Unit) bool {
		// Check if we can afford one more
		if !b.player.CanAfford(cost) {
			return false
		}

		if u.BuildProgress < 1 || (len(u.Orders) > 0 && u.IsStructure()) || b.used[u.Tag] {
			return false
		}

		// Produce the unit and adjust available resources
		u.Order(train)
		b.used[u.Tag] = true
		b.player.Spend(cost)
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
	cost := b.ProductionCost(producer, train)
	if !b.player.CanAfford(cost) {
		return false
	}

	// Find the closest available producer
	u := b.getNearestBuilder(producer, pos)
	if u.IsNil() {
		return false
	}

	// Produce the unit and adjust available resources
	u.OrderPos(train, pos)
	b.used[u.Tag] = true
	b.player.Spend(cost)
	return true
}

// BuildUnitOn commands an available producer to use the train ability to build/morph/train/warp a unit on the given target.
// If the food, mineral, and vespene requirements are not met or no producer was found it does nothing and returns false.
func (b *Builder) BuildUnitOn(producer api.UnitTypeID, train api.AbilityID, target Unit) bool {
	// Check if we can afford one
	cost := b.ProductionCost(producer, train)
	if !b.player.CanAfford(cost) {
		return false
	}

	// Find the closest available producer
	u := b.getNearestBuilder(producer, target.Pos2D())
	if u.IsNil() {
		return false
	}

	// Produce the unit and adjust available resources
	u.OrderTarget(train, target)
	b.used[u.Tag] = true
	b.player.Spend(cost)
	return true
}

func (b *Builder) getNearestBuilder(producer api.UnitTypeID, pos api.Point2D) Unit {
	builders := b.units.Self[producer].Drop(func(u Unit) bool {
		return u.BuildProgress < 1 || (len(u.Orders) > 0 && u.IsStructure()) || b.used[u.Tag]
	})
	return builders.ClosestTo(pos)
}

// ProductionCost computes the Cost for producerType to train once.
func (b *Builder) ProductionCost(producerType api.UnitTypeID, train api.AbilityID) Cost {
	// Special-case archon's since they are weird (and free)
	if train == ability.Morph_Archon {
		return Cost{}
	}

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

	// Per-build cost for this unit
	cost := Cost{
		Minerals: target.MineralCost * multiplier,
		Vespene:  target.VespeneCost * multiplier,
		Food:     uint32(foodMult),
	}

	// Except that morphs include the entire production cost, so subtract out the base unit's cost
	if producerType == unit.Zerg_Drone || isMorph[train] {
		cost.Minerals -= producer.MineralCost * multiplier
		cost.Vespene -= producer.VespeneCost * multiplier
	}
	return cost
}

// Build abilities which consume their producer
var isMorph = map[api.AbilityID]bool{
	ability.Morph_BroodLord:         true,
	ability.Morph_GreaterSpire:      true,
	ability.Morph_Hive:              true,
	ability.Morph_Lair:              true,
	ability.Morph_Lurker:            true,
	ability.Morph_OrbitalCommand:    true,
	ability.Morph_OverlordTransport: true,
	ability.Morph_Overseer:          true,
	ability.Morph_PlanetaryFortress: true,
	ability.Morph_Ravager:           true,
	ability.Train_Baneling:          true,
	// ability.Morph_Archon:            true,
	// ability.Morph_Hellbat:           true,
	// ability.Morph_Hellion:           true,
	// ability.Morph_Mothership:        true,
	// ability.Morph_Gateway:           true,
	// ability.Morph_WarpGate:          true,
}

// Convenience methods for giving orders directly to units:

// BuildUnitAt ...
func (u Unit) BuildUnitAt(train api.AbilityID, pos api.Point2D) bool {
	if u.IsNil() {
		return false
	}
	b := u.ctx.bot

	// Check if we can afford one
	cost := b.ProductionCost(u.UnitType, train)
	if !b.player.CanAfford(cost) {
		return false
	}

	// Produce the unit and adjust available resources
	u.OrderPos(train, pos)
	b.used[u.Tag] = true
	b.player.Spend(cost)
	return true

}

// BuildUnitOn ...
func (u Unit) BuildUnitOn(train api.AbilityID, target Unit) bool {
	if u.IsNil() {
		return false
	}
	b := u.ctx.bot

	// Check if we can afford one
	cost := b.ProductionCost(u.UnitType, train)
	if !b.player.CanAfford(cost) {
		return false
	}

	// Produce the unit and adjust available resources
	u.OrderTarget(train, target)
	b.used[u.Tag] = true
	b.player.Spend(cost)
	return true

}
