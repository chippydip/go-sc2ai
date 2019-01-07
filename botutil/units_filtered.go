package botutil

import "github.com/chippydip/go-sc2ai/api"

const (
	filterFlying byte = 1 << iota
	filterGround
	filterCanAttack
	filterPassive
	filterUnits
	filterStructures
)

type filteredUnits struct {
	ctx      *UnitContext
	filter   func(Unit) bool
	alliance api.Alliance
	bits     byte
}

// ground passive
//   ground passive units (MULEs, burrowed units?)
//   ground passive structures (buildings, uprooted spines/spores?)
// ground attackers
//   ground attackers structures (static defenses)
//   ground attackers units
// flying attackers
//   flying attackers units
//   flying attackers structures (none?)
// flying passive
//   flying passive structures (lifted Terran buildings)
//   flying passive units (air support units, overlords, etc)

func filterToMask(filter byte) [9]bool {
	// Normalize filter bits
	if filter&filterFlying == 0 && filter&filterGround == 0 {
		filter |= filterFlying | filterGround
	}
	if filter&filterCanAttack == 0 && filter&filterPassive == 0 {
		filter |= filterCanAttack | filterPassive
	}
	if filter&filterUnits == 0 && filter&filterStructures == 0 {
		filter |= filterUnits | filterStructures
	}

	// Include all by default, and then exclude as needed
	r := [9]bool{true, true, true, true, true, true, true, true, false}
	if filter&filterGround == 0 {
		r[0], r[1], r[2], r[3] = false, false, false, false
	}
	if filter&filterFlying == 0 {
		r[4], r[5], r[6], r[7] = false, false, false, false
	}

	if filter&filterPassive == 0 {
		r[0], r[1], r[6], r[7] = false, false, false, false
	}
	if filter&filterCanAttack == 0 {
		r[2], r[3], r[4], r[5] = false, false, false, false
	}

	if filter&filterUnits == 0 {
		r[0], r[3], r[4], r[7] = false, false, false, false
	}
	if filter&filterStructures == 0 {
		r[1], r[2], r[5], r[6] = false, false, false, false
	}
	return r
}

func newFilter(m map[api.UnitTypeID]Units, alliance api.Alliance) filteredUnits {
	return filteredUnits{ctx: m[0].ctx, alliance: alliance}
}

func (f filteredUnits) Flying() filteredUnits     { f.bits |= filterFlying; return f }
func (f filteredUnits) Ground() filteredUnits     { f.bits |= filterGround; return f }
func (f filteredUnits) CanAttack() filteredUnits  { f.bits |= filterCanAttack; return f }
func (f filteredUnits) Passive() filteredUnits    { f.bits |= filterPassive; return f }
func (f filteredUnits) Units() filteredUnits      { f.bits |= filterUnits; return f }
func (f filteredUnits) Structures() filteredUnits { f.bits |= filterStructures; return f }

func (f filteredUnits) Choose(filter func(Unit) bool) filteredUnits {
	if f.filter == nil {
		f.filter = filter
	} else {
		prev := f.filter
		f.filter = func(u Unit) bool {
			return prev(u) && filter(u)
		}
	}
	return f
}

func (f filteredUnits) All() Units {
	var raw []Unit

	include, start, ai := false, 0, 8*allianceIndex(f.alliance)
	for i, ok := range filterToMask(f.bits) {
		if ok == include {
			continue
		}
		include = ok

		if include {
			start = f.ctx.groups[ai+i]
		} else {
			end := f.ctx.groups[ai+i]
			if start == end {
				continue
			}

			s := f.ctx.wrapped[start:end]
			if f.filter != nil {
				// TODO: Check ranges for bulk appends?
				for _, u := range s {
					if f.filter(u) {
						raw = append(raw, u)
					}
				}
			} else if raw != nil {
				raw = append(raw, s...)
			} else {
				raw = s
			}
		}
	}

	return Units{ctx: f.ctx, raw: raw}
}

func (f filteredUnits) First() Unit {
	include, start, ai := false, 0, 8*allianceIndex(f.alliance)
	for i, ok := range filterToMask(f.bits) {
		if ok == include {
			continue
		}
		include = ok

		if include {
			start = f.ctx.groups[ai+i]
		} else {
			end := f.ctx.groups[ai+i]
			if start == end {
				continue
			}

			s := f.ctx.raw[start:end]
			if f.filter != nil {
				for _, u := range s {
					data := f.ctx.data[u.UnitType]
					wrapped := Unit{Unit: u, UnitTypeData: data}
					if f.filter(wrapped) {
						return wrapped
					}
				}
			} else {
				u := s[0]
				data := f.ctx.data[u.UnitType]
				wrapped := Unit{Unit: u, UnitTypeData: data}
				return wrapped
			}
		}
	}

	return Unit{}
}

