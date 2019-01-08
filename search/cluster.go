package search

import (
	"math"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
)

// Cluster breaks a list of units into clusters based on the given clustering distance.
func Cluster(units botutil.Units, distance float32) []UnitCluster {
	maxDistance := distance * distance

	// TODO: replace this with a better algorithm
	var clusters []UnitCluster
	units.Each(func(u botutil.Unit) {
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
	})
	return clusters
}

// UnitCluster is a cluster of units and the associated center of mass.
type UnitCluster struct {
	sum   api.Vec2D
	units botutil.Units
}

// Add adds a new unit to the cluster and updates the center of mass.
func (c *UnitCluster) Add(u botutil.Unit) {
	c.sum = c.sum.Add(api.Vec2D(u.Pos2D()))
	c.units.Append(u)
}

// Center is the center of mass of the cluster.
func (c *UnitCluster) Center() api.Point2D {
	return api.Point2D(c.sum.Div(float32(c.units.Len())))
}

// Units is the list of units in the cluster.
func (c *UnitCluster) Units() botutil.Units {
	return c.units
}
