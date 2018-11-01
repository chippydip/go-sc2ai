package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	// Checkout a temp copy of the API files
	dir, err := ioutil.TempDir("", "s2client-proto")
	check(err)
	defer os.RemoveAll(dir)

	// // Preserve the test directory to look at
	// defer func() {
	// 	if err := recover(); err != nil {
	// 		fmt.Println(err)
	// 	}

	// 	fmt.Print("Press 'Enter' to continue...")
	// 	bufio.NewReader(os.Stdin).ReadBytes('\n')
	// }()

	_, err = git.PlainClone(dir, false, &git.CloneOptions{
		URL:      "https://github.com/Blizzard/s2client-proto",
		Progress: os.Stdout,
	})
	check(err)

	// Get all the .proto files
	protoDir := filepath.Join(dir, "s2clientprotocol")
	files, err := ioutil.ReadDir(protoDir)
	check(err)

	protocArgs := []string{
		"-I=" + os.Getenv("GOPATH") + "/src/github.com/gogo/protobuf/gogoproto",
		"-I=" + os.Getenv("GOPATH") + "/src/github.com/gogo/protobuf/protobuf",
		"--proto_path=" + protoDir,
		"--gogofaster_out=api",
	}
	for _, file := range files {
		if filepath.Ext(file.Name()) != ".proto" {
			continue
		}
		path := filepath.Join(protoDir, file.Name())

		// Upgrade the file to proto3 and fix the package name
		writeLines(path, upgradeProto(path))

		// Add the file to the list of command line args for protoc
		protocArgs = append(protocArgs, path)
	}

	// Make sure we mapped all the expected types
	if len(typeMap) != 0 {
		fmt.Println("Not all types were mapped, missing:")
		for key := range typeMap {
			fmt.Println(key)
		}
	}

	// Generate go code from the .proto files
	out, err := exec.Command("protoc", protocArgs...).CombinedOutput()
	fmt.Print(string(out))
	check(err)
}

// Thing we want to use twice
const (
	importPrefix   = "import \"s2clientprotocol/"
	optionalPrefix = "optional "
	enumPrefix     = "enum "
	messagePrefix  = "message "
)

func upgradeProto(path string) []string {
	file, err := os.Open(path)
	check(err)
	defer file.Close()

	propPath := []string{}
	var lines []string

	// Read line by line, making modifications as needed
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		// Get the line and trim comments and whitespace to make matching easier
		line := scanner.Text()
		if comment := strings.Index(line, "//"); comment > 0 {
			line = line[:comment]
		}
		line = strings.TrimSpace(line)

		switch {
		// Upgrade to proto3 and set the go package name
		case line == "syntax = \"proto2\";":
			lines = append(lines, "syntax = \"proto3\";", "option go_package = \"api\";\nimport \"gogo.proto\";")

		// Remove subdirectory of the import so the output path isn't nested
		case strings.HasPrefix(line, importPrefix):
			lines = append(lines, "import \""+line[len(importPrefix):])

		// Remove "optional" prefixes (they are implicit in proto3)
		case strings.HasPrefix(line, optionalPrefix):
			lines = append(lines, mapTypes(propPath, line[len(optionalPrefix):]))

		// Track where we are in the path
		case strings.HasSuffix(line, " {"):
			id := strings.Split(line, " ")[1] // "<type> Identifier {"
			propPath = append(propPath, id)

			lines = append(lines, line)

			// Enums must have a zero value in proto3 (and unfortunately they must be unique due to C++ scoping rules)
			if strings.HasPrefix(line, enumPrefix) && line != "enum Race {" {
				lines = append(lines, line[len(enumPrefix):len(line)-2]+"_nil = 0 [(gogoproto.enumvalue_customname) = \"nil\"];")
			}

		// Pop the last path element
		case line == "}":
			propPath = propPath[:len(propPath)-1]
			lines = append(lines, line)

		// Everything else just gets copied to the output
		default:
			lines = append(lines, mapTypes(propPath, line))
		}
	}

	return lines
}

func mapTypes(path []string, line string) string {
	parts := strings.Split(line, " ")
	if len(parts) < 4 {
		return line // need at least "<type> <name> = <num>;"
	}

	key := strings.Join(path, ".") + "." + parts[len(parts)-3]
	if value, ok := typeMap[key]; ok {
		// Add the casttype option
		opt := fmt.Sprintf("[(gogoproto.casttype) = \"%v\"];", value)

		last := parts[len(parts)-1]
		parts = append(parts[:len(parts)-1], last[:len(last)-1], opt)
		delete(typeMap, key) // track which ones have been processed
		return strings.Join(parts, " ")
	}

	return line
}

