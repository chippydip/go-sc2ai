package botutil

import (
	"github.com/chippydip/go-sc2ai/client"
)

// Bot ...
type Bot struct {
	client.AgentInfo

	*Player
	*UnitContext
	*Actions
	*Builder
}

// NewBot ...
func NewBot(info client.AgentInfo) Bot {
	bot := Bot{AgentInfo: info}

	bot.Player = NewPlayer(info)
	bot.Actions = NewActions(info)
	bot.UnitContext = NewUnitContext(info, bot.Actions)
	bot.Builder = NewBuilder(info, bot.Player, bot.UnitContext)

	return bot
}

// NewBotTemp creates a new temporariy bot without registering for step updates.
func NewBotTemp(info client.AgentInfo) Bot {
	return NewBot(&tempAgentInfo{info})
}

// tempAgentInfo is a wrapper that disables event registration.
type tempAgentInfo struct {
	client.AgentInfo
}

func (info *tempAgentInfo) OnBeforeStep(func()) {}
func (info *tempAgentInfo) OnAfterStep(func())  {}
