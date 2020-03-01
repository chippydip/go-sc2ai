package botutil

import (
	"fmt"
	"log"
	"strings"

	"github.com/chippydip/go-sc2ai/client"
	"github.com/chippydip/go-sc2ai/runner"
)

// Bot ...
type Bot struct {
	client.AgentInfo
	GameLoop uint32

	*Player
	*UnitContext
	*Actions
	*Builder
}

// NewBot ...
func NewBot(info client.AgentInfo) *Bot {
	bot := &Bot{AgentInfo: info}

	bot.Player = NewPlayer(info)
	bot.Actions = NewActions(info)
	bot.UnitContext = NewUnitContext(info, bot)
	bot.Builder = NewBuilder(info, bot.Player, bot.UnitContext)

	update := func() {
		bot.GameLoop = bot.Observation().GetObservation().GetGameLoop()
		log.SetPrefix(fmt.Sprintf("[%v] ", bot.GameLoop))

		if bot.GameLoop == 224 {
			bot.checkVersion()
		}
	}
	update()
	bot.OnAfterStep(update)

	return bot
}

func (bot *Bot) checkVersion() {
	if c, ok := bot.AgentInfo.(*client.Client); !ok {
		log.Print("Skipping version check") // Should only happen when AgentInfo is mocked
	} else {
		// Check the game version, this should be less important but still worth reporting
		cVersion := formatVersion(c.GameVersion, c.BaseBuild)
		if c.BaseBuild == BaseBuild && c.GameVersion == GameVersion {
			bot.ChatTeam(fmt.Sprintf("(sc2) %v (thumbsup)", cVersion))
		} else {
			bot.ChatTeam(fmt.Sprintf("(sc2) %v (thumbsdown) (%v)", cVersion, formatVersion(GameVersion, BaseBuild)))
		}

		// It's critical that data versions match, however, or generated IDs may be wrong
		if c.DataBuild != DataBuild || c.DataVersion != DataVersion {
			bot.ChatTeam(fmt.Sprintf("(poo) (poo) (angry) (poo) (poo) %v:%v (scared) %v:%v",
				c.DataBuild, c.DataVersion, DataBuild, DataVersion))
		}
	}
}

func formatVersion(gameVersion string, baseBuild uint32) string {
	if strings.HasSuffix(gameVersion, fmt.Sprintf(".%v", baseBuild)) {
		return gameVersion
	}
	return fmt.Sprintf("%v(%v)", gameVersion, baseBuild)
}

// SetGameVersion sets the default base build and data version to the values
// last used to generate IDs.
func SetGameVersion() {
	runner.SetGameVersion(DataBuild, DataVersion)
}