func writeLines(path string, lines []string) {
	file, err := os.Create(path)
	check(err)
	defer file.Close()

	writer := bufio.NewWriter(file)
	for _, line := range lines {
		fmt.Fprintln(writer, line)
	}
	check(writer.Flush())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

// Make the API more type-safe
var typeMap = map[string]string{
	// data.proto
	"AbilityData.ability_id":           "AbilityID",
	"AbilityData.remaps_to_ability_id": "AbilityID",
	"UnitTypeData.unit_id":             "UnitTypeID",
	"UnitTypeData.ability_id":          "AbilityID",
	"UnitTypeData.tech_alias":          "UnitTypeID",
	"UnitTypeData.unit_alias":          "UnitTypeID",
	"UnitTypeData.tech_requirement":    "UnitTypeID",
	"UpgradeData.upgrade_id":           "UpgradeID",
	"UpgradeData.ability_id":           "AbilityID",
	"BuffData.buff_id":                 "BuffID",
	"EffectData.effect_id":             "EffectID",
	// debug.proto
	"DebugCreateUnit.unit_type":  "UnitTypeID",
	"DebugCreateUnit.owner":      "PlayerID",
	"DebugKillUnit.tag":          "UnitTag",
	"DebugSetUnitValue.unit_tag": "UnitTag",
	// error.proto
	// query.proto
	"RequestQueryPathing.start.unit_tag":             "UnitTag",
	"RequestQueryAvailableAbilities.unit_tag":        "UnitTag",
	"ResponseQueryAvailableAbilities.unit_tag":       "UnitTag",
	"ResponseQueryAvailableAbilities.unit_type_id":   "UnitTypeID",
	"RequestQueryBuildingPlacement.ability_id":       "AbilityID",
	"RequestQueryBuildingPlacement.placing_unit_tag": "UnitTag",
	// raw.proto
	"PowerSource.tag":                             "UnitTag",
	"PlayerRaw.upgrade_ids":                       "UpgradeID",
	"UnitOrder.ability_id":                        "AbilityID",
	"UnitOrder.target.target_unit_tag":            "UnitTag",
	"PassengerUnit.tag":                           "UnitTag",
	"PassengerUnit.unit_type":                     "UnitTypeID",
	"Unit.tag":                                    "UnitTag",
	"Unit.unit_type":                              "UnitTypeID",
	"Unit.owner":                                  "PlayerID",
	"Unit.add_on_tag":                             "UnitTag",
	"Unit.buff_ids":                               "BuffID",
	"Unit.engaged_target_tag":                     "UnitTag",
	"Event.dead_units":                            "UnitTag",
	"Effect.effect_id":                            "EffectID",
	"ActionRawUnitCommand.ability_id":             "AbilityID",
	"ActionRawUnitCommand.target.target_unit_tag": "UnitTag",
	"ActionRawUnitCommand.unit_tags":              "UnitTag",
	"ActionRawToggleAutocast.ability_id":          "AbilityID",
	"ActionRawToggleAutocast.unit_tags":           "UnitTag",
	// sc2api.proto
	"RequestJoinGame.participation.observed_player_id": "PlayerID",
	"ResponseJoinGame.player_id":                       "PlayerID",
	"RequestStartReplay.observed_player_id":            "PlayerID",
	"ChatReceived.player_id":                           "PlayerID",
	"PlayerInfo.player_id":                             "PlayerID",
	"PlayerCommon.player_id":                           "PlayerID",
	"ActionError.unit_tag":                             "UnitTag",
	"ActionError.ability_id":                           "AbilityID",
	"ActionObserverPlayerPerspective.player_id":        "PlayerID",
	"ActionObserverCameraFollowPlayer.player_id":       "PlayerID",
	"ActionObserverCameraFollowUnits.unit_tags":        "UnitTag",
	"PlayerResult.player_id":                           "PlayerID",
	// score.proto
	// spatial.proto
	// ui.proto
	"ControlGroup.leader_unit_type":   "UnitTypeID",
	"UnitInfo.unit_type":              "UnitTypeID",
	"UnitInfo.player_relative":        "PlayerID", // TODO: is this correct?
	"BuildItem.ability_id":            "AbilityID",
	"ActionToggleAutocast.ability_id": "AbilityID",
}

// TODO: spatial.proto?
