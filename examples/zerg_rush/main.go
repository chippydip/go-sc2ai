package main

import (
	"github.com/chippydip/go-sc2ai/agent"
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

func main() {
	runner.Set("map", runner.Random1v1Map())
	runner.Set("ComputerOpponent", "true")
	runner.Set("ComputerRace", "random")
	runner.Set("ComputerDifficulty", "Medium")

	// Create the agent and then start the game
	agent := agent.AgentFunc(runAgent)
	runner.RunAgent(client.NewParticipant(api.Race_Zerg, agent, "ZergRush"))
}
