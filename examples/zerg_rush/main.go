package main

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

func main() {
	// Play a random map against a medium difficulty computer
	runner.SetMap(runner.Random1v1Map())
	runner.SetComputer(api.Race_Random, api.Difficulty_Medium, api.AIBuild_RandomBuild)

	// Create the agent and then start the game
	agent := client.AgentFunc(runAgent)
	runner.RunAgent(client.NewParticipant(api.Race_Zerg, agent, "ZergRush"))
}
