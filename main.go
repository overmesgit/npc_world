// main.go
package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	world     *World
	camera    *Camera
	renderer  *Renderer
	inputHandler *InputHandler
}

func NewGame() *Game {
	return &Game{
		world:     NewWorld(),
		camera:    NewCamera(),
		renderer:  NewRenderer(),
		inputHandler: NewInputHandler(),
	}
}

func (g *Game) Update() error {
	g.inputHandler.HandleInput()
	g.world.Update()
	g.camera.Update()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.renderer.Render(screen, g.world, g.camera)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 640, 480
}

func main() {
	game := NewGame()
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("My 2D Top-Down Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}


