package main

import "github.com/hajimehoshi/ebiten/v2"

type InputHandler struct {}

func NewInputHandler() *InputHandler {
	return &InputHandler{}
}

func (ih *InputHandler) HandleInput(world *World) {
	player := world.GetPlayerCharacter()
	if player == nil {
		return
	}

	if ebiten.IsKeyPressed(ebiten.KeyArrowLeft) {
		player.Move(-1, 0, world)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowRight) {
		player.Move(1, 0, world)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		player.Move(0, -1, world)
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		player.Move(0, 1, world)
	}
}