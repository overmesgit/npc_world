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
	// Create some initial characters
	world.AddCharacter(NewCharacter(100, 100, "Player"))
	world.AddCharacter(NewCharacter(200, 200, "NPC1"))
	world.AddCharacter(NewCharacter(300, 300, "NPC2"))

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
	fmt.Println(g.world.GetPlayerCharacter())
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