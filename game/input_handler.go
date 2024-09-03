package game

import (
	"github.com/hajimehoshi/ebiten/v2"
    "github.com/hajimehoshi/ebiten/v2/inpututil"
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

	player.Move(dx, dy)

	// Handle attack input
	if inpututil.IsKeyJustPressed(ebiten.KeyControl) {
		player.Attack.TriggerAttack()
	}
}

