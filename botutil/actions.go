package botutil

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
)

// agentInfo is the subset of the client.AgentInfo interface required by Actions.
type agentInfo interface {
	OnBeforeStep(func())
	SendActions([]*api.Action) []api.ActionResult
}

// Actions provides convenience methods for queueing actions to be sent in a batch.
type Actions struct {
	info    agentInfo
	actions []*api.Action
}

// NewActions creates a new Actions manager. It's Send() method is registered to be
// automatcially called before each client Step().
func NewActions(info agentInfo) *Actions {
	a := &Actions{info: info}
	info.OnBeforeStep(a.Send)
	return a
}

// Send is called automatically to submit queued actions before each Step(). It may also be
// called manually at any point to send all queued actions immediately.
func (a *Actions) Send() {
	if len(a.actions) == 0 {
		return
	}

	for i, r := range a.info.SendActions(a.actions) {
		if r != api.ActionResult_Success {
			log.Print("ActionError: ", r, a.actions[i])
		}
	}
	a.actions = nil
}

// Chat sends a message that all players can see.
func (a *Actions) Chat(msg string) {
	a.actions = append(a.actions, &api.Action{
		ActionChat: &api.ActionChat{
			Channel: api.ActionChat_Broadcast,
			Message: msg,
		},
	})
}

// ChatTeam sends a message that only teammates (and observers) can see.
func (a *Actions) ChatTeam(msg string) {
	a.actions = append(a.actions, &api.Action{
		ActionChat: &api.ActionChat{
			Channel: api.ActionChat_Team,
			Message: msg,
		},
	})
}

// MoveCamera repositions the camera to center on the target point.
func (a *Actions) MoveCamera(pt api.Point2D) {
	p := pt.ToPoint()
	a.actions = append(a.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_CameraMove{
				CameraMove: &api.ActionRawCameraMove{
					CenterWorldSpace: &p,
				},
			},
		},
	})
}

// getTagger is the part of api.Unit's api required to send commands.
// An interface is used here so that custom unit types can also be used
// here as long as they implement this method.
type getTagger interface {
	GetTag() api.UnitTag
}

// soloTagser implements the tagser interface for a single unit. This is
// done by wrapping the tag so that no additional allocations are required.
type soloTagser api.UnitTag

func (t soloTagser) Tags() []api.UnitTag {
	return []api.UnitTag{api.UnitTag(t)}
}

// UnitCommand orders a unit to use an ability.
func (a *Actions) UnitCommand(u getTagger, ability api.AbilityID) {
	a.UnitsCommand(soloTagser(u.GetTag()), ability)
}

// UnitCommandOnTarget orders a unit to use an ability on a target unit.
func (a *Actions) UnitCommandOnTarget(u getTagger, ability api.AbilityID, target getTagger) {
	a.UnitsCommandOnTarget(soloTagser(u.GetTag()), ability, target)
}

// UnitCommandAtPos orders a unit to use an ability at a target location.
func (a *Actions) UnitCommandAtPos(u getTagger, ability api.AbilityID, target *api.Point2D) {
	a.UnitsCommandAtPos(soloTagser(u.GetTag()), ability, target)
}

// tagser provides access to a slice of unit tags to issue orders to.
type tagser interface {
	Tags() []api.UnitTag
}

// UnitsCommand orders units to all use an ability.
func (a *Actions) UnitsCommand(units tagser, ability api.AbilityID) {
	unitTags := units.Tags()
	if len(unitTags) == 0 {
		return
	}

	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
	})
}

// UnitsCommandOnTarget orders units to all use an ability on a target unit.
func (a *Actions) UnitsCommandOnTarget(units tagser, ability api.AbilityID, target getTagger) {
	unitTags := units.Tags()
	if len(unitTags) == 0 {
		return
	}

	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
		Target: &api.ActionRawUnitCommand_TargetUnitTag{
			TargetUnitTag: target.GetTag(),
		},
	})
}

// UnitsCommandAtPos orders units to all use an ability at a target location.
func (a *Actions) UnitsCommandAtPos(units tagser, ability api.AbilityID, target *api.Point2D) {
	unitTags := units.Tags()
	if len(unitTags) == 0 {
		return
	}

	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
		Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
			TargetWorldSpacePos: target,
		},
	})
}

// unitCommand finishes wrapping an ActionRawUnitCommand and adds it to the command list.
func (a *Actions) unitCommand(cmd *api.ActionRawUnitCommand) {
	a.actions = append(a.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: cmd,
			},
		},
	})
}
