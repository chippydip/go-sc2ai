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
					UnitTags:  []uint64{uint64(unit.Tag)}, // UnitTag should be accepted here
				}}}})
}

func (bot *zergRush) unitCommandTargetTag(unit *api.Unit, ability abilityType.Ability, target api.UnitTag) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []uint64{uint64(unit.Tag)}, // UnitTag should be accepted here
					Target: &api.ActionRawUnitCommand_TargetUnitTag{
						TargetUnitTag: target,
					}}}}})
}

func (bot *zergRush) unitCommandTargetPos(unit *api.Unit, ability abilityType.Ability, target *Point) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []uint64{uint64(unit.Tag)}, // UnitTag should be accepted here
					Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
						TargetWorldSpacePos: target.To2D(),
					}}}}})
}

func (bot *zergRush) unitsCommandTargetPos(units []*api.Unit, ability abilityType.Ability, target *Point) {
	// I hope, we can avoid this conversion in future
	uTags := []uint64{}
	for _, unit := range units {
		uTags = append(uTags, uint64(unit.Tag))
	}
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  uTags,
					Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
						TargetWorldSpacePos: target.To2D(),
					}}}}})
}
