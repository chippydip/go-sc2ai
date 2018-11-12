package botutil

import (
	"testing"

	"github.com/chippydip/go-sc2ai/api"
)

func BenchmarkInterfaces(b *testing.B) {
	u := Units{
		[]*api.Unit{&api.Unit{}},
		[]*api.UnitTypeData{&api.UnitTypeData{}},
	}
	for i := 0; i < b.N; i++ {
		//actions := Actions{}

		//actions.UnitsCommand(&u, ability.Attack)
		u.Choose(func(Unit) bool { return false })
	}
}
