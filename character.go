package main

import (
    "github.com/solarlune/resolv"
    "math"

    "github.com/hajimehoshi/ebiten/v2"
)

type Character struct {
    X, Y      float64
    Name      string
    Speed     float64
    IsPlayer  bool
    Sprite    *ebiten.Image
    Width     float64
    Height    float64
    Attack    Attack
    Health    int
    MaxHealth int
    collider  *resolv.Object
}

func NewCharacter(x, y float64, name string, sprite *ebiten.Image) Character {
    c := Character{
        X:         x,
        Y:         y,
        Name:      name,
        Speed:     2.0,
        IsPlayer:  name == "Player",
        Sprite:    sprite,
        Width:     float64(TileSize),
        Height:    float64(TileSize),
        Attack:    NewAttack(),
        Health:    100,
        MaxHealth: 100,
    }
    c.collider = resolv.NewObject(x, y, float64(TileSize), float64(TileSize))
    c.collider.SetShape(resolv.NewRectangle(0, 0, float64(TileSize), float64(TileSize)))
    c.collider.AddTags("character")
    return c
}

func (c *Character) TakeDamage(amount int) {
    c.Health -= amount
    if c.Health < 0 {
        c.Health = 0
    }
}

func (c *Character) Update(w *World) {
    c.Attack.Update()
    if c.Attack.IsAttacking && !c.Attack.HasDealtDamage {
        c.PerformAttack(w)
        c.Attack.HasDealtDamage = true // Set this flag after dealing damage
    }
}

func (c *Character) PerformAttack(w *World) {
    for _, monster := range w.monsters {
        dx := monster.X - c.X
        dy := monster.Y - c.Y
        distance := math.Sqrt(dx*dx + dy*dy)

        if distance <= c.Attack.Range {
            monster.TakeDamage(c.Attack.Damage)
        }
    }
}

func (c *Character) Move(dx, dy float64, w *World) {
    if dx != 0 && dy != 0 {
        magnitude := math.Sqrt(dx*dx + dy*dy)
        dx /= magnitude
        dy /= magnitude
    }

    newX := c.X + dx*c.Speed
    newY := c.Y + dy*c.Speed
    c.collider.Position.X = c.X
    c.collider.Position.Y = c.Y

    if collision := c.collider.Check(dx*c.Speed, dy*c.Speed, "mountain", "boundary"); collision == nil {
        c.X = newX
        c.Y = newY
        c.collider.Position.X = c.X
        c.collider.Position.Y = c.Y
    } else {
        // Move as far as possible before collision
        //        diff := collision.ContactWithObject(collision.Objects[0])
        //        c.X += diff.X
        //        c.Y += diff.Y
        //        c.collider.Position.X = c.X
        //        c.collider.Position.Y = c.Y
    }
    c.collider.Update()
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
