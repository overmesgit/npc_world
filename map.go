package main

import (
    "log"
    "os"
    "strconv"
    "strings"
)

const (
    TileSize = 32
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
    tiles, width, height := loadMapFromFile("map.txt")
    return &GameMap{
        Tiles:  tiles,
        Width:  width,
        Height: height,
    }
}

func loadMapFromFile(filename string) ([][]TileType, int, int) {
    content, err := os.ReadFile(filename)
    if err != nil {
        log.Fatal(err)
    }

    lines := strings.Split(string(content), "\n")
    height := len(lines)
    width := 0
    tiles := make([][]TileType, height)

    for y, line := range lines {
        if line == "" {
            continue
        }
        values := strings.Split(line, ",")
        if width == 0 {
            width = len(values)
        }
        tiles[y] = make([]TileType, width)
        for x, value := range values {
            tileType, err := strconv.Atoi(value)
            if err != nil {
                log.Fatal(err)
            }
            tiles[y][x] = TileType(tileType)
        }
    }

    return tiles, width, height
}
