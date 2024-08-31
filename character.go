package main

import (
    "github.com/solarlune/resolv"
    "math"

    "github.com/hajimehoshi/ebiten/v2"
)

type Character struct {
    Name      string
    Speed     float64
    IsPlayer  bool
    Sprite    *ebiten.Image
    Width     float64
    Height    float64
    Attack    Attack
    Health    int
    MaxHealth int
    Object    *resolv.Object
}

func NewCharacter(x, y float64, name string, sprite *ebiten.Image) Character {
    c := Character{
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
    c.Object = resolv.NewObject(x, y, float64(TileSize), float64(TileSize))
    c.Object.SetShape(resolv.NewRectangle(0, 0, float64(TileSize), float64(TileSize)))
    c.Object.AddTags("character")
    return c
}

func (c *Character) Move(dx, dy float64, w *World) {
    if dx != 0 && dy != 0 {
        magnitude := math.Sqrt(dx*dx + dy*dy)
        dx /= magnitude
        dy /= magnitude
    }

    newX := c.Object.Position.X + dx*c.Speed
    newY := c.Object.Position.Y + dy*c.Speed

    if collision := c.Object.Check(dx*c.Speed, dy*c.Speed, "mountain"); collision == nil {
        c.Object.Position.X = newX
        c.Object.Position.Y = newY
        c.Object.Update()
    }
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
        distance := c.Object.Position.Distance(monster.Object.Position)

        if distance <= c.Attack.Range {
            monster.TakeDamage(c.Attack.Damage)
        }
    }
}
