package main

import (
	"fmt"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

func main() {
	// Set some config values for single player (they can be overriden by command line flags)
	runner.Set("map", "C:/Program Files (x86)/StarCraft II/Maps/InterloperLE.SC2Map")
	runner.Set("ComputerOpponent", "true")
	runner.Set("ComputerRace", "protoss")
	runner.Set("ComputerDifficulty", "Easy")

	// Create the agent and then start the game
	runner.RunAgent(client.NewParticipant(api.Race_Terran, &exampleBot{}))
}

type exampleBot struct {
	client.AgentInfo
}

// OnGameStart is called once at the start of the game
func (bot *exampleBot) OnGameStart(info client.AgentInfo) {
	bot.AgentInfo = info

	// Send a friendly hello
	bot.SendActions([]*api.Action{
		&api.Action{
			ActionChat: &api.ActionChat{
				Channel: api.ActionChat_Broadcast,
				Message: "gl hf",
			},
		},
	})
}

// OnStep is called each game step (every game update by defaul)
func (bot *exampleBot) OnStep() {
	// Echo chat to the console
	for _, chat := range bot.Observation().GetChat() {
		fmt.Printf("[%v] %v\n", chat.GetPlayerId(), chat.GetMessage())
	}
}

// OnGameEnd is called once the game has ended
func (bot *exampleBot) OnGameEnd() {
	fmt.Println("gg")
}
