package main

import (
    "example.com/maj/game"
    "example.com/maj/units"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/solarlune/resolv"
    "log"

    "github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
    world        *game.World
    camera       *game.Camera
    renderer     *game.Renderer
    inputHandler *game.InputHandler
    space        *resolv.Space
}

func NewGame() *Game {
    monsterSprite, _, err := ebitenutil.NewImageFromFile("assets/monsters/tile_0_0.png")
    if err != nil {
        log.Fatal(err)
    }

    playerSprite, _, err := ebitenutil.NewImageFromFile("assets/rogues/tile_0_1.png")
    if err != nil {
        log.Fatal(err)
    }
    npcSprite, _, err := ebitenutil.NewImageFromFile("assets/rogues/tile_0_1.png")
    if err != nil {
        log.Fatal(err)
    }

    goblinDenSprite, _, err := ebitenutil.NewImageFromFile("assets/tiles/tile_7_16.png")
    if err != nil {
        log.Fatal(err)
    }

    world := game.NewWorld(monsterSprite, goblinDenSprite)

    world.AddCharacter(units.NewCharacter(float64(3*game.TileSize), float64(3*game.TileSize), "Player"))
    world.AddCharacter(units.NewCharacter(float64(1*game.TileSize), float64(1*game.TileSize), "NPC1"))
    world.AddCharacter(units.NewCharacter(float64(2*game.TileSize), float64(2*game.TileSize), "NPC2"))

    sprites := map[string]*ebiten.Image{
        "NPC":     npcSprite,
        "PLAYER":  playerSprite,
        "DEN":     goblinDenSprite,
        "MONSTER": monsterSprite,
    }

    return &Game{
        world:        world,
        camera:       game.NewCamera(),
        renderer:     game.NewRenderer(sprites),
        inputHandler: game.NewInputHandler(),
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
