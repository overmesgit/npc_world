package main

import (
    "math"

    "github.com/hajimehoshi/ebiten/v2"
)

type Character struct {
    X, Y     float64
    Name     string
    Speed    float64
    IsPlayer bool
    Sprite   *ebiten.Image
    Width    float64
    Height   float64
}

func NewCharacter(x, y float64, name string, sprite *ebiten.Image) Character {
    return Character{
        X:        x,
        Y:        y,
        Name:     name,
        Speed:    2.0,
        IsPlayer: name == "Player",
        Sprite:   sprite,
        Width:    float64(TileSize), // Assuming character size is same as tile size
        Height:   float64(TileSize),
    }
}

func (c *Character) Update(w *World) {
    // Update character logic
}

func (c *Character) Move(dx, dy float64, w *World) {
    // Normalize diagonal movement
    if dx != 0 && dy != 0 {
        magnitude := math.Sqrt(dx*dx + dy*dy)
        dx /= magnitude
        dy /= magnitude
    }

    newX := c.X + dx*c.Speed
    newY := c.Y + dy*c.Speed

    // Check collision for X movement
    if !c.collidesWithMountain(newX, c.Y, w) {
        c.X = newX
    }

    // Check collision for Y movement
    if !c.collidesWithMountain(c.X, newY, w) {
        c.Y = newY
    }
}

func (c *Character) collidesWithMountain(x, y float64, w *World) bool {
    // Check all four corners of the character
    corners := [][2]float64{
        {x, y},                              // Top-left
        {x + c.Width - 1, y},                // Top-right
        {x, y + c.Height - 1},               // Bottom-left
        {x + c.Width - 1, y + c.Height - 1}, // Bottom-right
    }

    for _, corner := range corners {
        tileX := int(corner[0] / TileSize)
        tileY := int(corner[1] / TileSize)
        if !w.gameMap.IsTileWalkable(tileX, tileY) {
            return true
        }
    }

    return false
}
