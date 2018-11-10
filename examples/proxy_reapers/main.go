package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"

	"github.com/chippydip/go-sc2ai/agent"
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
	"github.com/chippydip/go-sc2ai/search"
)

type proxyReapers struct {
	agent agent.Agent
	info  client.AgentInfo

	actions            []*api.Action
	myStartLocation    api.Point2D
	enemyStartLocation api.Point2D
	baseLocations      []api.Point2D
	units              UnitsByTypes
	mineralFields      UnitsByTypes
	vespeneGeysers     UnitsByTypes
	neutralUnits       UnitsByTypes
	enemyUnits         UnitsByTypes
	orders             map[api.AbilityID]int

	minerals int
	vespene  int
	foodCap  int
	foodUsed int
	foodLeft int

	positionsForSupplies []api.Point2D
	positionsForBarracks []api.Point2D

	okTargets   Units
	goodTargets Units
	builder1    api.UnitTag
	builder2    api.UnitTag
	retreat     map[api.UnitTag]bool
}

func (bot *proxyReapers) parseObservation() {
	bot.minerals = int(bot.info.Observation().Observation.PlayerCommon.Minerals)
	bot.vespene = int(bot.info.Observation().Observation.PlayerCommon.Vespene)
	bot.foodCap = int(bot.info.Observation().Observation.PlayerCommon.FoodCap)
	bot.foodUsed = int(bot.info.Observation().Observation.PlayerCommon.FoodUsed)
	bot.foodLeft = bot.foodCap - bot.foodUsed
}

func (bot *proxyReapers) parseUnits() {
	bot.units = UnitsByTypes{}
	bot.mineralFields = UnitsByTypes{}
	bot.vespeneGeysers = UnitsByTypes{}
	bot.neutralUnits = UnitsByTypes{}
	bot.enemyUnits = UnitsByTypes{}
	for _, unit := range bot.info.Observation().Observation.RawData.Units {
		var units *UnitsByTypes
		switch unit.Alliance {
		case api.Alliance_Self:
			units = &bot.units
		case api.Alliance_Enemy:
			units = &bot.enemyUnits
		case api.Alliance_Neutral:
			if unit.MineralContents > 0 {
				units = &bot.mineralFields
			} else if unit.VespeneContents > 0 {
				units = &bot.vespeneGeysers
			} else {
				units = &bot.neutralUnits
			}
		default:
			fmt.Fprintln(os.Stderr, "Not supported alliance: ", unit)
			continue
		}
		units.AddFromApi(unit.UnitType, unit)
	}
}

func (bot *proxyReapers) parseOrders() {
	bot.orders = map[api.AbilityID]int{}
	for _, unitTypes := range bot.units {
		for _, unit := range unitTypes {
			for _, order := range unit.Orders {
				bot.orders[order.AbilityId]++
			}
		}
	}
}

// OnGameStart is called once at the start of the game
func (bot *proxyReapers) OnGameStart() {
	defer recoverPanic()

	bot.info = bot.agent.Info()
	InitUnits(bot.info.Data().Units)
	bot.parseUnits()
	bot.initLocations()
	for _, uc := range search.CalculateExpansionLocations(bot.info, false) {
		bot.baseLocations = append(bot.baseLocations, uc.Center())
	}
	bot.findBuildingsPositions()
	bot.retreat = map[api.UnitTag]bool{}

	// Send a friendly hello
	bot.chatSend("(glhf)")
}

// OnStep is called each game step (every game update by defaul)
func (bot *proxyReapers) OnStep() {
	defer recoverPanic()

	bot.info = bot.agent.Info()
	bot.parseObservation()
	bot.parseUnits()
	bot.parseOrders()

	bot.strategy()
	bot.tactics()

	if len(bot.actions) > 0 {
		bot.info.SendActions(bot.actions)
		bot.actions = nil
	}
}

// OnGameEnd is called once the game has ended
func (bot *proxyReapers) OnGameEnd() {
	bot.chatSend("(gg)")
}

func runAgent(a agent.Agent) {
	bot := proxyReapers{agent: a}
	bot.OnGameStart()

	for bot.info.IsInGame() {
		bot.OnStep()
		bot.info.Update(1)
	}
	bot.OnGameEnd()
}

func main() {
	maps := []string{"AcidPlantLE", "BlueshiftLE", "CeruleanFallLE", "DreamcatcherLE",
		"FractureLE", "LostAndFoundLE", "ParaSiteLE"}

	rand.Seed(time.Now().UnixNano())
	runner.Set("map", maps[rand.Intn(len(maps))]+".SC2Map")
	// runner.Set("map", "ParaSiteLE.SC2Map")
	runner.Set("ComputerOpponent", "true")
	runner.Set("ComputerRace", "random")
	runner.Set("ComputerDifficulty", "VeryHard")

	// Create the agent and then start the game
	runner.RunAgent(client.NewParticipant(api.Race_Terran, agent.AgentFunc(runAgent)))
}
