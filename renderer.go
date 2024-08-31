package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/text"
    "image/color"
    "golang.org/x/image/font"
    "golang.org/x/image/font/opentype"
    "golang.org/x/image/font/gofont/goregular"
    "log"
)

type Renderer struct {
    font font.Face
}

func NewRenderer() *Renderer {
    tt, err := opentype.Parse(goregular.TTF)
    if err != nil {
        log.Fatal(err)
    }
    const dpi = 72
    font, err := opentype.NewFace(tt, &opentype.FaceOptions{
        Size:    12,
        DPI:     dpi,
        Hinting: font.HintingFull,
    })
    if err != nil {
        log.Fatal(err)
    }

    return &Renderer{
        font: font,
    }
}

func (r *Renderer) Render(screen *ebiten.Image, world *World, camera *Camera) {
    // Clear the screen
    screen.Fill(color.RGBA{135, 206, 235, 255}) // Sky blue background

    // Draw game map
    for y := 0; y < world.gameMap.Height; y++ {
        for x := 0; x < world.gameMap.Width; x++ {
            r.drawTile(screen, x, y, world.gameMap.Tiles[y][x], camera)
        }
    }

    // Draw characters
    for _, char := range world.GetCharacters() {
        r.drawCharacter(screen, &char, camera)
    }

    // Draw monsters
    for _, monster := range world.monsters {
        r.drawMonster(screen, monster, camera)
    }
}

func (r *Renderer) drawTile(screen *ebiten.Image, x, y int, tileType TileType, camera *Camera) {
    worldX := float64(x * TileSize)
    worldY := float64(y * TileSize)
    screenX, screenY := camera.WorldToScreen(worldX, worldY)

    switch tileType {
    case TileGrass:
        ebitenutil.DrawRect(screen, screenX, screenY, float64(TileSize), float64(TileSize), color.RGBA{34, 139, 34, 255}) // Forest green
    case TileMountain:
        ebitenutil.DrawRect(screen, screenX, screenY, float64(TileSize), float64(TileSize), color.RGBA{139, 69, 19, 255}) // Saddle brown
    }
}

func (r *Renderer) drawCharacter(screen *ebiten.Image, char *Character, camera *Camera) {
    pos := char.Object.Position
    screenX, screenY := camera.WorldToScreen(pos.X, pos.Y)

    // Draw character sprite
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(screenX, screenY)
    screen.DrawImage(char.Sprite, op)

    // Draw health bar
    r.drawHealthBar(screen, screenX, screenY-10, char.Width, 5, char.Health, char.MaxHealth)

    // Draw character name
    text.Draw(screen, char.Name, r.font, int(screenX), int(screenY)-15, color.White)

    // Draw attack message if attacking
    if char.Attack.IsAttacking {
        text.Draw(screen, char.Attack.Message, r.font, int(screenX), int(screenY)-30, color.RGBA{255, 0, 0, 255})
    }
}

func (r *Renderer) drawMonster(screen *ebiten.Image, monster *Monster, camera *Camera) {
    pos := monster.Object.Position
    screenX, screenY := camera.WorldToScreen(pos.X, pos.Y)

    // Draw monster sprite
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(screenX, screenY)
    screen.DrawImage(monster.Sprite, op)

    // Draw health bar
    r.drawHealthBar(screen, screenX, screenY-10, monster.Width, 5, monster.Health, monster.MaxHealth)
}

func (r *Renderer) drawHealthBar(screen *ebiten.Image, x, y, width, height float64, health, maxHealth int) {
    // Draw background (empty health bar)
    ebitenutil.DrawRect(screen, x, y, width, height, color.RGBA{255, 0, 0, 255})

    // Draw filled portion of health bar
    healthPercentage := float64(health) / float64(maxHealth)
    filledWidth := width * healthPercentage
    ebitenutil.DrawRect(screen, x, y, filledWidth, height, color.RGBA{0, 255, 0, 255})
}
