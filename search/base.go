package search

import (
	"log"
	"strings"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
)

// Base ...
type Base struct {
	m *Map
	i int

	ResourceCenter api.Point2D
	MineralCenter  api.Point2D
	Minerals       []botutil.Unit
	Geysers        []botutil.Unit

	Location     api.Point2D
	TownHall     botutil.Unit
	GasBuildings map[api.Point2D]botutil.Unit
	SelfWorkers  map[api.UnitTag]bool
	OtherWorkers map[api.UnitTag]bool
}

func newBase(m *Map, i int, loc BaseLocation) *Base {
	// Re-compute center with 4x weight on geysers to better represent unbalanced gas bases
	cluster := UnitCluster{}
	minerals := UnitCluster{}
	for _, u := range loc.Resources.Units() {
		if u.HasVespene {
			cluster.Add(u)
			cluster.Add(u)
			cluster.Add(u)
		} else {
			minerals.Add(u)
		}
		cluster.Add(u)
	}

	return &Base{
		m:              m,
		i:              i,
		ResourceCenter: cluster.Center(),
		MineralCenter:  minerals.Center(),
		Minerals:       make([]botutil.Unit, 0, loc.Resources.Count()),
		Geysers:        make([]botutil.Unit, 0, 2),
		Location:       loc.Location,
		GasBuildings:   map[api.Point2D]botutil.Unit{},
		SelfWorkers:    map[api.UnitTag]bool{},
		OtherWorkers:   map[api.UnitTag]bool{},
	}
}

func (base *Base) updateResource(u botutil.Unit) {
	switch {
	case u.HasMinerals:
		base.Minerals = base.updateOrAdd(base.Minerals, u)
	case u.HasVespene:
		base.Geysers = base.updateOrAdd(base.Geysers, u)
	default:
		log.Panicf("unknown resource: %v", u)
	}
}

func (base *Base) update(bot *botutil.Bot) {
	// Check for exhausted minerals
	for i := 0; i < len(base.Minerals); i++ {
		u := base.Minerals[i]
		if !bot.WasObserved(u.Tag) {
			copy(base.Minerals[i:], base.Minerals[i+1:])
			base.Minerals = base.Minerals[:len(base.Minerals)-1]
			i--
		}
	}

	// Clear fields that are re-computed each loop
	base.TownHall = botutil.Unit{}
	for k := range base.GasBuildings {
		delete(base.GasBuildings, k)
	}
	for k := range base.SelfWorkers {
		delete(base.SelfWorkers, k)
	}
	for k := range base.OtherWorkers {
		delete(base.OtherWorkers, k)
	}
}

func (base *Base) updateOrAdd(units []botutil.Unit, u botutil.Unit) []botutil.Unit {
	for i, u2 := range units {
		if u2.Pos2D().Distance2(u.Pos2D()) < 1 {
			if u2.Pos2D() != u.Pos2D() {
				log.Panicf("%v != %v", u2.Pos2D(), u.Pos2D())
			}

			units[i] = u
			if u.IsSnapshot() {
				// TODO: Move this to botutil?
				// Not populated for snapshots
				// u.Health = u2.Health
				// u.HealthMax = u2.HealthMax
				// u.Shield = u2.Shield
				// u.ShieldMax = u2.ShieldMax
				// u.Energy = u2.Energy
				// u.EnergyMax = u2.EnergyMax
				u.MineralContents = u2.MineralContents
				u.VespeneContents = u2.VespeneContents
				// u.IsFlying = u2.IsFlying
				// u.IsBurrowed = u2.IsBurrowed

				// update Tag?
				//u.Tag = u2.Tag

				// facing
				// detect_range, radar_range, is_powered?
			}
			return units
		}
	}

	// Not found, append
	units = append(units, u)

	// Keep sorted
	uIsSmall, uDist := strings.HasSuffix(u.Name, "750"), u.Pos2D().Distance2(base.Location)
	for i := 0; i < len(units)-1; i++ {
		isSmall := strings.HasSuffix(units[i].Name, "750")
		if !isSmall && uIsSmall {
			continue // small patches after big ones, regardless of distance
		}

		dist := units[i].Pos2D().Distance2(base.Location)
		if uIsSmall != isSmall || uDist < dist {
			// found insertion point, shift back the rest and insert again
			copy(units[i+1:], units[i:])
			units[i] = u
			return units
		}
	}

	return units
}

// IsSelfOwned returns true if the current player owns the TownHall at this base.
func (base *Base) IsSelfOwned() bool {
	return !base.TownHall.IsNil() && base.TownHall.Alliance == api.Alliance_Self
}

// IsEnemyOwned returns true if the enemy player owns the TownHall at this base.
func (base *Base) IsEnemyOwned() bool {
	return !base.TownHall.IsNil() && base.TownHall.Alliance == api.Alliance_Enemy
}

// IsUnowned returns true if no player owns a TownHall at this base.
func (base *Base) IsUnowned() bool {
	return base.TownHall.IsNil()
}

// Natural returns the closest other base.
func (base *Base) Natural() *Base {
	best, minDist := (*Base)(nil), float32(256*256)
	for _, other := range base.m.Bases {
		if dist := base.WalkDistance(other); dist > 0 && dist < minDist {
			best, minDist = other, dist
		}
	}
	return best
}

// WalkDistance returns the ground pathfinding distances between the bases.
func (base *Base) WalkDistance(other *Base) float32 {
	return base.m.distance(base.i, other.i)
}
