package main

// world.go
type World struct {
	characters []Character
	monsters   []Monster
	gameMap        *GameMap
}

func NewWorld() *World {
	return &World{
		characters: make([]Character, 0),
		monsters:   make([]Monster, 0),
		gameMap:        NewGameMap(),
	}
}

func (w *World) Update() {
	for i := range w.characters {
		w.characters[i].Update()
	}
	for i := range w.monsters {
		w.monsters[i].Update()
	}
}
