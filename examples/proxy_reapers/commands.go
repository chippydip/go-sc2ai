package main

import (
	"github.com/chippydip/go-sc2ai/api"
)

func (bot *proxyReapers) chatSend(msg string) {
	bot.info.SendActions([]*api.Action{{
		ActionChat: &api.ActionChat{
			Channel: api.ActionChat_Broadcast,
			Message: msg,
		},
	}})
}

func (bot *proxyReapers) unitCommand(unit *Unit, ability api.AbilityID) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []api.UnitTag{unit.Tag},
				}}}})
}

func (bot *proxyReapers) unitCommandTargetTag(unit *Unit, ability api.AbilityID, target api.UnitTag, queue bool) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId:    ability,
					UnitTags:     []api.UnitTag{unit.Tag},
					QueueCommand: queue,
					Target: &api.ActionRawUnitCommand_TargetUnitTag{
						TargetUnitTag: target,
					}}}}})
}

func (bot *proxyReapers) unitCommandTargetPos(unit *Unit, ability api.AbilityID, target api.Point2D, queue bool) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId:    ability,
					UnitTags:     []api.UnitTag{unit.Tag},
					QueueCommand: queue,
					Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
						TargetWorldSpacePos: &target,
					}}}}})
}

func (bot *proxyReapers) unitsCommandTargetPos(units Units, ability api.AbilityID, target api.Point2D) {
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
