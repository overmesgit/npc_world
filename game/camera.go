package game

import (
    "example.com/maj/units"
    "github.com/solarlune/resolv"
)

type Camera struct {
    Position resolv.Vector
}

func NewCamera() *Camera {
    return &Camera{
        Position: resolv.NewVector(0, 0),
    }
}

func (c *Camera) Update(player *units.Character) {
    if player != nil {
        playerPos := player.Object.Position
        c.Position.X = playerPos.X - 320 // Assuming 640x480 screen
        c.Position.Y = playerPos.Y - 240
    }
}

func (c *Camera) WorldToScreen(worldX, worldY float64) (screenX, screenY float64) {
    return worldX - c.Position.X, worldY - c.Position.Y
}
