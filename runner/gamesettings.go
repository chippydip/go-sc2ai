package runner

import (
	"github.com/chippydip/go-sc2ai/api"
	"github.com/chippydip/go-sc2ai/client"
)

var processSettings = struct {
	realtime          bool
	processPath       string
	baseBuild         uint32
	dataVersion       string
	netAddress        string
	timeoutMS         int
	portStart         int
	extraCommandLines []string
	processInfo       []client.ProcessInfo
}{false, "", 0, "", "127.0.0.1", 120000, 8168, nil, nil}

var renderSettings = struct {
}{}

var featureLayerSettings = struct {
}{}

var interfaceOptions = &api.InterfaceOptions{Raw: true, Score: true}

var gameSettings = struct {
	mapName     string
	playerSetup []*api.PlayerSetup
	ports       client.Ports
}{}

var replaySettings = struct {
	dir    string
	files  []string
	player api.PlayerID
}{"", nil, 1}
