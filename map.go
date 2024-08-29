package main

const TileSize = 32 // Define the size of each tile in pixels

type GameMap struct {
	Tiles  [][]int
	Width  int // in tiles
	Height int // in tiles
}

func NewGameMap() *GameMap {
	width, height := 50, 50 // 50x50 tiles
	tiles := make([][]int, height)
	for i := range tiles {
		tiles[i] = make([]int, width)
	}
	return &GameMap{
		Tiles:  tiles,
		Width:  width,
		Height: height,
	}
}

