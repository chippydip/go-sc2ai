package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

// Enemy()
//   EnemyFlying()
//     EnmyFlyingPassive()
//       EnemyFlyingPassiveUnits()
//       EnemyFlyingPassiveStructures() \
//     EnemyFlyingAttackers()           | EnemyFlyingStructures()
//       EnemyFlyingAttackStructures()  /              \
//       EnemyFlyingAttackUnits() \                    |
//   EnemyGround()                | EnemyAttackUnits() | EnemyAttackers()
//     EnemyGroundAttackers()     |                    |
//       EnemyGroundAttackUnits() /                    |
//       EnemyGroundAttackStructures()  \              /
//     EnemyGroundPassive()             | EnemyGroundStructures()
//       EnemyGroundPassiveStructures() /
//       EnemyGroundPassiveUnits()
// EnemyStructures() -> EnemyGroundStructures() + EnemyFlyingStructures() ?

const (
	neutralAll       = 24
	neutralResources = 25
	neutralMinerals  = 25
	neutralVespene   = 26
)

// UnitContext stores shared state about units from an observation and provides filtered access to those units.
type UnitContext struct {
	raw     []*api.Unit
	data    []*api.UnitTypeData
	wrapped []Unit
	byTag   map[api.UnitTag]*Unit

	groups [28]int

	Self    self
	Ally    ally
	Enemy   enemy
	Neutral neutral

	bot *Bot

	dummy Units
}

// NewUnitContext creates a new context and registers it to update after each step.
func NewUnitContext(info client.AgentInfo, bot *Bot) *UnitContext {
	ctx := &UnitContext{
		byTag:   map[api.UnitTag]*Unit{},
		Self:    self{},
		Ally:    ally{},
		Enemy:   enemy{},
		Neutral: neutral{},
		bot:     bot,
	}
	ctx.dummy = Units{raw: []Unit{Unit{ctx: ctx}}}
	update := func() {
		// Load the latest observation
		ctx.raw = info.Observation().GetObservation().GetRawData().GetUnits()
		ctx.data = info.Data().GetUnits()

		ctx.update()
	}
	update()
	info.OnAfterStep(update)
	return ctx
}

func (ctx *UnitContext) update() {
	for k := range ctx.byTag {
		delete(ctx.byTag, k)
	}

	// Reset maps without clearing to eliminate most allocs
	ctx.clear(ctx.Self)
	ctx.clear(ctx.Ally)
	ctx.clear(ctx.Enemy)
	ctx.clear(ctx.Neutral)

	// This should never happen, but just in case...
	if len(ctx.raw) == 0 {
		return
	}

	// Pre-sorting appears to actually improve overall speed, presumably it improves
	// locality of reference enough when setting tags (all units of same type in a row)
	// and potentially reduces the time for the real sort due to grouping here that it
	// more than covers the cost of this extra sort call.
	sortUnits(&ctx.raw)

	// Sort the units in place so common queries can use direct slices and avoid allocation
	for _, u := range ctx.raw {
		setSortTag(u, ctx.data[u.UnitType])
	}
	sortUnits(&ctx.raw)

	// Allocate a new array for wrapped unit objects
	ctx.wrapped = make([]Unit, len(ctx.raw))

	// Slice up the sorted result
	(&grouper{}).group(ctx)
}

func (ctx *UnitContext) clear(m map[api.UnitTypeID]Units) {
	// Reset maps without clearing to eliminate most allocs
	for k := range m {
		m[k] = Units{}
	}

	// Stash a ctx pointer in an easy to find place
	m[0] = ctx.dummy
}

// Use the high bits of a UnitTypeID to allow us to specify sort criteria and stuff sort in-place.
// Top 5 bits are the alliance, next 3 are used for grouping, the remaining 24 hold the UnitTypeID.
const (
	idFlying    api.UnitTypeID = 1 << 26
	idWeapons   api.UnitTypeID = 1 << 25
	idStructure api.UnitTypeID = 1 << 24

	idVespene api.UnitTypeID = 2 << 24
	idMineral api.UnitTypeID = 1 << 24

	idMask      api.UnitTypeID = (1 << 24) - 1
	idGroupMask api.UnitTypeID = ^idMask
)

