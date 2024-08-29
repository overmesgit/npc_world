package main

import "github.com/hajimehoshi/ebiten/v2"

type Renderer struct{}

func NewRenderer() *Renderer {
    return &Renderer{}
}

func (r *Renderer) Render(screen *ebiten.Image, world *World, camera *Camera) {
    // Implement rendering logic
    // Draw map
    // Draw characters
    //	for _, char := range world.GetCharacters() {
    //		 Draw character at (char.X - camera.X, char.Y - camera.Y)
    //	}
    // Draw monsters
}
