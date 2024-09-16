package units

import (
    "example.com/maj/ai"
    "fmt"
    "github.com/solarlune/resolv"
    "time"
)

type Character struct {
    Name          string
    Speed         float64
    IsPlayer      bool
    Width         float64
    Height        float64
    Attack        Attack
    Health        int
    MaxHealth     int
    Object        *resolv.Object
    Planner       *ai.GOAPPlanner
    CurrentPlan   []ai.GOAPAction
    TargetMonster *Monster
    WanderTarget  resolv.Vector
    WanderTime    time.Time
}

func NewCharacter(x, y float64, name string) *Character {
    c := &Character{
        Name:      name,
        Speed:     2.0,
        IsPlayer:  name == "Player",
        Width:     float64(32),
        Height:    float64(32),
        Attack:    NewAttack(2 * 32),
        Health:    100,
        MaxHealth: 100,
    }
    c.Object = resolv.NewObject(x, y, float64(32), float64(32))
    c.Object.SetShape(resolv.NewRectangle(0, 0, float64(32), float64(32)))
    c.Object.AddTags("character")
    c.Object.Data = c

    if !c.IsPlayer {
        InitNPCGOAP(c)
    }
    return c
}

func (c *Character) MoveToPoint(target resolv.Vector) bool {
    if c.Object.Center().Distance(target) <= c.Speed {
        c.Object.Position = target.Sub(c.Object.Center().Sub(c.Object.Position))
        c.Object.Update()
        return true
    }

    direction := target.Sub(c.Object.Center()).Unit()
    return c.Move(direction)
}

func (c *Character) Move(direction resolv.Vector) bool {
    direction = direction.Unit()

    step := direction.Mult(resolv.NewVector(c.Speed, c.Speed))
    if collision := c.Object.Check(step.X, step.Y, "mountain", "goblin_den"); collision == nil {
        c.Object.Position = c.Object.Position.Add(step)
        c.Object.Update()
        return true
    } else {
        return false
    }
}

func (c *Character) TakeDamage(amount int) {
    c.Health -= amount
    if c.Health < 0 {
        c.Health = 0
    }
}

func (c *Character) Update() {
    c.Attack.Update()

    if c.IsPlayer {
        // Player update logic (controlled by input)
        if c.Attack.IsAttacking && !c.Attack.HasDealtDamage {
            c.PerformAttack()
            c.Attack.HasDealtDamage = true
        }
    } else {
        // NPC update logic using GOAP
        currentState := c.UpdateGOAPState()
        goalState := c.GenerateGOAPGoal(currentState)

        fmt.Println(currentState, goalState)
        c.CurrentPlan = c.Planner.Plan(currentState, goalState)
        if c.CurrentPlan == nil {
            return
        }

        action := c.CurrentPlan[0]
        fmt.Println(c.Name, action)
        c.ExecuteGOAPAction(action)
        c.CurrentPlan = c.CurrentPlan[1:]
    }
}

func (c *Character) PerformAttack() {
    checkX := c.Object.Position.X - c.Attack.Range
    checkY := c.Object.Position.Y - c.Attack.Range
    checkSize := c.Attack.Range * 2
    nearbyObjects := c.Object.Space.CheckWorld(checkX, checkY, checkSize, checkSize,
        "monster", "goblin_den")

    seenObj := make(map[*resolv.Object]bool, 0)
    for _, obj := range nearbyObjects {
        if obj == c.Object || seenObj[obj] {
            continue
        }
        seenObj[obj] = true

        distance := c.Object.Center().Distance(obj.Center())
        if distance <= c.Attack.Range {
            switch {
            case obj.HasTags("monster"):
                if monster, ok := obj.Data.(*Monster); ok {
                    monster.TakeDamage(c.Attack.Damage)
                }
            case obj.HasTags("goblin_den"):
                if den, ok := obj.Data.(*GoblinDen); ok {
                    den.TakeDamage(c.Attack.Damage)
                }
            }
        }
    }
}

func (c *Character) Take() {
    collisions := c.Object.Check(0, 0, "mushroom")
    if collisions != nil {
        for _, obj := range collisions.Objects {
            switch {
            case obj.HasTags("mushroom"):
                c.Health = min(c.Health+20, 120)
                obj.Space.Remove(obj)
            }
        }
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}
