package main

import (
    "math"
    "math/rand"
)

type Monster struct {
    X, Y      float64
    Width     float64
    Height    float64
    Speed     float64
    Direction struct{ X, Y float64 }
}

func NewMonster(x, y float64) *Monster {
    return &Monster{
        X:      x,
        Y:      y,
        Width:  float64(TileSize),
        Height: float64(TileSize),
        Speed:  1.0,
        Direction: struct{ X, Y float64 }{
            X: rand.Float64()*2 - 1,
            Y: rand.Float64()*2 - 1,
        },
    }
}

func (m *Monster) Update(w *World) {
    newX := m.X + m.Direction.X*m.Speed
    newY := m.Y + m.Direction.Y*m.Speed

    // Check all four corners of the monster
    corners := [][2]float64{
        {newX, newY},                              // Top-left
        {newX + m.Width - 1, newY},                // Top-right
        {newX, newY + m.Height - 1},               // Bottom-left
        {newX + m.Width - 1, newY + m.Height - 1}, // Bottom-right
    }

    canMove := true
    for _, corner := range corners {
        tileX := int(corner[0] / TileSize)
        tileY := int(corner[1] / TileSize)
        if !w.gameMap.IsTileWalkable(tileX, tileY) {
            canMove = false
            break
        }
    }

    if canMove {
        m.X = newX
        m.Y = newY
    } else {
        // Change direction if hit an obstacle
        m.Direction.X = rand.Float64()*2 - 1
        m.Direction.Y = rand.Float64()*2 - 1
    }
}

func (m *Monster) NormalizeDirection() {
    magnitude := math.Sqrt(m.Direction.X*m.Direction.X + m.Direction.Y*m.Direction.Y)
    if magnitude != 0 {
        m.Direction.X /= magnitude
        m.Direction.Y /= magnitude
    }
}


