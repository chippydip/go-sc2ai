package main

import "github.com/chippydip/go-sc2ai/api"

type UnitsByTypes map[api.UnitTypeID]Units

func (ut UnitsByTypes) Add(utype api.UnitTypeID, unit *Unit) {
	if unit != nil {
		ut[utype] = append(ut[utype], unit)
	}
}

func (ut UnitsByTypes) AddFromApi(utype api.UnitTypeID, unit *api.Unit) {
	if unit != nil {
		ut[utype] = append(ut[utype], NewUnit(unit))
	}
}

func (ut UnitsByTypes) OfType(ids ...api.UnitTypeID) Units {
	u := Units{}
	for _, id := range ids {
		u = append(u, ut[id]...)
	}
	return u
}

func (ut UnitsByTypes) Units() Units {
	u := Units{}
	for _, units := range ut {
		u = append(u, units...)
	}
	return u
}
