package main

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/botutil"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

type bot struct {
	*botutil.Bot
}

func main() {
	// Play a random map against a medium difficulty computer
	runner.Set("map", runner.Random1v1Map())
	runner.Set("ComputerOpponent", "true")
	runner.Set("ComputerRace", api.Race_Random.String())
	runner.Set("ComputerDifficulty", api.Difficulty_Easy.String())

	// Create the agent and then start the game
	agent := client.AgentFunc(runAgent)
	runner.RunAgent(client.NewParticipant(api.Race_Protoss, agent, "StubBot"))
}

func runAgent(info client.AgentInfo) {
	bot := bot{Bot: botutil.NewBot(info)}
	bot.LogActionErrors()

	bot.init()
	for bot.IsInGame() {
		bot.doSmt()

		if err := bot.Step(1); err != nil {
			log.Print(err)
			break
		}
	}
}

func (bot *bot) init() {
	// Send a friendly hello
	bot.Chat("(glhf)")
}

func (bot *bot) doSmt() {

}
