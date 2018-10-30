package main

import abilityType "github.com/chippydip/go-sc2ai/api/ability"
import "github.com/chippydip/go-sc2ai/api"

func (bot *zergRush) ChatSend(msg string) {
	bot.SendActions([]*api.Action{{
		ActionChat: &api.ActionChat{
			Channel: api.ActionChat_Broadcast,
			Message: msg,
		},
	}})
}

func (bot *zergRush) unitCommand(unit *api.Unit, ability abilityType.Ability) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []api.UnitTag{unit.Tag},
				}}}})
}

func (bot *zergRush) unitCommandTargetTag(unit *api.Unit, ability abilityType.Ability, target api.UnitTag) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []api.UnitTag{unit.Tag},
					Target: &api.ActionRawUnitCommand_TargetUnitTag{
						TargetUnitTag: target,
					}}}}})
}

func (bot *zergRush) unitCommandTargetPos(unit *api.Unit, ability abilityType.Ability, target api.Point2D) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []api.UnitTag{unit.Tag},
					Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
						TargetWorldSpacePos: &target,
					}}}}})
}

func (bot *zergRush) unitsCommandTargetPos(units []*api.Unit, ability abilityType.Ability, target api.Point2D) {
	// I hope, we can avoid this conversion in future
	uTags := []api.UnitTag{}
	for _, unit := range units {
		uTags = append(uTags, unit.Tag)
	}
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  uTags,
					Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
						TargetWorldSpacePos: &target,
					}}}}})
}
