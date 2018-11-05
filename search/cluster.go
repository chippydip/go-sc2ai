package search

import (
	"math"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/filter"
)

// Cluster breaks a list of units into clusters based on the given clustering distance.
func Cluster(units filter.Units, distance float32) []UnitCluster {
	maxDistance := distance * distance

	// TODO: replace this with a better algorithm
	var clusters []UnitCluster
	for _, u := range units {
		// Find the nearest cluster
		minDist := float32(math.MaxFloat32)
		clusterIndex := -1
		for i, cluster := range clusters {
			d := u.Pos.ToPoint2D().Distance2(cluster.Center())
			if d < minDist {
				minDist = d
				clusterIndex = i
			}
		}

		// If too far, add a new cluster
		if minDist > maxDistance || clusterIndex < 0 {
			clusterIndex = len(clusters)
			clusters = append(clusters, UnitCluster{})
		}

		clusters[clusterIndex].Add(u)
	}
	return clusters
}

// UnitCluster is a cluster of units and the associated center of mass.
type UnitCluster struct {
	sum   api.Vec2D
	units filter.Units
}

// Add adds a new unit to the cluster and updates the center of mass.
func (c *UnitCluster) Add(u *api.Unit) {
	c.sum = c.sum.Add(api.Vec2D(u.Pos.ToPoint2D()))
	c.units = append(c.units, u)
}

// Center is the center of mass of the cluster.
func (c *UnitCluster) Center() api.Point2D {
	return api.Point2D(c.sum.Div(float32(len(c.units))))
}

// Units is the list of units in the cluster.
func (c *UnitCluster) Units() filter.Units {
	return c.units
}
