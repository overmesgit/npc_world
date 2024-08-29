package main

import (
    "github.com/hajimehoshi/ebiten/v2"
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
    screen.Fill(color.RGBA{0, 0, 0, 255})

    // Draw characters
    for _, char := range world.GetCharacters() {
        r.drawCharacter(screen, &char, camera)
    }
}

func (r *Renderer) drawCharacter(screen *ebiten.Image, char *Character, camera *Camera) {
    // Character size
    const charSize = 20

    // Calculate screen position
    screenX := char.X - camera.X
    screenY := char.Y - camera.Y

    // Draw green square
    square := ebiten.NewImage(charSize, charSize)
    square.Fill(color.RGBA{0, 255, 0, 255})

    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(screenX, screenY)
    screen.DrawImage(square, op)

    // Draw character name
    text.Draw(screen, char.Name, r.font, int(screenX), int(screenY)-5, color.White)
}
