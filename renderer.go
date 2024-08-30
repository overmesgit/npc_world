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

func (r *Renderer) drawMonster(screen *ebiten.Image, monster *Monster, camera *Camera) {
    screenX := monster.X - camera.X
    screenY := monster.Y - camera.Y

    // Draw monster sprite
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(screenX, screenY)
    screen.DrawImage(monster.Sprite, op)

    // Draw health bar
    healthBarWidth := monster.Width
    healthBarHeight := 5.0
    healthBarY := screenY - healthBarHeight - 2 // 2 pixels above the monster

    // Draw background (empty health bar)
    ebitenutil.DrawRect(screen, screenX, healthBarY, healthBarWidth, healthBarHeight, color.RGBA{255, 0, 0, 255})

    // Draw filled portion of health bar
    healthPercentage := float64(monster.Health) / float64(monster.MaxHealth)
    filledWidth := healthBarWidth * healthPercentage
    ebitenutil.DrawRect(screen, screenX, healthBarY, filledWidth, healthBarHeight, color.RGBA{0, 255, 0, 255})
}

func (r *Renderer) drawTile(screen *ebiten.Image, x, y int, tileType TileType, camera *Camera) {
    screenX := float64(x*TileSize) - camera.X
    screenY := float64(y*TileSize) - camera.Y

    switch tileType {
    case TileGrass:
        ebitenutil.DrawRect(screen, screenX, screenY, float64(TileSize), float64(TileSize), color.RGBA{34, 139, 34, 255}) // Forest green
    case TileMountain:
        ebitenutil.DrawRect(screen, screenX, screenY, float64(TileSize), float64(TileSize), color.RGBA{139, 69, 19, 255}) // Saddle brown
    }
}

func (r *Renderer) drawCharacter(screen *ebiten.Image, char *Character, camera *Camera) {
    // Calculate screen position
    screenX := char.X - camera.X
    screenY := char.Y - camera.Y

    // Draw character sprite
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(screenX, screenY)
    screen.DrawImage(char.Sprite, op)

    // Draw character name
    text.Draw(screen, char.Name, r.font, int(screenX), int(screenY)-5, color.White)

    // Draw attack message if attacking
    if char.Attack.IsAttacking {
        text.Draw(screen, char.Attack.Message, r.font, int(screenX), int(screenY)-20, color.RGBA{255, 0, 0, 255})
    }
}
