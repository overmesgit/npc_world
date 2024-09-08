package game

import (
    gamemap "example.com/maj/map"
    "example.com/maj/units"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/ebitenutil"
    "github.com/hajimehoshi/ebiten/v2/text"
    "github.com/hajimehoshi/ebiten/v2/vector"
    "golang.org/x/image/font"
    "golang.org/x/image/font/gofont/goregular"
    "golang.org/x/image/font/opentype"
    "image"
    "image/color"
    "log"
)

type Renderer struct {
    font    font.Face
    sprites Sprites
}

type Sprites struct {
    Monsters    *ebiten.Image
    Characteres *ebiten.Image
    Tiles       *ebiten.Image
}

func NewRenderer(sprites Sprites) *Renderer {
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
        sprites: sprites,
        font:    font,
    }
}

func (r *Renderer) Render(screen *ebiten.Image, world *World, camera *Camera) {
    // Clear the screen
    screen.Fill(color.RGBA{135, 206, 235, 255}) // Sky blue background

    // Draw game map
    for y := 0; y < world.GameMap.Height; y++ {
        for x := 0; x < world.GameMap.Width; x++ {
            r.drawTile(screen, x, y, world.GameMap.Tiles[y][x], camera)
        }
    }

    for _, obj := range world.Space.Objects() {
        switch obj.Data.(type) {
        case *units.Character:
            character := obj.Data.(*units.Character)
            r.drawCharacter(screen, character, camera)
            r.drawViewField(screen, character, camera)
        case *units.Monster:
            monster := obj.Data.(*units.Monster)
            r.drawMonster(screen, monster, camera)
        case *units.GoblinDen:
            den := obj.Data.(*units.GoblinDen)
            r.drawGoblinDen(screen, den, camera)
        case *units.Mushroom:
            mushroom := obj.Data.(*units.Mushroom)
            r.drawMushroom(screen, mushroom, camera)
        }
    }
}

func (r *Renderer) drawSprite(screen *ebiten.Image, sheet *ebiten.Image, indexX, indexY int, x, y float64) {
    sx := indexX * 32
    sy := indexY * 32
    op := &ebiten.DrawImageOptions{}
    op.GeoM.Translate(x, y)
    screen.DrawImage(sheet.SubImage(image.Rect(sx, sy, sx+32, sy+32)).(*ebiten.Image), op)
}

func (r *Renderer) drawGoblinDen(screen *ebiten.Image, den *units.GoblinDen, camera *Camera) {
    pos := den.Object.Position
    screenX, screenY := camera.WorldToScreen(pos.X, pos.Y)

    r.drawSprite(screen, r.sprites.Tiles, 0, 16, screenX, screenY)

    // Draw health bar for goblin den
    r.drawHealthBar(screen, screenX, screenY-10, float64(gamemap.TileSize), 5, den.Health, den.MaxHealth)
}

func (r *Renderer) drawMushroom(screen *ebiten.Image, mushroom *units.Mushroom, camera *Camera) {
    pos := mushroom.Object.Position
    screenX, screenY := camera.WorldToScreen(pos.X, pos.Y)
    r.drawSprite(screen, r.sprites.Tiles, 0, 20, screenX, screenY)
}

func (r *Renderer) drawTile(screen *ebiten.Image, x, y int, tileType gamemap.TileType, camera *Camera) {
    worldX := float64(x * gamemap.TileSize)
    worldY := float64(y * gamemap.TileSize)
    screenX, screenY := camera.WorldToScreen(worldX, worldY)

    switch tileType {
    case gamemap.TileGrass:
        ebitenutil.DrawRect(screen, screenX, screenY, float64(gamemap.TileSize), float64(gamemap.TileSize), color.RGBA{34, 139, 34, 255}) // Forest green
    case gamemap.TileMountain:
        ebitenutil.DrawRect(screen, screenX, screenY, float64(gamemap.TileSize), float64(gamemap.TileSize), color.RGBA{139, 69, 19, 255}) // Saddle brown
    }

    borderColor := color.RGBA{0, 0, 0, 255} // Black border
    vector.StrokeRect(screen,
        float32(screenX),
        float32(screenY),
        float32(gamemap.TileSize),
        float32(gamemap.TileSize),
        1,
        borderColor,
        false)
}

func (r *Renderer) drawCharacter(screen *ebiten.Image, char *units.Character, camera *Camera) {
    pos := char.Object.Position
    screenX, screenY := camera.WorldToScreen(pos.X, pos.Y)

    // Draw character sprite
    r.drawSprite(screen, r.sprites.Characteres, 0, 1, screenX, screenY)

    // Draw health bar
    r.drawHealthBar(screen, screenX, screenY-10, char.Width, 5, char.Health, char.MaxHealth)

    // Draw character name
    text.Draw(screen, char.Name, r.font, int(screenX), int(screenY)-15, color.White)

    // Draw attack message if attacking
    if char.Attack.IsAttacking {
        text.Draw(screen, char.Attack.Message, r.font, int(screenX), int(screenY)-30, color.RGBA{255, 0, 0, 255})
    }
}

func (r *Renderer) drawMonster(screen *ebiten.Image, monster *units.Monster, camera *Camera) {
    pos := monster.Object.Position
    screenX, screenY := camera.WorldToScreen(pos.X, pos.Y)

    // Draw monster sprite
    r.drawSprite(screen, r.sprites.Monsters, 0, 0, screenX, screenY)

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

func (r *Renderer) drawViewField(screen *ebiten.Image, character *units.Character, camera *Camera) {
    viewRange := float64(3 * gamemap.TileSize)

    checkX := character.Object.Center().X - viewRange
    checkY := character.Object.Center().Y - viewRange
    checkSize := viewRange * 2

    screenX, screenY := camera.WorldToScreen(checkX, checkY)

    //    ebitenutil.DrawRect(
    //        screen,
    //        screenX,
    //        screenY,
    //        checkSize,
    //        checkSize,
    //        color.RGBA{0, 255, 0, 64},
    //    )

    vector.StrokeRect(screen,
        float32(screenX), float32(screenY),
        float32(checkSize), float32(checkSize),
        1,                          // Line width
        color.RGBA{0, 255, 0, 255}, // Solid green for the border
        false) // Disable anti-aliasing for a crisp border
}