func setSortTag(u *api.Unit, d *api.UnitTypeData) {
	id := u.UnitType
	id |= api.UnitTypeID(allianceIndex(u.Alliance) << 27) // top 5 bits

	if u.Alliance == api.Alliance_Neutral {
		switch {
		case d.HasVespene:
			id |= idVespene
		case d.HasMinerals:
			id |= idMineral
		}
	} else {
		grp := api.UnitTypeID(0)

		for _, a := range d.Attributes {
			if a == api.Attribute_Structure {
				grp |= idStructure
				break
			}
		}

		if len(d.Weapons) > 0 {
			grp = idWeapons - 1 - grp // invert unit/structure order for attackers
			grp |= idWeapons
		}

		if u.IsFlying {
			grp = idFlying - 1 - grp // invert order for flying units
			grp |= idFlying
		}

		id |= grp & idGroupMask
	}

	u.UnitType = id
}

// Stack tracking while finding the boundaries of the sorted groups.

type grouper struct {
	lastGroup int
	typeStart int
	prevType  api.UnitTypeID
}

func (g *grouper) group(ctx *UnitContext) {
	g.prevType = ctx.raw[0].UnitType
	for i, u := range ctx.raw {
		if u.UnitType != g.prevType {
			g.updateMap(ctx, i)

			grp := int(u.UnitType >> 24)
			if grp != g.lastGroup {
				g.updateGroups(ctx, i, grp)
			}

			g.prevType = u.UnitType
		}

		// Revert the unit type so it can be used for data lookup again
		u.UnitType &= idMask

		// Wrap the unit
		ctx.wrapped[i] = Unit{ctx, u, ctx.data[u.UnitType]}
		ctx.byTag[u.Tag] = &ctx.wrapped[i]
	}
	g.updateMap(ctx, len(ctx.raw))
	g.updateGroups(ctx, len(ctx.raw), len(ctx.groups)-1)
}

func (g *grouper) updateMap(ctx *UnitContext, i int) {
	m := ctx.alliance(g.prevType >> 27)
	t := g.prevType & idMask
	r := m[t].raw
	n := ctx.wrapped[g.typeStart:i]
	if r == nil {
		// Normal case
		m[t] = Units{raw: n}
	} else {
		// Only happens if there are flying and ground units of the same type (locust?)
		s := make([]Unit, len(r)+len(n))
		copy(s, r)
		copy(s[:len(r)], n)
		m[t] = Units{raw: s}
	}
	g.typeStart = i
}

func (g *grouper) updateGroups(ctx *UnitContext, i int, grp int) {
	for ii := g.lastGroup + 1; ii <= grp; ii++ {
		ctx.groups[ii] = i // start index of group
	}
	g.lastGroup = grp
}

func allianceIndex(alliance api.Alliance) int {
	switch alliance {
	case api.Alliance_Self:
		return 0
	case api.Alliance_Ally:
		return 1
	case api.Alliance_Enemy:
		return 2
	case api.Alliance_Neutral:
		return 3
	default:
		return -1
	}
}

func (ctx *UnitContext) alliance(index api.UnitTypeID) map[api.UnitTypeID]Units {
	switch index {
	case 0:
		return ctx.Self
	case 1:
		return ctx.Ally
	case 2:
		return ctx.Enemy
	case 3:
		return ctx.Neutral
	default:
		return nil
	}
}

// WasObserved returns true if a unit with the given tag was inlucded in the last observation.
func (ctx *UnitContext) WasObserved(tag api.UnitTag) bool {
	return ctx.byTag[tag] != nil
}

// UnitByTag returns a the unit with the given tag if it was in included in the last observation.
func (ctx *UnitContext) UnitByTag(tag api.UnitTag) Unit {
	if ptr := ctx.byTag[tag]; ptr != nil {
		return *ptr
	}
	return Unit{}
}

// AllUnits returns all units from the most recent observation.
func (ctx *UnitContext) AllUnits() Units {
	return Units{raw: ctx.wrapped}
}
