package main

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

func main() {
	// Play a random map against a VeryHard difficulty computer
	runner.Set("map", runner.Random1v1Map())
	runner.Set("ComputerOpponent", "true")
	runner.Set("ComputerRace", "random")
	runner.Set("ComputerDifficulty", "VeryHard")

	// Create the agent and then start the game
	agent := client.AgentFunc(runAgent)
	runner.RunAgent(client.NewParticipant(api.Race_Terran, agent, "ProxyReapers"))
}
