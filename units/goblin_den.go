package units

import (
    "github.com/solarlune/resolv"
    "math"
    "math/rand"
    "time"
)

type GoblinDen struct {
    Object          *resolv.Object
    SpawnCooldown   time.Duration
    LastSpawnTime   time.Time
    MaxMonsters     int
    CurrentMonsters int
    Health          int // Add this line
    MaxHealth       int // Add this line
}

func NewGoblinDen(x, y float64) *GoblinDen {
    den := &GoblinDen{
        SpawnCooldown:   time.Second * 10,
        LastSpawnTime:   time.Now(),
        MaxMonsters:     5,
        CurrentMonsters: 0,
        Health:          100, // Add this line
        MaxHealth:       100, // Add this line
    }
    size := float64(32)
    den.Object = resolv.NewObject(x, y, size, size)
    den.Object.SetShape(resolv.NewRectangle(0, 0, size, size))
    den.Object.AddTags("goblin_den")
    den.Object.Data = den
    return den
}

func (d *GoblinDen) Update() []*Monster {
    if time.Since(d.LastSpawnTime) >= d.SpawnCooldown && d.CurrentMonsters < d.MaxMonsters {
        monster := d.SpawnMonster()
        d.LastSpawnTime = time.Now()
        d.CurrentMonsters++
        return []*Monster{monster}
    }
    return nil
}

func (d *GoblinDen) SpawnMonster() *Monster {
    spawnRadius := float64(32 * 2)
    angle := rand.Float64() * 2 * math.Pi
    x := d.Object.Position.X + math.Cos(angle)*spawnRadius
    y := d.Object.Position.Y + math.Sin(angle)*spawnRadius
    return NewMonster(x, y, d)
}

func (d *GoblinDen) TakeDamage(amount int) {
    d.Health -= amount
    if d.Health < 0 {
        d.Health = 0
    }
}
