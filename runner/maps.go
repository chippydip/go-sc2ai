package runner

import (
	"math/rand"
	"time"
)

var maps2018s3 = []string{
	"AcidPlantLE",
	"BlueshiftLE",
	"CeruleanFallLE",
	"DreamcatcherLE",
	"FractureLE",
	"LostAndFoundLE",
	"ParaSiteLE",
}

var maps2018s4 = []string{
	"AutomatonLE",
	"BlueshiftLE",
	"CeruleanFallLE",
	"DarknessSanctuaryLE",
	"KairosJunctionLE",
	"ParaSiteLE",
	"PortAleksanderLE",
}

// TODO: check for current ladder pool maps, download if missing?

// Random1v1Map returns a random map name from the current 1v1 ladder map pool.
func Random1v1Map() string {
	currentMaps := maps2018s4

	rand.Seed(time.Now().UnixNano())
	return currentMaps[rand.Intn(len(currentMaps))] + ".SC2Map"
}
