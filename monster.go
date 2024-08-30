package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/solarlune/resolv"
    "math"
    "math/rand"
    "time"
)

type Monster struct {
    X, Y           float64
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
    collider       *resolv.Object
}

func NewMonster(x, y float64, sprite *ebiten.Image) *Monster {
    m := &Monster{
        X:              x,
        Y:              y,
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
    m.collider = resolv.NewObject(x, y, float64(TileSize), float64(TileSize))
    m.collider.SetShape(resolv.NewRectangle(0, 0, float64(TileSize), float64(TileSize)))
    m.collider.AddTags("monster")
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
        // Move towards the character if they're within twice the attack range
        m.MoveTowards(nearestChar.X, nearestChar.Y, w)
    } else {
        m.MoveRandomly()
    }
}

func (m *Monster) FindNearestCharacter(w *World) (*Character, float64) {
    var nearestChar *Character
    minDistance := math.Inf(1)

    for i := range w.characters {
        char := &w.characters[i]
        dx := char.X - m.X
        dy := char.Y - m.Y
        distance := math.Sqrt(dx*dx + dy*dy)

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

func (m *Monster) MoveTowards(targetX, targetY float64, w *World) {
    dx := targetX - m.X
    dy := targetY - m.Y
    distance := math.Sqrt(dx*dx + dy*dy)

    if distance > 0 {
        m.Direction.X = dx / distance
        m.Direction.Y = dy / distance
    }

    newX := m.X + m.Direction.X*m.Speed
    newY := m.Y + m.Direction.Y*m.Speed

    m.TryMove(newX, newY)
}

func (m *Monster) TryMove(newX, newY float64) bool {
    dx := newX - m.X
    dy := newY - m.Y

    if collision := m.collider.Check(dx, dy, "mountain", "boundary"); collision == nil {
        m.X = newX
        m.Y = newY
        m.collider.Position.X = m.X
        m.collider.Position.Y = m.Y
        m.collider.Update()
        return true
    } else {
        return false
    }
}
func (m *Monster) MoveRandomly() {
    newX := m.X + m.Direction.X*m.Speed
    newY := m.Y + m.Direction.Y*m.Speed

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
        m.collider.Space.Remove(m.collider)
    }
}
