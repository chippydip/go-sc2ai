package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

// Query ...
type Query struct {
	info    client.AgentInfo
	request api.RequestQuery
}

// NewQuery ...
func NewQuery(info client.AgentInfo) Query {
	return Query{info: info}
}

// IgnoreResourceRequirements ...
func (q *Query) IgnoreResourceRequirements() {
	q.request.IgnoreResourceRequirements = true
}

// PathingFrom ...
func (q *Query) PathingFrom(start, end api.Point2D) int {
	q.request.Pathing = append(q.request.Pathing, &api.RequestQueryPathing{
		Start: &api.RequestQueryPathing_StartPos{
			StartPos: &start,
		},
		EndPos: &end,
	})
	return len(q.request.Pathing) - 1
}

// UnitPathing ...
func (q *Query) UnitPathing(tag api.UnitTag, end api.Point2D) int {
	q.request.Pathing = append(q.request.Pathing, &api.RequestQueryPathing{
		Start: &api.RequestQueryPathing_UnitTag{
			UnitTag: tag,
		},
		EndPos: &end,
	})
	return len(q.request.Pathing) - 1
}

// UnitAbilities ...
func (q *Query) UnitAbilities(tag api.UnitTag) int {
	q.request.Abilities = append(q.request.Abilities, &api.RequestQueryAvailableAbilities{
		UnitTag: tag,
	})
	return len(q.request.Abilities) - 1
}

// Placement ...
func (q *Query) Placement(ability api.AbilityID, pos api.Point2D) int {
	q.request.Placements = append(q.request.Placements, &api.RequestQueryBuildingPlacement{
		AbilityId: ability,
		TargetPos: &pos,
	})
	return len(q.request.Placements) - 1
}

// PlacementWithUnit ...
func (q *Query) PlacementWithUnit(tag api.UnitTag, ability api.AbilityID, pos api.Point2D) int {
	q.request.Placements = append(q.request.Placements, &api.RequestQueryBuildingPlacement{
		AbilityId:      ability,
		TargetPos:      &pos,
		PlacingUnitTag: tag,
	})
	return len(q.request.Placements) - 1
}

// Execute runs the query and returns the result
func (q *Query) Execute() QueryResult {
	return QueryResult{
		request:  q.request,
		response: q.info.Query(q.request),
	}
}

// QueryResult ...
type QueryResult struct {
	request  api.RequestQuery
	response *api.ResponseQuery
}

// Pathing ...
func (q QueryResult) Pathing() []*api.ResponseQueryPathing {
	return q.response.GetPathing()
}

// PathingQuery ...
func (q QueryResult) PathingQuery(i int) *api.RequestQueryPathing {
	return q.request.Pathing[i]
}

// Abilities ...
func (q QueryResult) Abilities() []*api.ResponseQueryAvailableAbilities {
	return q.response.GetAbilities()
}

// AbilitiesQuery ...
func (q QueryResult) AbilitiesQuery(i int) *api.RequestQueryAvailableAbilities {
	return q.request.Abilities[i]
}

// Placements ...
func (q QueryResult) Placements() []*api.ResponseQueryBuildingPlacement {
	return q.response.GetPlacements()
}

// PlacementQuery ...
func (q QueryResult) PlacementQuery(i int) *api.RequestQueryBuildingPlacement {
	return q.request.Placements[i]
}
