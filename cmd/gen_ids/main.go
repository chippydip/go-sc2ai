package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

func main() {
	runner.SetMap(runner.Random1v1Map())

	agent := client.AgentFunc(generate)
	runner.RunAgent(client.NewParticipant(api.Race_Random, agent, "NilBot"))
}

func generate(info client.AgentInfo) {
	dumpAbilities(info.Data().GetAbilities(), info.Data().GetUnits())
	dumpBuffs(info.Data().GetBuffs())
	dumpEffects(info.Data().GetEffects())
	dumpUnits(info.Data().GetUnits())
	dumpUpgrades(info.Data().GetUpgrades())
}

func dumpAbilities(abilities []*api.AbilityData, units []*api.UnitTypeData) {
	// Detect base abilities of things with assigned hotkeys
	remaps := map[api.AbilityID]bool{}
	for _, ability := range abilities {
		if ability.GetAvailable() && ability.ButtonName != "" {
			if ability.RemapsToAbilityId != 0 && ability.Hotkey != "" {
				remaps[ability.RemapsToAbilityId] = true
			}
		}
	}

	// Find values to export and detect duplicate names
	byName := map[string]int{}
	for _, ability := range abilities {
		if ability.GetAvailable() && ability.ButtonName != "" {
			if ability.Hotkey != "" || remaps[ability.AbilityId] {
				byName[ability.FriendlyName] = byName[ability.FriendlyName] + 1
			}
		}
	}

	// Generate the values
	names := []string{}
	values := map[string]uint32{}
	namesMap := map[uint32]string{}
	idRemaps := map[uint32]uint32{}
	for _, ability := range abilities {
		n := byName[ability.FriendlyName]
		if n == 0 {
			continue
		}

		if ability.GetAvailable() && ability.ButtonName != "" {
			if ability.Hotkey != "" || remaps[ability.AbilityId] {
				name := ability.FriendlyName
				if n > 1 {
					name = fmt.Sprintf("%v %v", name, uint32(ability.AbilityId))
				}
				name = makeID(name)

				names = append(names, name)
				values[name] = uint32(ability.AbilityId)
				namesMap[uint32(ability.AbilityId)] = name

				if ability.RemapsToAbilityId != 0 {
					idRemaps[uint32(ability.AbilityId)] = uint32(ability.RemapsToAbilityId)
				}

			}
		}
	}
	sort.Strings(names)

	values["Invalid"] = 0
	values["Smart"] = 1
	writeEnum("ability", "AbilityID", append([]string{"Invalid", "Smart"}, names...), values)

	// Map to built units
	mapAbilityToProducedUnit(names, units)

	// Remap abilities to generic versions
	remapAbilities(idRemaps, namesMap)
}

func dumpBuffs(buffs []*api.BuffData) {
	names := []string{}
	values := map[string]uint32{}
	for _, buff := range buffs {
		if name := makeID(buff.GetName()); name != "" {
			names = append(names, name)
			values[name] = uint32(buff.BuffId)
		}
	}
	//sort.Strings(names)

	values["Invalid"] = 0
	writeEnum("buff", "BuffID", append([]string{"Invalid"}, names...), values)
}

func dumpEffects(effects []*api.EffectData) {
	names := []string{}
	values := map[string]uint32{}
	for _, effect := range effects {
		if name := makeID(effect.GetFriendlyName()); name != "" {
			names = append(names, name)
			values[name] = uint32(effect.EffectId)
		}
	}
	//sort.Strings(names)

	values["Invalid"] = 0
	writeEnum("effect", "EffectID", append([]string{"Invalid"}, names...), values)
}

func dumpUnits(units []*api.UnitTypeData) {
	names := []string{}
	values := map[string]uint32{}
	namesByRace := map[string][]string{}
	valuesByRace := map[string]map[string]uint32{}
	for _, unit := range units {
		if unit.GetAvailable() && unit.Name != "" {
			race := unit.Race.String()
			if race == "NoRace" {
				race = "Neutral"
			}
			name := makeID(race + "_" + unit.Name)

			names = append(names, name)
			values[name] = uint32(unit.UnitId)

			namesByRace[race] = append(namesByRace[race], unit.Name)
			if valuesByRace[race] == nil {
				valuesByRace[race] = make(map[string]uint32)
			}
			valuesByRace[race][unit.Name] = uint32(unit.UnitId)
		}
	}
	sort.Strings(names)

	values["Invalid"] = 0
	writeEnum("unit", "UnitTypeID", append([]string{"Invalid"}, names...), values)

	for race, names := range namesByRace {
		sort.Strings(names)

		writeEnum(strings.ToLower(race), "UnitTypeID", names, valuesByRace[race])
	}
}

func dumpUpgrades(upgrades []*api.UpgradeData) {
	names := []string{}
	values := map[string]uint32{}
	for _, upgrade := range upgrades {
		if name := makeID(upgrade.GetName()); name != "" {
			names = append(names, name)
			values[name] = uint32(upgrade.UpgradeId)
		}
	}
	//sort.Strings(names)

	values["Invalid"] = 0
	writeEnum("upgrade", "UpgradeID", append([]string{"Invalid"}, names...), values)
}

func makeID(id string) string {
	id = strings.Replace(id, " ", "_", -1)
	for _, c := range id {
		if !unicode.IsLetter(c) {
			return "A_" + id
		}
		if unicode.IsLower(c) {
			return string(unicode.ToUpper(c)) + id[1:]
		}
		break
	}
	return id
}

func writeEnum(name string, apiType string, names []string, values map[string]uint32) {
	pkgName := strings.ToLower(name)
	fmtString := "\t%-*v api." + apiType + " = %v\n"

	maxLen, maxVal := 0, uint32(0)
	for _, name := range names {
		if len(name) > maxLen {
			maxLen = len(name)
		}
		if val := values[name]; val > maxVal {
			maxVal = val
		}
	}

	check(os.MkdirAll("enums/"+pkgName, 0777))
	file, err := os.Create("enums/" + pkgName + "/enum.go")
	check(err)
	defer file.Close()

	w := bufio.NewWriter(file)

	fmt.Fprint(w, "// Code generated by gen_ids. DO NOT EDIT.\npackage "+
		pkgName+"\n\nimport \"github.com/chippydip/go-sc2ai/api\"\n\nconst (\n")

	for _, name := range names {
		fmt.Fprintf(w, fmtString, maxLen, name, values[name])
	}
	fmt.Fprint(w, ")\n")
	check(w.Flush())

	if !strings.HasPrefix(strings.ToLower(apiType), name) {
		return
	}

	// String() function
	fmtString2 := "\t%-*v \"%v\",\n"
	file2, err := os.Create("enums/" + pkgName + "/strings.go")
	check(err)
	defer file.Close()

	w2 := bufio.NewWriter(file2)

	fmt.Fprint(w2, "// Code generated by gen_ids. DO NOT EDIT.\npackage "+pkgName+
		"\n\nimport \"github.com/chippydip/go-sc2ai/api\"\n\n"+
		"func String(e api."+apiType+") string {\n\treturn strings[uint32(e)]\n}\n\nvar strings = map[uint32]string{\n")

	maxDigits := int(math.Ceil(math.Log10(float64(maxVal)))) + 1
	for _, name := range names {
		fmt.Fprintf(w2, fmtString2, maxDigits, strconv.Itoa(int(values[name]))+":", name)
	}
	fmt.Fprint(w2, "}\n")
	check(w2.Flush())
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}
