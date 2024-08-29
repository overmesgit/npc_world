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

    // Draw game map (placeholder)
    for y := 0; y < world.gameMap.Height; y++ {
        for x := 0; x < world.gameMap.Width; x++ {
            r.drawTile(screen, x, y, camera)
        }
    }

    // Draw characters
    for _, char := range world.GetCharacters() {
        r.drawCharacter(screen, &char, camera)
    }
}

func (r *Renderer) drawTile(screen *ebiten.Image, x, y int, camera *Camera) {
    // Placeholder: Draw a simple grid
    tileColor := color.RGBA{200, 200, 200, 255} // Light gray
    ebitenutil.DrawRect(screen,
        float64(x*TileSize)-camera.X,
        float64(y*TileSize)-camera.Y,
        float64(TileSize),
        float64(TileSize),
        tileColor)
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
}
