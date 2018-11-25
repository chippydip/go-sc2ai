package botutil

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

// Actions provides convenience methods for queueing actions to be sent in a batch.
type Actions struct {
	info    client.AgentInfo
	actions []*api.Action
}

// NewActions creates a new Actions manager. It's Send() method is registered to be
// automatcially called before each client Step().
func NewActions(info client.AgentInfo) *Actions {
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

// UnitOrder orders a unit to use an ability.
func (a *Actions) UnitOrder(u getTagger, ability api.AbilityID) {
	a.unitsOrder([]api.UnitTag{u.GetTag()}, ability)
}

// UnitOrderTarget orders a unit to use an ability on a target unit.
func (a *Actions) UnitOrderTarget(u getTagger, ability api.AbilityID, target getTagger) {
	a.unitsOrderTarget([]api.UnitTag{u.GetTag()}, ability, target)
}

// UnitOrderPos orders a unit to use an ability at a target location.
func (a *Actions) UnitOrderPos(u getTagger, ability api.AbilityID, target *api.Point2D) {
	a.unitsOrderPos([]api.UnitTag{u.GetTag()}, ability, target)
}

// tagser provides access to a slice of unit tags to issue orders to.
type tagser interface {
	Tags() []api.UnitTag
}

// UnitsOrder orders units to all use an ability.
func (a *Actions) UnitsOrder(units tagser, ability api.AbilityID) {
	a.unitsOrder(units.Tags(), ability)
}

// UnitsOrderTarget orders units to all use an ability on a target unit.
func (a *Actions) UnitsOrderTarget(units tagser, ability api.AbilityID, target getTagger) {
	a.unitsOrderTarget(units.Tags(), ability, target)
}

// UnitsOrderPos orders units to all use an ability at a target location.
func (a *Actions) UnitsOrderPos(units tagser, ability api.AbilityID, target *api.Point2D) {
	a.unitsOrderPos(units.Tags(), ability, target)
}

// unitsOrder orders units to all use an ability.
func (a *Actions) unitsOrder(unitTags []api.UnitTag, ability api.AbilityID) {
	if len(unitTags) == 0 {
		return
	}

	a.unitOrder(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
	})
}

// unitsOrderTarget orders units to all use an ability on a target unit.
func (a *Actions) unitsOrderTarget(unitTags []api.UnitTag, ability api.AbilityID, target getTagger) {
	if len(unitTags) == 0 {
		return
	}

	a.unitOrder(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
		Target: &api.ActionRawUnitCommand_TargetUnitTag{
			TargetUnitTag: target.GetTag(),
		},
	})
}

// unitsOrderPos orders units to all use an ability at a target location.
func (a *Actions) unitsOrderPos(unitTags []api.UnitTag, ability api.AbilityID, target *api.Point2D) {
	if len(unitTags) == 0 {
		return
	}

	a.unitOrder(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
		Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
			TargetWorldSpacePos: target,
		},
	})
}

// unitOrder finishes wrapping an ActionRawUnitCommand and adds it to the command list.
func (a *Actions) unitOrder(cmd *api.ActionRawUnitCommand) {
	a.actions = append(a.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: cmd,
			},
		},
	})
}

// Convenience methods for giving orders directly to units:

// Order ...
func (units Units) Order(ability api.AbilityID) {
	if len(units.raw) > 0 {
		units.ctx.actions.unitsOrder(units.Tags(), ability)
	}
}

// OrderTarget ...
func (units Units) OrderTarget(ability api.AbilityID, target getTagger) {
	if len(units.raw) > 0 {
		units.ctx.actions.unitsOrderTarget(units.Tags(), ability, target)
	}
}

// OrderPos ...
func (units Units) OrderPos(ability api.AbilityID, target *api.Point2D) {
	if len(units.raw) > 0 {
		units.ctx.actions.unitsOrderPos(units.Tags(), ability, target)
	}
}

// Order ...
func (u Unit) Order(ability api.AbilityID) {
	if !u.IsNil() {
		u.ctx.actions.UnitOrder(u, ability)
	}
}

// OrderTarget ...
func (u Unit) OrderTarget(ability api.AbilityID, target getTagger) {
	if !u.IsNil() {
		u.ctx.actions.UnitOrderTarget(u, ability, target)
	}
}

// OrderPos ...
func (u Unit) OrderPos(ability api.AbilityID, target *api.Point2D) {
	if !u.IsNil() {
		u.ctx.actions.UnitOrderPos(u, ability, target)
	}
}
