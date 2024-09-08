package units

import (
    "github.com/solarlune/resolv"
    "math"
    "math/rand"
    "time"
)

type Monster struct {
    Width, Height float64
    Speed         float64
    Direction     struct{ X, Y float64 }
    Health        int
    MaxHealth     int
    Attack        Attack
    Object        *resolv.Object
    Den           *GoblinDen
    WanderRadius  float64
}

func NewMonster(x, y float64, den *GoblinDen) *Monster {
    attack := NewAttack(float64(32 * 1.5))
    attack.Damage = 10
    attack.CooldownDuration = time.Second * 2
    m := &Monster{
        Width:        float64(32),
        Height:       float64(32),
        Speed:        1.0,
        Direction:    struct{ X, Y float64 }{X: rand.Float64()*2 - 1, Y: rand.Float64()*2 - 1},
        Health:       100,
        MaxHealth:    100,
        Attack:       attack,
        Den:          den,
        WanderRadius: float64(32 * 5), // 5 tiles radius
    }
    m.Object = resolv.NewObject(x, y, float64(32), float64(32))
    m.Object.SetShape(resolv.NewRectangle(0, 0, float64(32), float64(32)))
    m.Object.AddTags("monster")
    m.Object.Data = m
    return m
}

func (m *Monster) Update() {
    if m.Health <= 0 {
        return
    }

    nearestChar, distance := m.FindNearestCharacter()
    if nearestChar != nil && distance <= m.Attack.Range {
        m.AttackCharacter(nearestChar)
    } else if nearestChar != nil && distance <= m.Attack.Range*4 {
        m.MoveTowards(nearestChar.Object)
    } else {
        m.WanderNearDen()
    }
}

func (m *Monster) FindNearestCharacter() (*Character, float64) {
    nearestObj, distance := FindNearest(m.Object, 2*m.Attack.Range, "character")
    if nearestObj != nil {
        return nearestObj.Data.(*Character), distance
    }
    return nil, 0
}

func (m *Monster) WanderNearDen() {
    denPos := m.Den.Object.Center()
    monsterPos := m.Object.Center()
    distanceToDen := monsterPos.Distance(denPos)

    if distanceToDen > m.WanderRadius {
        // Move back towards den
        m.MoveTowards(m.Den.Object)
    } else {
        // Wander randomly
        m.MoveRandomly()
    }
}

func (m *Monster) AttackCharacter(char *Character) {
    m.Attack.TriggerAttack()
    if m.Attack.IsAttacking && !m.Attack.HasDealtDamage {
        char.TakeDamage(m.Attack.Damage)
        m.Attack.HasDealtDamage = true
    }
}

func (m *Monster) MoveTowards(object *resolv.Object) {
    position := m.Object.Position
    distance := m.Object.Center().Distance(object.Center())

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

    collision := m.Object.Check(dx, dy, "mountain")
    if collision == nil {
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
    }
}
