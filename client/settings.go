package client

import (
	"github.com/chippydip/go-sc2ai/api"
)

// ProcessInfo ...
type ProcessInfo struct {
	Path string
	PID  int
	Port int
}

// PlayerSetup ...
type PlayerSetup struct {
	PlayerType api.PlayerType
	Agent      Agent
	Race       api.Race
	Difficulty api.Difficulty
}

// NewParticipant ...
func NewParticipant(race api.Race, agent Agent) PlayerSetup {
	return PlayerSetup{
		api.PlayerType_Participant,
		agent,
		race,
		api.Difficulty_Easy,
	}
}

// NewComputer ...
func NewComputer(race api.Race, difficulty api.Difficulty) PlayerSetup {
	return PlayerSetup{
		api.PlayerType_Computer,
		nil,
		race,
		difficulty,
	}
}

// Ports ...
type Ports struct {
	ServerPorts *api.PortSet
	ClientPorts []*api.PortSet
	SharedPort  int32
}

func newPorts() Ports {
	return Ports{&api.PortSet{GamePort: -1, BasePort: -1}, []*api.PortSet{}, -1}
}

func (p Ports) isValid() bool {
	if p.SharedPort < 1 || !portSetIsValid(p.ServerPorts) || len(p.ClientPorts) < 1 {
		return false
	}

	for _, ps := range p.ClientPorts {
		if !portSetIsValid(ps) {
			return false
		}
	}

	return true
}

func portSetIsValid(ps *api.PortSet) bool {
	return ps != nil && ps.GamePort > 0 && ps.BasePort > 0
}
