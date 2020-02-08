package search

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
)

// DBSCAN ...
type DBSCAN struct {
	Units map[api.UnitTag]botutil.Unit

	clustered map[api.UnitTag]bool
	neighbors []api.UnitTag

	clusters []UnitCluster
	outliers []botutil.Unit
}

// NewDBSCAN ...
func NewDBSCAN() *DBSCAN {
	return &DBSCAN{
		Units:     map[api.UnitTag]botutil.Unit{},
		clustered: map[api.UnitTag]bool{},
	}
}

func (db *DBSCAN) setup() {
	for k := range db.clustered {
		delete(db.clustered, k)
	}

	for i := range db.clusters {
		db.clusters[i].Clear()
	}

	for i := range db.outliers {
		db.outliers[i] = botutil.Unit{}
	}
	db.outliers = db.outliers[:0]
}

// Cluster ...
func (db *DBSCAN) Cluster(minPts int, eps float32) ([]UnitCluster, []botutil.Unit) {
	db.setup()

	c, eps2 := 0, eps*eps
	for k, u := range db.Units {
		if db.clustered[k] {
			continue
		}

		db.getNeighbors(u.Pos2D(), eps2)
		if len(db.neighbors) < minPts {
			db.outliers = append(db.outliers, u)
			continue
		}

		// new cluster
		if len(db.clusters) <= c {
			db.clusters = append(db.clusters, UnitCluster{})
		}

		cluster := &db.clusters[c]
		c++

		cluster.Add(u)
		db.clustered[u.Tag] = true
		db.addNeighbors(cluster)

		for i := 0; i < len(cluster.Units()); i++ {
			u := cluster.Units()[i]

			db.getNeighbors(u.Pos2D(), eps2)
			if len(db.neighbors) >= minPts {
				db.addNeighbors(cluster)
			}
		}
	}

	return db.clusters[:c], db.outliers
}

func (db *DBSCAN) getNeighbors(pos api.Point2D, eps2 float32) {
	db.neighbors = db.neighbors[:0]
	for k2, v2 := range db.Units {
		if pos.Distance2(v2.Pos2D()) <= eps2 {
			db.neighbors = append(db.neighbors, k2)
		}
	}
}

func (db *DBSCAN) addNeighbors(cluster *UnitCluster) {
	for _, k := range db.neighbors {
		if !db.clustered[k] {
			cluster.Add(db.Units[k])
			db.clustered[k] = true
		}
	}
}
