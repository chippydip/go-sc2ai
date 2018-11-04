package filter

import (
	"github.com/chippydip/go-sc2ai/api"
)

// Units ...
type Units []*api.Unit

// Tags ...
func (units Units) Tags() []api.UnitTag {
	tags := make([]api.UnitTag, len(units))
	for i, u := range units {
		tags[i] = u.Tag
	}
	return tags
}

// Filter ...
func (units Units) Filter(filter func(*api.Unit) bool) Units {
	var result []*api.Unit
	for _, u := range units {
		if filter(u) {
			result = append(result, u)
		}
	}
	return result
}
