package botutil

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/unit"
)

type Builder struct {
	player  *Player
	units   *Units
	actions *Actions
}

func NewBuilder(info client.AgentInfo, player *Player, units *Units, actions *Actions) *Builder {
	b := &Builder{player, units, actions}
	// TODO: Only if player is Zerg
	if true {
		update := func() {
			n, data := 0, b.units.data
			for _, u := range b.units.units {
				// Count number of units that consume half a food
				if u.Alliance == api.Alliance_Self && data[u.UnitType].FoodRequired == 0.5 {
					n++
				}
			}
			// The game rounds fractional food down, but should really round up since you
			// this makes it seem like you can build units when you actually can't.
			if n%2 != 0 {
				b.player.FoodUsed++
			}
		}
		update()
		info.OnAfterStep(update)
	}
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

	// Loop until done
	origCount := count
	for i := 0; count > 0; {
		// Check if we can afford one more
		if !b.canAfford(cost) {
			break
		}

		u := b.getNextProducer(&i, producer)
		if u == nil {
			break
		}

		// Produce the unit and adjust available resources
		b.actions.UnitCommand(u, train)
		b.spend(cost)
		count--
	}

	return origCount - count
}

// TODO: BuildUnitsWithAddon

// BuildUnitAt ...
func (b *Builder) BuildUnitAt(producer api.UnitTypeID, train api.AbilityID, pos api.Point2D) bool {
	// Check if we can afford one
	cost := b.getProductionCost(producer, train)
	if !b.canAfford(cost) {
		return false
	}

	var i int
	u := b.getNextProducer(&i, producer)
	if u == nil {
		return false
	}

	// Produce the unit and adjust available resources
	b.actions.UnitCommandAtPos(u, train, &pos)
	b.spend(cost)
	return true
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

func (b *Builder) getNextProducer(i *int, producer api.UnitTypeID) *api.Unit {
	// Find the next available producer
	for ; *i < len(b.units.units); *i++ {
		u := b.units.units[*i]
		if u.Alliance != api.Alliance_Self {
			continue
		}
		// TODO: Take reactors into accountfor len(u.Orders)? u.AddOnTag -> Unit -> isReactor
		if u.UnitType == producer && u.BuildProgress == 1 && (len(u.Orders) == 0 || IsWorker(b.units.wrap(u))) {
			*i++
			return u
		}
	}
	return nil
}
