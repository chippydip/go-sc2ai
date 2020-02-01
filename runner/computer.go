package runner

import (
	"fmt"
	"strings"

	"github.com/chippydip/go-sc2ai/api"
)

var (
	computerOpponent   = false
	computerRace       = api.Race_Terran
	computerDifficulty = api.Difficulty_Easy
	computerBuild      = api.AIBuild_RandomBuild
)

func init() {
	// Testing Flags
	flagBool("ComputerOpponent", &computerOpponent, "If we set up a computer opponent")
	flagVar("ComputerRace", (*raceFlag)(&computerRace), "Race of computer opponent")
	flagVar("ComputerDifficulty", (*difficultyFlag)(&computerDifficulty), "Difficulty of computer opponent")
	flagVar("ComputerBuild", (*buildFlag)(&computerBuild), "Build of computer opponent")
}

// SetComputer sets the default computer opponent flags (can still be overridden on the command line).
func SetComputer(race api.Race, difficulty api.Difficulty, build api.AIBuild) {
	Set("ComputerOpponent", "true")
	Set("ComputerRace", api.Race_name[int32(race)])
	Set("ComputerDifficulty", api.Difficulty_name[int32(difficulty)])
	Set("ComputerBuild", api.AIBuild_name[int32(build)])
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

type buildFlag api.AIBuild

func (f *buildFlag) Set(value string) error {
	if v, ok := api.AIBuild_value[value]; ok {
		*f = buildFlag(v)
		return nil
	}
	return fmt.Errorf("Unknown build: %v", value)
}

func (f *buildFlag) String() string {
	if v, ok := api.AIBuild_name[int32(*f)]; ok {
		return v
	}
	return ""
}
