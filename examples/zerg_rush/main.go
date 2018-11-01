package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/enums/ability"
	"github.com/chippydip/go-sc2ai/enums/zerg"
	"github.com/chippydip/go-sc2ai/runner"
)

type zergRush struct {
	client.AgentInfo

	actions []*api.Action

	myStartLocation    api.Point2D
	enemyStartLocation api.Point2D
	observation        *api.Observation
	units              map[api.UnitTypeID][]*api.Unit
	neutralUnits       map[api.UnitTypeID][]*api.Unit
	enemyUnits         map[api.UnitTypeID][]*api.Unit

	typeData       []*api.UnitTypeData // actually, a map[unitType.Unit]
	unitAttributes map[api.UnitTypeID]map[api.Attribute]bool

	okTargets   []*api.Unit
	goodTargets []*api.Unit

	minerals uint32
	vespene  uint32
	foodCap  int
	foodUsed int
	foodLeft int
}

func New() *zergRush {
	bot := new(zergRush)
	bot.unitAttributes = map[api.UnitTypeID]map[api.Attribute]bool{}
	return bot
}

func (bot *zergRush) parseUnits() {
	bot.units = map[api.UnitTypeID][]*api.Unit{}
	bot.neutralUnits = map[api.UnitTypeID][]*api.Unit{}
	bot.enemyUnits = map[api.UnitTypeID][]*api.Unit{}
	for _, unit := range bot.AgentInfo.Observation().GetObservation().GetRawData().GetUnits() {
		switch unit.Alliance {
		case api.Alliance_Self:
			bot.units[unit.UnitType] = append(bot.units[unit.UnitType], unit)
		case api.Alliance_Enemy:
			bot.enemyUnits[unit.UnitType] = append(bot.enemyUnits[unit.UnitType], unit)
		case api.Alliance_Neutral:
			bot.neutralUnits[unit.UnitType] = append(bot.neutralUnits[unit.UnitType], unit)
		default:
			fmt.Println("Not supported alliance: ", unit)
		}
	}
}

// OnGameStart is called once at the start of the game
func (bot *zergRush) OnGameStart(info client.AgentInfo) {
	bot.AgentInfo = info
	bot.parseUnits()
	bot.typeData = bot.AgentInfo.Data().GetUnits()
	for _, utd := range bot.typeData {
		for _, attribute := range utd.Attributes {
			bot.unitAttributes[utd.UnitId] = map[api.Attribute]bool{}
			bot.unitAttributes[utd.UnitId][attribute] = true
		}
	}

	// My hatchery is on start position
	bot.myStartLocation = bot.units[zerg.Hatchery][0].Pos.ToPoint2D()
	bot.enemyStartLocation = *bot.AgentInfo.GameInfo().GetStartRaw().GetStartLocations()[0]

	// Send a friendly hello
	bot.ChatSend("(glhf)")
}

func (bot *zergRush) ParseObservation() {
	bot.observation = bot.AgentInfo.Observation().GetObservation()
	bot.parseUnits()

	bot.okTargets = nil
	bot.goodTargets = nil
	for _, units := range bot.enemyUnits {
		for _, unit := range units {
			if !unit.IsFlying && unit.UnitType != zerg.Larva && unit.UnitType != zerg.Egg {
				bot.okTargets = append(bot.okTargets, unit)
				if !bot.unitAttributes[unit.UnitType][api.Attribute_Structure] {
					bot.goodTargets = append(bot.goodTargets, unit)
				}
			}
		}
	}

	bot.minerals = bot.observation.PlayerCommon.Minerals
	bot.vespene = bot.observation.PlayerCommon.Vespene
	bot.foodCap = int(bot.observation.PlayerCommon.FoodCap)
	bot.foodUsed = int(bot.observation.PlayerCommon.FoodUsed)
	bot.foodLeft = bot.foodCap - bot.foodUsed
}

func (bot *zergRush) AlreadyTraining(abilityID api.AbilityID) int {
	count := 0
	units := bot.units[zerg.Egg]
	if abilityID == ability.Train_Queen {
		units = bot.units[zerg.Hatchery]
	}
	for _, unit := range units {
		if len(unit.Orders) > 0 && unit.Orders[0].AbilityId == abilityID {
			count++
		}
	}
	return count
}

