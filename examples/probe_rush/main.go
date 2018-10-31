package main

import (
	"math/rand"
	"time"

	"github.com/chippydip/go-sc2ai/api"
	abilityType "github.com/chippydip/go-sc2ai/api/ability"
	unitType "github.com/chippydip/go-sc2ai/api/unit" // using unit.Protoss_Nexus is not convenient, because everywhere you wish to iterate []*Units as unit
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

type probeRush struct {
	client.AgentInfo

	actions []*api.Action

	myStartLocation    api.Point2D
	enemyStartLocation api.Point2D
	observation        *api.Observation
	units              []*api.Unit
	myUnits            []*api.Unit
	enemyUnits         []*api.Unit
	mineralFields      []*api.Unit
	typeData           []*api.UnitTypeData

	homeMineral *api.Unit
	nexus       *api.Unit
	probes      []*api.Unit
	okTargets   []*api.Unit
	goodTargets []*api.Unit

	minerals uint32
	vespene  uint32
}

// User should check that he receives not nil
func closestUnit(pos api.Point2D, units []*api.Unit) *api.Unit {
	var closest *api.Unit
	for _, unit := range units {
		if closest == nil ||
			pos.Distance2(closest.Pos.ToPoint2D()) > pos.Distance2(unit.Pos.ToPoint2D()) {
			closest = unit
		}
	}
	return closest
}

// One of go flaws is that we need that function. Or is there better solution? Make map[api.Attribute]Bool?
func contains(s []api.Attribute, e api.Attribute) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (bot *probeRush) ChatSend(msg string) {
	bot.SendActions([]*api.Action{{
		ActionChat: &api.ActionChat{
			Channel: api.ActionChat_Broadcast,
			Message: msg,
		},
	}})
}

// OnGameStart is called once at the start of the game
func (bot *probeRush) OnGameStart(info client.AgentInfo) {
	bot.AgentInfo = info

	bot.enemyStartLocation = *bot.AgentInfo.GameInfo().GetStartRaw().GetStartLocations()[0]
	bot.units = bot.AgentInfo.Observation().GetObservation().GetRawData().GetUnits()
	// Find my nexus, it is on start position
	for _, unit := range bot.units {
		if unit.UnitType == unitType.Protoss_Nexus {
			bot.myStartLocation = unit.Pos.ToPoint2D()
			break
		}
	}
	/*fmt.Println(bot.myStartLocation)
	fmt.Println(bot.enemyStartLocation)*/

	bot.typeData = bot.AgentInfo.Data().GetUnits()

	// Send a friendly hello
	bot.ChatSend("(glhf)")
}

func (bot *probeRush) ParseObservation() {
	bot.observation = bot.AgentInfo.Observation().GetObservation()
	bot.units = bot.observation.GetRawData().GetUnits()
	bot.myUnits = nil
	bot.enemyUnits = nil
	bot.mineralFields = nil
	bot.nexus = nil
	bot.probes = nil
	bot.okTargets = nil
	bot.goodTargets = nil
	for _, unit := range bot.units {
		if unit.Alliance == api.Alliance_Self {
			bot.myUnits = append(bot.myUnits, unit)
			if unit.UnitType == unitType.Protoss_Nexus {
				bot.nexus = unit
			} else if unit.UnitType == unitType.Protoss_Probe {
				bot.probes = append(bot.probes, unit)
			}
		} else if unit.Alliance == api.Alliance_Enemy {
			bot.enemyUnits = append(bot.enemyUnits, unit)
			if !unit.IsFlying && unit.UnitType != unitType.Zerg_Larva && unit.UnitType != unitType.Zerg_Egg {
				bot.okTargets = append(bot.okTargets, unit)
				if !contains(bot.typeData[unit.UnitType].Attributes, api.Attribute_Structure) {
					bot.goodTargets = append(bot.goodTargets, unit)
				}
			}
		} else if unit.MineralContents > 0 {
			bot.mineralFields = append(bot.mineralFields, unit)
		}
	}

	bot.minerals = bot.observation.PlayerCommon.Minerals
	bot.vespene = bot.observation.PlayerCommon.Vespene
}

func (bot *probeRush) FirstStep() {
	bot.homeMineral = closestUnit(bot.myStartLocation, bot.mineralFields)
	// fmt.Println(bot.homeMineral)
}

func (bot *probeRush) UnitCommand(unit *api.Unit, ability abilityType.Ability) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []api.UnitTag{unit.Tag},
				}}}})
}

func (bot *probeRush) UnitCommandTargetTag(unit *api.Unit, ability abilityType.Ability, target api.UnitTag) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []api.UnitTag{unit.Tag},
					Target: &api.ActionRawUnitCommand_TargetUnitTag{
						TargetUnitTag: target,
					}}}}})
}

func (bot *probeRush) UnitCommandTargetPos(unit *api.Unit, ability abilityType.Ability, target api.Point2D) {
	bot.actions = append(bot.actions, &api.Action{
		ActionRaw: &api.ActionRaw{
			Action: &api.ActionRaw_UnitCommand{
				UnitCommand: &api.ActionRawUnitCommand{
					AbilityId: ability,
					UnitTags:  []api.UnitTag{unit.Tag},
					Target: &api.ActionRawUnitCommand_TargetWorldSpacePos{
						TargetWorldSpacePos: &target,
					}}}}})
}

func (bot *probeRush) UnitsCommandTargetPos(units []*api.Unit, ability abilityType.Ability, target api.Point2D) {
	// I hope, we can avoid this conversion in future
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

func (bot *probeRush) Strategy() {
	// Build probe if can
	if bot.nexus != nil && len(bot.nexus.Orders) == 0 && bot.minerals >= 50 {
		bot.UnitCommand(bot.nexus, abilityType.Train_Probe)
	}
	// Chronoboost self if building something
	if bot.nexus != nil && len(bot.nexus.Orders) > 0 && bot.nexus.Energy >= 50 {
		bot.UnitCommandTargetTag(bot.nexus, abilityType.Effect_ChronoBoostEnergyCost, bot.nexus.Tag)
	}
}

func (bot *probeRush) Tactics() {
	if len(bot.okTargets) == 0 && bot.probes != nil {
		// Attack enemy base position
		bot.UnitsCommandTargetPos(bot.probes, abilityType.Attack, bot.enemyStartLocation)
	} else {
		// To see battle better
		if len(bot.goodTargets) > 0 {
			time.Sleep(time.Millisecond * 20)
		}
		for _, probe := range bot.probes {
			if probe.Shield == 0 {
				bot.UnitCommandTargetTag(probe, abilityType.Harvest_Gather, bot.homeMineral.Tag)
			} else {
				target := new(api.Unit)
				if len(bot.goodTargets) > 0 {
					target = closestUnit(probe.Pos.ToPoint2D(), bot.goodTargets)
				} else {
					target = closestUnit(probe.Pos.ToPoint2D(), bot.okTargets)
				}
				bot.UnitCommandTargetPos(probe, abilityType.Attack, target.Pos.ToPoint2D())
			}
		}
	}
}

// OnStep is called each game step (every game update by defaul)
func (bot *probeRush) OnStep() {
	bot.ParseObservation()
	if bot.observation.GameLoop == 1 {
		bot.FirstStep()
	}

	bot.Strategy()
	bot.Tactics()

	if len(bot.actions) > 0 {
		bot.SendActions(bot.actions)
		bot.actions = nil
	}
}

// OnGameEnd is called once the game has ended
func (bot *probeRush) OnGameEnd() {
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
	runner.RunAgent(client.NewParticipant(api.Race_Protoss, &probeRush{}))
}
