package main

import (
	"log"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

func main() {
	// Set some config values for single player (they can be overriden by command line flags)
	runner.Set("map", "InterloperLE.SC2Map")
	runner.Set("ComputerOpponent", "true")
	runner.Set("ComputerRace", "protoss")
	runner.Set("ComputerDifficulty", "Easy")

	// Create the agent and then start the game
	runner.RunAgent(client.NewParticipant(api.Race_Terran, client.AgentFunc(runBot), "NilBot"))
}

type bot struct {
	client.AgentInfo
}

func runBot(info client.AgentInfo) {
	bot := bot{info}
	bot.init()

	for bot.IsInGame() {
		bot.update()

		bot.Step(1)
	}

	// Alread out of the game at this point, so can't send this as a chat message
	log.Print("gg")
}

func (bot *bot) init() {
	// Send a friendly hello
	bot.SendActions([]*api.Action{
		&api.Action{
			ActionChat: &api.ActionChat{
				Channel: api.ActionChat_Broadcast,
				Message: "(glhf)",
			},
		},
	})
}

// OnStep is called each game step (every game update by defaul)
func (bot *bot) update() {
	// Echo chat to the console
	for _, chat := range bot.Observation().GetChat() {
		log.Printf("[%v] %v\n", chat.GetPlayerId(), chat.GetMessage())
	}
}