type self map[api.UnitTypeID]Units

func (m self) Flying() filteredUnits     { return newFilter(m, api.Alliance_Self).Flying() }
func (m self) Ground() filteredUnits     { return newFilter(m, api.Alliance_Self).Ground() }
func (m self) CanAttack() filteredUnits  { return newFilter(m, api.Alliance_Self).CanAttack() }
func (m self) Passive() filteredUnits    { return newFilter(m, api.Alliance_Self).Passive() }
func (m self) Units() filteredUnits      { return newFilter(m, api.Alliance_Self).Units() }
func (m self) Structures() filteredUnits { return newFilter(m, api.Alliance_Self).Structures() }
func (m self) All() Units                { return newFilter(m, api.Alliance_Self).All() }
func (m self) First() Unit               { return newFilter(m, api.Alliance_Self).First() }
func (m self) Choose(filter func(Unit) bool) filteredUnits {
	return newFilter(m, api.Alliance_Self).Choose(filter)
}

func (m self) TechAlias(unitType api.UnitTypeID) Units {
	units, aliases := m[unitType], m[0].ctx.data[unitType].TechAlias
	for _, alias := range aliases {
		other := m[alias]
		units.Concat(&other)
	}
	return units
}

func (m self) Count(unitType api.UnitTypeID) int {
	units := m[unitType]
	return units.Len()
}

func (m self) CountInProduction(unitType api.UnitTypeID) int {
	n, abil := 0, m[0].ctx.data[unitType].AbilityId
	for _, u := range m.All().raw {
		for _, order := range u.Orders {
			if order.AbilityId == abil {
				n++
			}
		}
	}
	return n
}

func (m self) CountAll(unitType api.UnitTypeID) int {
	return m.Count(unitType) + m.CountInProduction(unitType)
}

func (m self) CountIf(predicate func(Unit) bool) int {
	n := 0
	m.All().Each(func(u Unit) {
		if predicate(u) {
			n++
		}
	})
	return n
}

type ally map[api.UnitTypeID]Units

func (m ally) Flying() filteredUnits     { return newFilter(m, api.Alliance_Ally).Flying() }
func (m ally) Ground() filteredUnits     { return newFilter(m, api.Alliance_Ally).Ground() }
func (m ally) CanAttack() filteredUnits  { return newFilter(m, api.Alliance_Ally).CanAttack() }
func (m ally) Passive() filteredUnits    { return newFilter(m, api.Alliance_Ally).Passive() }
func (m ally) Units() filteredUnits      { return newFilter(m, api.Alliance_Ally).Units() }
func (m ally) Structures() filteredUnits { return newFilter(m, api.Alliance_Ally).Structures() }
func (m ally) All() Units                { return newFilter(m, api.Alliance_Ally).All() }
func (m ally) First() Unit               { return newFilter(m, api.Alliance_Ally).First() }
func (m ally) Choose(filter func(Unit) bool) filteredUnits {
	return newFilter(m, api.Alliance_Ally).Choose(filter)
}

type enemy map[api.UnitTypeID]Units

func (m enemy) Flying() filteredUnits     { return newFilter(m, api.Alliance_Enemy).Flying() }
func (m enemy) Ground() filteredUnits     { return newFilter(m, api.Alliance_Enemy).Ground() }
func (m enemy) CanAttack() filteredUnits  { return newFilter(m, api.Alliance_Enemy).CanAttack() }
func (m enemy) Passive() filteredUnits    { return newFilter(m, api.Alliance_Enemy).Passive() }
func (m enemy) Units() filteredUnits      { return newFilter(m, api.Alliance_Enemy).Units() }
func (m enemy) Structures() filteredUnits { return newFilter(m, api.Alliance_Enemy).Structures() }
func (m enemy) All() Units                { return newFilter(m, api.Alliance_Enemy).All() }
func (m enemy) First() Unit               { return newFilter(m, api.Alliance_Enemy).First() }
func (m enemy) Choose(filter func(Unit) bool) filteredUnits {
	return newFilter(m, api.Alliance_Enemy).Choose(filter)
}

type neutral map[api.UnitTypeID]Units

func newNeutral(m neutral, start, length int) Units {
	ctx := m[0].ctx
	start, end := ctx.groups[start], ctx.groups[start+length]
	return Units{ctx: ctx, raw: ctx.wrapped[start:end]}
}

func (m neutral) Minerals() Units  { return newNeutral(m, neutralMinerals, 1) }
func (m neutral) Vespene() Units   { return newNeutral(m, neutralVespene, 1) }
func (m neutral) Resources() Units { return newNeutral(m, neutralResources, 2) }
func (m neutral) All() Units       { return newNeutral(m, neutralAll, 3) }
