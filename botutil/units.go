package botutil

import (
	"math"

	"github.com/chippydip/go-sc2ai/api"
)

// Units ...
type Units struct {
	ctx    *UnitContext
	raw    []Unit
	filter func(Unit) bool
}

// Init ...
func (units *Units) Init(ctx *UnitContext, capacity int) {
	*units = Units{ctx, make([]Unit, 0, capacity), nil}
}

func (units *Units) applyFilter() {
	if len(units.raw) == 0 {
		units.raw = nil // make sure it's actually nil
	} else if units.filter != nil {
		raw := make([]Unit, 0, len(units.raw))
		for _, u := range units.raw {
			if units.filter(u) {
				raw = append(raw, u)
			}
		}
		if len(raw) == 0 {
			raw = nil
		}
		units.raw = raw
	}
	units.filter = nil
}

func (units *Units) ensureOwns() {
	units.applyFilter()

	if len(units.raw) == 0 {
		return // we may not even have a ctx to compare to
	}

	// Don't mess the the ctx's slice
	if sliceID(units.raw) == sliceID(units.ctx.wrapped) {
		tmp := make([]Unit, len(units.raw), 2*len(units.raw))
		copy(tmp, units.raw)
		units.raw = tmp
	}
}

func sliceID(s []Unit) *Unit {
	if cap(s) > 0 {
		return &s[:cap(s)][cap(s)-1]
	}
	return nil
}

// Len returns the length of the underlying slice of units.
func (units *Units) Len() int {
	units.applyFilter()
	return len(units.raw)
}

// Raw returns the underlying slice of api Units.
func (units *Units) Raw() []Unit {
	units.ensureOwns()
	return units.raw
}

// Data returns the underlying UnitType data.
func (units *Units) Data() []*api.UnitTypeData {
	if units.ctx == nil {
		return nil
	}
	return units.ctx.data
}

// Append adds the given unit to the slice.
func (units *Units) Append(u Unit) {
	units.ensureOwns()
	units.ctx = u.ctx // in case unitx.ctx was nil
	units.raw = append(units.raw, u)
}

// Concat ...
func (units *Units) Concat(other *Units) {
	if len(units.raw) == 0 {
		*units = *other
	} else if len(other.raw) > 0 {
		units.ensureOwns()
		other.applyFilter()
		if len(other.raw) > 0 {
			units.ctx = other.ctx // in case units.ctx was nil
			units.raw = append(units.raw, other.raw...)
		}
	}
}

// Tags ...
func (units *Units) Tags() []api.UnitTag {
	tags := make([]api.UnitTag, 0, len(units.raw))
	for _, u := range units.raw {
		if units.filter == nil || units.filter(u) {
			tags = append(tags, u.Tag)
		}
	}
	return tags
}

// Each ...
func (units Units) Each(f func(Unit)) {
	for _, u := range units.raw {
		if units.filter == nil || units.filter(u) {
			f(u)
		}
	}
}

// EachWhile calls f until it returns false or runs out of elements.
// Returns the last result of f (false on early return).
func (units Units) EachWhile(f func(Unit) bool) bool {
	for _, u := range units.raw {
		if units.filter == nil || units.filter(u) {
			if !f(u) {
				return false
			}
		}
	}
	return true
}

// EachUntil calls f until it returns true or runs out of elements.
// Returns the last result of f (true on early return).
func (units Units) EachUntil(f func(Unit) bool) bool {
	for _, u := range units.raw {
		if units.filter == nil || units.filter(u) {
			if f(u) {
				return true
			}
		}
	}
	return false
}

// Choose returns a new list with only the units for which filter returns true.
func (units Units) Choose(filter func(Unit) bool) Units {
	// Don't filter empty lists
	if len(units.raw) == 0 {
		return units
	}

	// If this is the first filter, just set and return it
	if units.filter == nil {
		return Units{units.ctx, units.raw, filter}
	}

	// Otherwise, union the two filters
	prev := units.filter
	return Units{units.ctx, units.raw, func(u Unit) bool {
		return prev(u) && filter(u)
	}}
}

// Drop returns a new list without the units for which filter returns true.
func (units Units) Drop(filter func(Unit) bool) Units {
	return units.Choose(func(u Unit) bool {
		return !filter(u)
	})
}

// First returns the first unit in the list or a Unit.IsNil() if the list is empty.
func (units Units) First() Unit {
	if len(units.raw) > 0 {
		if units.filter == nil {
			return units.raw[0]
		}

		for _, u := range units.raw {
			if units.filter(u) {
				return u
			}
		}
	}
	return Unit{}
}

// ClosestTo returns the closest unit from the latest observation.
func (units Units) ClosestTo(pos api.Point2D) Unit {
	minDist := float32(math.Inf(1))
	var closest Unit
	for _, u := range units.raw {
		if units.filter == nil || units.filter(u) {
			dist := pos.Distance2(u.Pos.ToPoint2D())
			if dist < minDist {
				closest = u
				minDist = dist
			}
		}
	}
	return closest
}

// Tagged ...
func (units Units) Tagged(m map[api.UnitTag]bool) Units {
	return units.Choose(func(u Unit) bool {
		return m[u.Tag]
	})
}

// NotTagged ...
func (units Units) NotTagged(m map[api.UnitTag]bool) Units {
	return units.Choose(func(u Unit) bool {
		return !m[u.Tag]
	})
}

// HasEnergy ...
func (units Units) HasEnergy(energy float32) Units {
	return units.Choose(func(u Unit) bool {
		return u.Energy >= energy
	})
}

// IsBuilt ...
func (units Units) IsBuilt() Units {
	return units.Choose(func(u Unit) bool {
		return u.BuildProgress == 1
	})
}

// NoBuff ...
func (units Units) NoBuff(buffID api.BuffID) Units {
	return units.Choose(func(u Unit) bool {
		for _, b := range u.BuffIds {
			if b == buffID {
				return false
			}
		}
		return true
	})
}
