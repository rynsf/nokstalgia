package driver

var config = map[string]uint32{
	"maze":       0,
	"level":      0,
	"rules":      0,
	"challenges": 0,
}

var selectedGame string

func GetSelectedGame() string {
	return selectedGame
}

func SetSelectedGame(game string) {
	selectedGame = game
}

func GetConfig(id string) uint32 {
	return config[id]
}

func SetConfig(id string, val uint32) {
	config[id] = val
}
