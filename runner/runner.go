package runner

import (
	"fmt"
	"log"
	"strings"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

var gamePort = 0
var startPort = 0
var ladderServer = ""
var opponentID = ""

var computerOpponent = false
var computerRace = api.Race_Terran
var computerDifficulty = api.Difficulty_Easy

func init() {
	// Ladder Flags
	flagInt("GamePort", &gamePort, "Port of client to connect to")
	flagInt("StartPort", &startPort, "Starting server port")
	flagStr("LadderServer", &ladderServer, "Ladder server address")
	flagStr("OpponentId", &opponentID, "Ladder ID of the opponent (for learning bots)")

	// Testing Flags
	flagBool("ComputerOpponent", &computerOpponent, "If we set up a computer opponent")
	flagVar("ComputerRace", (*raceFlag)(&computerRace), "Race of computer opponent")
	flagVar("ComputerDifficulty", (*difficultyFlag)(&computerDifficulty), "Difficulty of computer opponent")
}

// RunAgent ...
func RunAgent(agent client.PlayerSetup) {
	if !LoadSettings() {
		return
	}

	// fmt.Println(gamePort, startPort, ladderServer, computerOpponent, computerRace, computerDifficulty)
	// fmt.Println(processSettings, gameSettings)

	var numAgents = 1
	if computerOpponent && gamePort == 0 {
		SetParticipants(agent, client.NewComputer(computerRace, computerDifficulty))
	} else {
		numAgents = 2
		SetParticipants(agent)
	}

	if gamePort > 0 {
		log.Print("Connecting to port ", gamePort)
		Connect(gamePort)
		SetupPorts(numAgents, startPort, false)
		JoinGame()
		processSettings.timeoutMS = 10000
		log.Print(" Successfully joined game")
	} else {
		LaunchStarcraft()

		if len(replaySettings.files) > 0 {
			for _, file := range replaySettings.files {
				if StartReplay(file) {
					Run()
				}
			}
			return
		}

		StartGame(gameSettings.mapName)
	}

	Run()
}

type raceFlag api.Race

func (f *raceFlag) Set(value string) error {
	// Uppercase first character
	if len(value) > 0 {
		value = strings.ToUpper(value[:1]) + value[1:]
	}

	if v, ok := api.Race_value[value]; ok {
		*f = raceFlag(v)
		return nil
	}
	return fmt.Errorf("Unknown race: %v", value)
}

func (f *raceFlag) String() string {
	if v, ok := api.Difficulty_name[int32(*f)]; ok {
		return strings.ToLower(v)
	}
	return ""
}

type difficultyFlag api.Difficulty

func (f *difficultyFlag) Set(value string) error {
	if v, ok := api.Difficulty_value[value]; ok {
		*f = difficultyFlag(v)
		return nil
	}
	return fmt.Errorf("Unknown difficulty: %v", value)
}

func (f *difficultyFlag) String() string {
	if v, ok := api.Difficulty_name[int32(*f)]; ok {
		return v
	}
	return ""
}
