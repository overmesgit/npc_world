package main

import (
	"math/rand"
)

const (
	TileSize      = 32
	AreaSize      = 16 // Size of each area in tiles
	MountainWidth = 2  // Width of mountain ranges
	PassageSize   = 4  // Size of the passages at mountain intersections
)

type TileType int

const (
	TileGrass TileType = iota
	TileMountain
)

type GameMap struct {
	Tiles  [][]TileType
	Width  int // in tiles
	Height int // in tiles
}

func NewGameMap() *GameMap {
	width, height := AreaSize*3+MountainWidth*2, AreaSize*3+MountainWidth*2
	tiles := make([][]TileType, height)
	for i := range tiles {
		tiles[i] = make([]TileType, width)
	}

	// Generate the map
	for y := 0; y < height; y++ {
		for x := 0; x < width; x++ {
			// Check if we're on a mountain range
			if isMountainTile(x, y) {
				tiles[y][x] = TileMountain
			} else {
				tiles[y][x] = TileGrass
			}
		}
	}

	return &GameMap{
		Tiles:  tiles,
		Width:  width,
		Height: height,
	}
}

func isMountainTile(x, y int) bool {
	// Horizontal mountain ranges
	isHorizontalMountain := (y >= AreaSize && y < AreaSize+MountainWidth) ||
		(y >= AreaSize*2+MountainWidth && y < AreaSize*2+MountainWidth*2)

	// Vertical mountain ranges
	isVerticalMountain := (x >= AreaSize && x < AreaSize+MountainWidth) ||
		(x >= AreaSize*2+MountainWidth && x < AreaSize*2+MountainWidth*2)

	// Check if we're in one of the open areas at intersections
	isOpenArea := isInOpenArea(x, y)

	return (isHorizontalMountain || isVerticalMountain) && !isOpenArea
}

func isInOpenArea(x, y int) bool {
	// Top-left open area
	if x >= AreaSize-PassageSize/2 && x < AreaSize+MountainWidth+PassageSize/2 &&
		y >= AreaSize-PassageSize/2 && y < AreaSize+MountainWidth+PassageSize/2 {
		return true
	}

	// Top-right open area
	if x >= AreaSize*2+MountainWidth-PassageSize/2 && x < AreaSize*2+MountainWidth*2+PassageSize/2 &&
		y >= AreaSize-PassageSize/2 && y < AreaSize+MountainWidth+PassageSize/2 {
		return true
	}

	// Bottom-left open area
	if x >= AreaSize-PassageSize/2 && x < AreaSize+MountainWidth+PassageSize/2 &&
		y >= AreaSize*2+MountainWidth-PassageSize/2 && y < AreaSize*2+MountainWidth*2+PassageSize/2 {
		return true
	}

	// Bottom-right open area
	if x >= AreaSize*2+MountainWidth-PassageSize/2 && x < AreaSize*2+MountainWidth*2+PassageSize/2 &&
		y >= AreaSize*2+MountainWidth-PassageSize/2 && y < AreaSize*2+MountainWidth*2+PassageSize/2 {
		return true
	}

	return false
}

func (m *GameMap) IsTileWalkable(x, y int) bool {
	if x < 0 || x >= m.Width || y < 0 || y >= m.Height {
		return false
	}
	return m.Tiles[y][x] != TileMountain
}

// Add some random features to each area
func (m *GameMap) AddRandomFeatures() {
	for ay := 0; ay < 3; ay++ {
		for ax := 0; ax < 3; ax++ {
			m.addFeaturesToArea(ax, ay)
		}
	}
}

func (m *GameMap) addFeaturesToArea(areaX, areaY int) {
	startX := areaX * (AreaSize + MountainWidth)
	startY := areaY * (AreaSize + MountainWidth)

	// Add some random mountain tiles
	for i := 0; i < 5; i++ {
		x := startX + rand.Intn(AreaSize)
		y := startY + rand.Intn(AreaSize)
		if !isMountainTile(x, y) {
			m.Tiles[y][x] = TileMountain
		}
	}
}