package main

// map.go
type GameMap struct {
	Tiles [][]int
	// Add other map properties
}

func NewGameMap() *GameMap {
	// Initialize map
	return &GameMap{}
}