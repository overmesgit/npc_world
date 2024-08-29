package main
type GameMap struct {
	Tiles  [][]int
	Width  int
	Height int
}

func NewGameMap() *GameMap {
	width, height := 50, 50
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

