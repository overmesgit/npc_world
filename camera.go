package main
type Camera struct {
	X, Y float64
}

func NewCamera() *Camera {
	return &Camera{}
}

func (c *Camera) Update(player *Character) {
	if player != nil {
		c.X = player.X - 320 // Assuming 640x480 screen
		c.Y = player.Y - 240
	}
}

