package main

import (
    "fmt"
    "math"
)

type Character struct {
	X, Y     float64
	Name     string
	Speed    float64
	IsPlayer bool
}

func NewCharacter(x, y float64, name string) Character {
	return Character{
		X:        x,
		Y:        y,
		Name:     name,
		Speed:    2.0,
		IsPlayer: name == "Player",
	}
}

func (c *Character) Update(w *World) {
	// Update character logic
	// This could include AI for NPCs, or be empty for the player character
	// as their movement is handled by input
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
	fmt.Println(newX, newY)
	c.X = newX
	c.Y = newY
	// Simple collision detection with game boundaries
	if newX >= 0 && newX < float64(w.gameMap.Width) {
		c.X = newX
	}
	if newY >= 0 && newY < float64(w.gameMap.Height) {
		c.Y = newY
	}
}

