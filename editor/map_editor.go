package main

import (
    "fmt"
    "image/color"
    "log"
    "os"
    "strconv"
    "strings"

    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
    ScreenWidth  = 800
    ScreenHeight = 600
    TileSize     = 32
    MapWidth     = 20
    MapHeight    = 20
)

type TileType int

const (
    TileGrass TileType = iota
    TileMountain
)

type Editor struct {
	tiles         [][]TileType
	currentTile   TileType
	saveRequested bool
	saveMessage   string
	messageTimer  int
}

func NewEditor() *Editor {
	tiles := make([][]TileType, MapHeight)
	for i := range tiles {
		tiles[i] = make([]TileType, MapWidth)
	}
	return &Editor{
		tiles:       tiles,
		currentTile: TileGrass,
	}
}

func (e *Editor) Update() error {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		x, y := ebiten.CursorPosition()
		tileX, tileY := x/TileSize, y/TileSize
		if tileX >= 0 && tileX < MapWidth && tileY >= 0 && tileY < MapHeight {
			e.tiles[tileY][tileX] = e.currentTile
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.Key1) {
		e.currentTile = TileGrass
	} else if inpututil.IsKeyJustPressed(ebiten.Key2) {
		e.currentTile = TileMountain
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyS) {
		err := e.SaveMap("map.txt")
		if err != nil {
			e.saveMessage = "Error saving map: " + err.Error()
		} else {
			e.saveMessage = "Map saved successfully!"
		}
		e.messageTimer = 180 // Show message for 3 seconds (60 frames per second)
	}

	if e.messageTimer > 0 {
		e.messageTimer--
	}

	return nil
}

func (e *Editor) Draw(screen *ebiten.Image) {
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			switch e.tiles[y][x] {
			case TileGrass:
				ebitenutil.DrawRect(screen, float64(x*TileSize), float64(y*TileSize), TileSize, TileSize, color.RGBA{34, 139, 34, 255})
			case TileMountain:
				ebitenutil.DrawRect(screen, float64(x*TileSize), float64(y*TileSize), TileSize, TileSize, color.RGBA{139, 69, 19, 255})
			}
		}
	}

	ebitenutil.DebugPrint(screen, fmt.Sprintf("Current Tile: %v (1:Grass, 2:Mountain) | Press 'S' to save", e.currentTile))

	if e.messageTimer > 0 {
		ebitenutil.DebugPrintAt(screen, e.saveMessage, 10, ScreenHeight-20)
	}
}

func (e *Editor) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return ScreenWidth, ScreenHeight
}

func (e *Editor) SaveMap(filename string) error {
	var sb strings.Builder
	for y := 0; y < MapHeight; y++ {
		for x := 0; x < MapWidth; x++ {
			sb.WriteString(strconv.Itoa(int(e.tiles[y][x])))
			if x < MapWidth-1 {
				sb.WriteString(",")
			}
		}
		sb.WriteString("\n")
	}

	return os.WriteFile(filename, []byte(sb.String()), 0644)
}

func main() {
	editor := NewEditor()
	ebiten.SetWindowSize(ScreenWidth, ScreenHeight)
	ebiten.SetWindowTitle("Map Editor")

	if err := ebiten.RunGame(editor); err != nil {
		log.Fatal(err)
	}
}