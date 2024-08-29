package main

import (
    "fmt"
    "log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	world        *World
	camera       *Camera
	renderer     *Renderer
	inputHandler *InputHandler
}

func NewGame() *Game {
	world := NewWorld()
	// Create some initial characters (now using tile coordinates)
	world.AddCharacter(NewCharacter(float64(3*TileSize), float64(3*TileSize), "Player"))
	world.AddCharacter(NewCharacter(float64(6*TileSize), float64(6*TileSize), "NPC1"))
	world.AddCharacter(NewCharacter(float64(9*TileSize), float64(9*TileSize), "NPC2"))

	return &Game{
		world:        world,
		camera:       NewCamera(),
		renderer:     NewRenderer(),
		inputHandler: NewInputHandler(),
	}
}

func (g *Game) Update() error {
	g.inputHandler.HandleInput(g.world)
	g.world.Update()
	g.camera.Update(g.world.GetPlayerCharacter())
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Render(screen, g.world, g.camera)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("My 2D Top-Down Game")
	if err := ebiten.RunGame(NewGame()); err != nil {
		log.Fatal(err)
	}
}