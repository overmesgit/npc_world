package main

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type InputHandler struct{}

func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

func (ih *InputHandler) HandleInput(world *World) {
	player := world.GetPlayerCharacter()
	if player == nil {
		return
	}

	dx, dy := 0.0, 0.0

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		dx -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		dx += 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		dy -= 1
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		dy += 1
	}

	player.Move(dx, dy, world)
}

