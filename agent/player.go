package agent

import "github.com/chippydip/go-sc2ai/api"

// PlayerId        PlayerID `protobuf:"varint,1,opt,name=player_id,json=playerId,proto3,casttype=PlayerID" json:"player_id,omitempty"`
// Minerals        uint32   `protobuf:"varint,2,opt,name=minerals,proto3" json:"minerals,omitempty"`
// Vespene         uint32   `protobuf:"varint,3,opt,name=vespene,proto3" json:"vespene,omitempty"`
// FoodCap         uint32   `protobuf:"varint,4,opt,name=food_cap,json=foodCap,proto3" json:"food_cap,omitempty"`
// FoodUsed        uint32   `protobuf:"varint,5,opt,name=food_used,json=foodUsed,proto3" json:"food_used,omitempty"`
// FoodArmy        uint32   `protobuf:"varint,6,opt,name=food_army,json=foodArmy,proto3" json:"food_army,omitempty"`
// FoodWorkers     uint32   `protobuf:"varint,7,opt,name=food_workers,json=foodWorkers,proto3" json:"food_workers,omitempty"`
// IdleWorkerCount uint32   `protobuf:"varint,8,opt,name=idle_worker_count,json=idleWorkerCount,proto3" json:"idle_worker_count,omitempty"`
// ArmyCount       uint32   `protobuf:"varint,9,opt,name=army_count,json=armyCount,proto3" json:"army_count,omitempty"`
// WarpGateCount   uint32   `protobuf:"varint,10,opt,name=warp_gate_count,json=warpGateCount,proto3" json:"warp_gate_count,omitempty"`
// LarvaCount      uint32   `protobuf:"varint,11,opt,name=larva_count,json=larvaCount,proto3" json:"larva_count,omitempty"`

func (a *Agent) PlayerID() api.PlayerID {
	return a.info.PlayerID()
}

func (a *Agent) playerCommon() *api.PlayerCommon {
	return a.info.Observation().Observation.PlayerCommon
}

func (a *Agent) Minerals() int {
	return int(a.playerCommon().Minerals)
}

func (a *Agent) Vespene() int {
	return int(a.playerCommon().Vespene)
}

func (a *Agent) FoodCap() int {
	return int(a.playerCommon().FoodCap)
}

func (a *Agent) FoodUsed() int {
	return int(a.playerCommon().FoodUsed)
}

func (a *Agent) FoodArmy() int {
	return int(a.playerCommon().FoodArmy)
}

func (a *Agent) FoodWorkers() int {
	return int(a.playerCommon().FoodWorkers)
}
