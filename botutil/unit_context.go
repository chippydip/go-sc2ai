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
	raw   []*api.Unit
	data  []*api.UnitTypeData
	byTag map[api.UnitTag]*api.Unit
	loop  uint32

	groups [28]int

	Self    self
	Ally    ally
	Enemy   enemy
	Neutral neutral

	actions *Actions
}

// NewUnitContext creates a new context and registers it to update after each step.
func NewUnitContext(info client.AgentInfo, actions *Actions) *UnitContext {
	u := &UnitContext{
		byTag:   map[api.UnitTag]*api.Unit{},
		Self:    self{},
		Ally:    ally{},
		Enemy:   enemy{},
		Neutral: neutral{},
		actions: actions,
	}
	update := func() {
		// Load the latest observation
		u.raw = info.Observation().GetObservation().GetRawData().GetUnits()
		u.data = info.Data().GetUnits()
		u.loop = info.Observation().GetObservation().GetGameLoop()
		u.update()
	}
	update()
	info.OnAfterStep(update)
	return u
}

func (ctx *UnitContext) wrap(u *api.Unit) Unit {
	if u == nil {
		return Unit{}
	}
	return Unit{ctx, u, ctx.data[u.UnitType]}
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
	for _, uu := range ctx.raw {
		ctx.byTag[uu.Tag] = uu
		setSortTag(uu, ctx.data[uu.UnitType])
	}
	sortUnits(&ctx.raw)

	// Slice up the sorted result
	(&grouper{}).group(ctx)
}

func (ctx *UnitContext) clear(m map[api.UnitTypeID]Units) {
	// Reset maps without clearing to eliminate most allocs
	for k := range m {
		m[k] = Units{}
	}

	// Stash a ctx pointer in an easy to find place
	m[0] = Units{ctx: ctx}
}

// Use the high bits of a UnitTypeID to allow us to specify sort criteria and stuff sort in-place.
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
			grp = idWeapons - grp // invert unit/structure order for attackers
			grp |= idWeapons
		}

		if u.IsFlying {
			grp = idFlying - grp // invert order for flying units
			grp |= idFlying
		}

		id |= grp
	}

	u.UnitType = id
}

// Stack tracking while finding the boundaries of the sorted groups.

type grouper struct {
	lastGroup int
	typeStart int
	prevType  api.UnitTypeID
}

func (g *grouper) group(u *UnitContext) {
	g.prevType = u.raw[0].UnitType
	for i, uu := range u.raw {
		if uu.UnitType != g.prevType {
			g.updateMap(u, i)

			grp := int(uu.UnitType >> 24)
			if grp != g.lastGroup {
				g.updateGroups(u, i, grp)
			}

			g.prevType = uu.UnitType
		}

		// Revert the unit type so it can be used for data lookup again
		uu.UnitType &= idMask
	}
	g.updateMap(u, len(u.raw))
	g.updateGroups(u, len(u.raw), len(u.groups)-1)
}

func (g *grouper) updateMap(u *UnitContext, i int) {
	m := u.alliance(g.prevType >> 27)
	t := g.prevType & idMask
	r := m[t].raw
	n := u.raw[g.typeStart:i]
	if r == nil {
		// Normal case
		m[t] = Units{ctx: u, raw: n}
	} else {
		// Only happens if there are flying and ground units of the same type (locust?)
		s := make([]*api.Unit, len(r)+len(n))
		copy(s, r)
		copy(s[:len(r)], n)
		m[t] = Units{ctx: u, raw: s}
	}
	g.typeStart = i
}

func (g *grouper) updateGroups(u *UnitContext, i int, grp int) {
	for ii := g.lastGroup + 1; ii <= grp; ii++ {
		u.groups[ii] = i // start index of group
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
