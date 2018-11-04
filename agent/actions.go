package agent

import (
	"github.com/chippydip/go-sc2ai/api"
)

func (a *Agent) sendActions() {
	if len(a.actions) > 0 {
		// TODO: can we automatically call this from the client via a callback when in step mode?
		// TODO: should this be async?
		a.info.SendActions(a.actions)
		a.actions = nil
	}
}

func (a *Agent) ChatAll(msg string) {
	a.actions = append(a.actions, &api.Action{
		ActionChat: &api.ActionChat{
			Channel: api.ActionChat_Broadcast,
			Message: msg,
		},
	})
}

func (a *Agent) ChatTeam(msg string) {
	a.actions = append(a.actions, &api.Action{
		ActionChat: &api.ActionChat{
			Channel: api.ActionChat_Team,
			Message: msg,
		},
	})
}

func (a *Agent) unitCommand(cmd *api.ActionRawUnitCommand) {
	a.actions = append(a.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: cmd,
			},
		},
	})
}

func (a *Agent) UnitCommand(unitTag api.UnitTag, ability api.AbilityID) {
	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  []api.UnitTag{unitTag},
	})
}

func (a *Agent) UnitCommandAtTarget(unitTag api.UnitTag, ability api.AbilityID, target api.UnitTag) {
	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  []api.UnitTag{unitTag},
		Target: &api.ActionRawUnitCommand_TargetUnitTag{
			TargetUnitTag: target,
		},
	})
}

func (a *Agent) UnitCommandAtPos(unitTag api.UnitTag, ability api.AbilityID, target api.Point2D) {
	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  []api.UnitTag{unitTag},
		Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
			TargetWorldSpacePos: &target,
		},
	})
}

func (a *Agent) UnitsCommand(unitTags []api.UnitTag, ability api.AbilityID) {
	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
	})
}

func (a *Agent) UnitsCommandAtTarget(unitTags []api.UnitTag, ability api.AbilityID, target api.UnitTag) {
	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
		Target: &api.ActionRawUnitCommand_TargetUnitTag{
			TargetUnitTag: target,
		},
	})
}

func (a *Agent) UnitsCommandAtPos(unitTags []api.UnitTag, ability api.AbilityID, target api.Point2D) {
	a.unitCommand(&api.ActionRawUnitCommand{
		AbilityId: ability,
		UnitTags:  unitTags,
		Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
			TargetWorldSpacePos: &target,
		},
	})
}
