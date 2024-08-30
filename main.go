package main

import (
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
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
    monsterSprite, _, err := ebitenutil.NewImageFromFile("assets/monsters/tile_0_0.png")
	if err != nil {
		log.Fatal(err)
	}

	world := NewWorld(monsterSprite)
    world.gameMap.AddRandomFeatures()

    // Load sprites
    playerSprite, _, err := ebitenutil.NewImageFromFile("assets/rogues/tile_0_0.png")
    if err != nil {
        log.Fatal(err)
    }
    npcSprite, _, err := ebitenutil.NewImageFromFile("assets/rogues/tile_0_1.png")
    if err != nil {
        log.Fatal(err)
    }

    // Create characters with sprites
    world.AddCharacter(NewCharacter(float64(3*TileSize), float64(3*TileSize), "Player", playerSprite))
    world.AddCharacter(NewCharacter(float64(6*TileSize), float64(6*TileSize), "NPC1", npcSprite))
    world.AddCharacter(NewCharacter(float64(9*TileSize), float64(9*TileSize), "NPC2", npcSprite))

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
    ebiten.SetWindowSize(1280, 960)
    ebiten.SetWindowTitle("My 2D Top-Down Game")
    if err := ebiten.RunGame(NewGame()); err != nil {
        log.Fatal(err)
    }
}