func (bot *zergRush) Strategy() {
	// Wait until we have enough minerals to build spawning pool
	if bot.minerals >= 200 && len(bot.units[zerg.SpawningPool]) == 0 &&
		len(bot.units[zerg.Drone]) > 0 {
		builder := bot.units[zerg.Drone][0]
		pos := bot.myStartLocation.Offset(bot.enemyStartLocation, 5)
		bot.unitCommandTargetPos(builder, ability.Build_SpawningPool, pos)
		return
	}
	// We are building spawning pool (or it is ready)
	if len(bot.units[zerg.SpawningPool]) > 0 {
		// Build drones
		if len(bot.units[zerg.Drone]) < 14 && len(bot.units[zerg.Larva]) > 0 && bot.minerals >= 50 {
			bot.unitCommand(bot.units[zerg.Larva][0], ability.Train_Drone)
			return
		}
		// Build overlords
		if bot.foodLeft < 2 && len(bot.units[zerg.Larva]) > 0 && bot.minerals >= 100 &&
			bot.AlreadyTraining(ability.Train_Overlord) == 0 {
			bot.unitCommand(bot.units[zerg.Larva][0], ability.Train_Overlord)
			return
		}
		// If pool is ready
		if bot.units[zerg.SpawningPool][0].BuildProgress == 1 {
			// Build zerglings
			if len(bot.units[zerg.Larva]) > 0 && bot.minerals >= 50 {
				bot.unitCommand(bot.units[zerg.Larva][0], ability.Train_Zergling)
				return
			}
			if len(bot.units[zerg.Queen]) == 0 && bot.minerals >= 150 &&
				bot.AlreadyTraining(ability.Train_Queen) == 0 && len(bot.units[zerg.Hatchery]) > 0 {
				bot.unitCommand(bot.units[zerg.Hatchery][0], ability.Train_Queen)
				return
			}
		}
	}
}

func (bot *zergRush) Tactics() {
	lings := bot.units[zerg.Zergling]
	if len(lings) >= 6 {
		if len(bot.okTargets) == 0 {
			bot.unitsCommandTargetPos(lings, ability.Attack, bot.enemyStartLocation)
		} else {
			// To see battle better
			if len(bot.goodTargets) > 0 {
				time.Sleep(time.Millisecond * 20)
			}
			for _, ling := range lings {
				if len(bot.goodTargets) > 0 {
					target := ClosestUnit(ling.Pos.ToPoint2D(), bot.goodTargets)
					if ling.Pos.ToPoint2D().Distance2(target.Pos.ToPoint2D()) > 4*4 {
						// If target is far, attack it as unit, ling will run ignoring everything else
						bot.unitCommandTargetTag(ling, ability.Attack, target.Tag)
					} else {
						// Attack as position, ling will choose best target around
						bot.unitCommandTargetPos(ling, ability.Attack, target.Pos.ToPoint2D())
					}
				} else {
					target := ClosestUnit(ling.Pos.ToPoint2D(), bot.okTargets)
					bot.unitCommandTargetPos(ling, ability.Attack, target.Pos.ToPoint2D())
				}
			}
		}
	}

	queens := bot.units[zerg.Queen]
	// If queen can inject, do it
	if len(queens) > 0 && queens[0].Energy >= 25 && len(bot.units[zerg.Hatchery]) > 0 {
		bot.unitCommandTargetTag(queens[0], ability.Effect_InjectLarva, bot.units[zerg.Hatchery][0].Tag)
	}
}

// OnStep is called each game step (every game update by defaul)
func (bot *zergRush) OnStep() {
	bot.ParseObservation()

	bot.Strategy()
	bot.Tactics()

	if len(bot.actions) > 0 {
		bot.SendActions(bot.actions)
		bot.actions = nil
	}
}

// OnGameEnd is called once the game has ended
func (bot *zergRush) OnGameEnd() {
	bot.ChatSend("(gg)")
}

func main() {
	maps := []string{"AcidPlantLE", "BlueshiftLE", "CeruleanFallLE", "DreamcatcherLE",
		"FractureLE", "LostAndFoundLE", "ParaSiteLE"}

	rand.Seed(time.Now().UnixNano())
	runner.Set("map", maps[rand.Intn(len(maps))]+".SC2Map")
	runner.Set("ComputerOpponent", "true")
	runner.Set("ComputerRace", "random")
	runner.Set("ComputerDifficulty", "Medium")

	// Create the agent and then start the game
	runner.RunAgent(client.NewParticipant(api.Race_Zerg, New()))
}
