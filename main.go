package main

import (
    "example.com/maj/game"
    gamemap "example.com/maj/map"
    "example.com/maj/units"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
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
    isPaused     bool
}

func NewGame() *Game {

    world := game.NewWorld()
    world.AddCharacter(units.NewCharacter(float64(3*gamemap.TileSize), float64(3*gamemap.TileSize), "Player"))
    world.AddCharacter(units.NewCharacter(float64(4*gamemap.TileSize), float64(4*gamemap.TileSize), "NPC1"))
    world.AddCharacter(units.NewCharacter(float64(4*gamemap.TileSize), float64(4*gamemap.TileSize), "NPC2"))
    world.AddCharacter(units.NewCharacter(float64(4*gamemap.TileSize), float64(4*gamemap.TileSize), "NPC3"))
    world.AddCharacter(units.NewCharacter(float64(4*gamemap.TileSize), float64(4*gamemap.TileSize), "NPC4"))
    world.AddCharacter(units.NewCharacter(float64(4*gamemap.TileSize), float64(4*gamemap.TileSize), "NPC5"))

    chars, _, err := ebitenutil.NewImageFromFile("assets/rogues.png")
    if err != nil {
        log.Fatal(err)
    }
    monsters, _, err := ebitenutil.NewImageFromFile("assets/monsters.png")
    if err != nil {
        log.Fatal(err)
    }
    tiles, _, err := ebitenutil.NewImageFromFile("assets/tiles.png")
    if err != nil {
        log.Fatal(err)
    }

    return &Game{
        world:  world,
        camera: game.NewCamera(),
        renderer: game.NewRenderer(game.Sprites{
            monsters,
            chars,
            tiles,
        }),
        inputHandler: game.NewInputHandler(),
    }
}

func (g *Game) Update() error {
    if inpututil.IsKeyJustPressed(ebiten.KeySpace) {
        g.isPaused = !g.isPaused
    }

    if !g.isPaused {
        g.inputHandler.HandleInput(g.world)
        g.world.Update()
        g.camera.Update(g.world.GetPlayerCharacter())
    }
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
