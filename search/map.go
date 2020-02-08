package search

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
)

// Map ...
type Map struct {
	bases

	StartLocation api.Point2D
	// EnemyStartLocations []api.Point2D
}

// NewMap ...
func NewMap(bot *botutil.Bot) *Map {
	m := &Map{}
	m.bases = newBases(m, bot)
	m.StartLocation = bot.Self.Structures().First().Pos2D()

	// locs := bot.GameInfo().GetStartRaw().GetStartLocations()
	// m.EnemyStartLocations = make([]api.Point2D, len(locs))
	// for i, l := range locs {
	// 	m.EnemyStartLocations[i] = *l
	// }

	update := func() {
		m.bases.update(bot)

		// if len(m.EnemyStartLocations) > 1 {
		// 	bot.
		// }
	}
	bot.OnAfterStep(update)
	update()

	return m
}
