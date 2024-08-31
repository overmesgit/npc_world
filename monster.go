package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/solarlune/resolv"
    "math"
    "math/rand"
    "time"
)

type Monster struct {
    Width, Height  float64
    Speed          float64
    Direction      struct{ X, Y float64 }
    Health         int
    MaxHealth      int
    Sprite         *ebiten.Image
    AttackRange    float64
    AttackDamage   int
    AttackCooldown time.Duration
    LastAttackTime time.Time
    Object         *resolv.Object
}

func NewMonster(x, y float64, sprite *ebiten.Image) *Monster {
    m := &Monster{
        Width:          float64(TileSize),
        Height:         float64(TileSize),
        Speed:          1.0,
        Direction:      struct{ X, Y float64 }{X: rand.Float64()*2 - 1, Y: rand.Float64()*2 - 1},
        Health:         100,
        MaxHealth:      100,
        Sprite:         sprite,
        AttackRange:    float64(TileSize * 1.5), // 1.5 tiles range
        AttackDamage:   10,
        AttackCooldown: time.Second * 2, // Attack every 2 seconds
    }
    m.Object = resolv.NewObject(x, y, float64(TileSize), float64(TileSize))
    m.Object.SetShape(resolv.NewRectangle(0, 0, float64(TileSize), float64(TileSize)))
    m.Object.AddTags("monster")
    return m
}

func (m *Monster) Update(w *World) {
    if m.Health <= 0 {
        return
    }

    nearestChar, distance := m.FindNearestCharacter(w)
    if nearestChar != nil && distance <= m.AttackRange {
        m.AttackCharacter(nearestChar)
    } else if nearestChar != nil && distance <= m.AttackRange*2 {
        m.MoveTowards(nearestChar.Object)
    } else {
        m.MoveRandomly()
    }
}

func (m *Monster) FindNearestCharacter(w *World) (*Character, float64) {
    var nearestChar *Character
    minDistance := math.Inf(1)

    for i := range w.characters {
        char := &w.characters[i]
        distance := char.Object.Position.Distance(m.Object.Position)

        if distance < minDistance {
            minDistance = distance
            nearestChar = char
        }
    }

    return nearestChar, minDistance
}

func (m *Monster) AttackCharacter(char *Character) {
    if time.Since(m.LastAttackTime) >= m.AttackCooldown {
        char.TakeDamage(m.AttackDamage)
        m.LastAttackTime = time.Now()
    }
}

func (m *Monster) MoveTowards(object *resolv.Object) {
    position := m.Object.Position
    distance := position.Distance(object.Position)

    if distance > 0 {
        sub := object.Position.Sub(position)
        m.Direction.X = sub.X / distance
        m.Direction.Y = sub.Y / distance
    }

    newX := position.X + m.Direction.X*m.Speed
    newY := position.Y + m.Direction.Y*m.Speed

    m.TryMove(newX, newY)
}

func (m *Monster) TryMove(newX, newY float64) bool {
    position := m.Object.Position
    dx := newX - position.X
    dy := newY - position.Y

    if collision := m.Object.Check(dx, dy, "mountain", "character"); collision == nil {
        m.Object.Position.X = newX
        m.Object.Position.Y = newY
        m.Object.Update()
        return true
    }
    return false
}

func (m *Monster) MoveRandomly() {
    position := m.Object.Position
    newX := position.X + m.Direction.X*m.Speed
    newY := position.Y + m.Direction.Y*m.Speed

    if !m.TryMove(newX, newY) {
        // Change direction if hit an obstacle
        m.Direction.X = rand.Float64()*2 - 1
        m.Direction.Y = rand.Float64()*2 - 1
        m.NormalizeDirection()
    }
}

func (m *Monster) NormalizeDirection() {
    magnitude := math.Sqrt(m.Direction.X*m.Direction.X + m.Direction.Y*m.Direction.Y)
    if magnitude != 0 {
        m.Direction.X /= magnitude
        m.Direction.Y /= magnitude
    }
}

func (m *Monster) TakeDamage(amount int) {
    m.Health -= amount
    if m.Health < 0 {
        m.Health = 0
        m.Object.Space.Remove(m.Object)
    }
}
