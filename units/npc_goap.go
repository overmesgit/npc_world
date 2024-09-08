package units

import (
    "example.com/maj/ai"
    "fmt"
    "github.com/solarlune/resolv"
    "math"
    "math/rand"
    "time"
)

const (
    RunToSafety     = "RunToSafety"
    LookForMushroom = "LookForMushroom"
    TakeMushroom    = "TakeMushroom"
    FindMonster     = "FindMonster"
    MoveToTarget    = "MoveToTarget"
    AttackMonster   = "AttackMonster"
    Wander          = "Wander"
)

func InitNPCGOAP(npc *Character) {
    npc.Planner = ai.NewGOAPPlanner()

    // Run to safety action
    npc.Planner.AddAction(ai.GOAPAction{
        Name: RunToSafety,
        CostFunc: func(state ai.GOAPState) float64 {
            return 1
        },
        Preconditions: ai.GOAPState{"lowHealth": true, "inDanger": true},
        Effects:       ai.GOAPState{"inDanger": false},
    })

    // Look for mushroom action
    npc.Planner.AddAction(ai.GOAPAction{
        Name: LookForMushroom,
        CostFunc: func(state ai.GOAPState) float64 {
            return 2
        },
        Preconditions: ai.GOAPState{"seeMushroom": true},
        Effects:       ai.GOAPState{"mushroomNear": true},
    })
    npc.Planner.AddAction(ai.GOAPAction{
        Name: TakeMushroom,
        CostFunc: func(state ai.GOAPState) float64 {
            return 2
        },
        Preconditions: ai.GOAPState{"mushroomNear": true},
        Effects:       ai.GOAPState{"hasFullHealth": true},
    })

    // Find monster action
    npc.Planner.AddAction(ai.GOAPAction{
        Name: FindMonster,
        CostFunc: func(state ai.GOAPState) float64 {
            return 3
        },
        Preconditions: ai.GOAPState{"hasTarget": false, "monstersArround": true},
        Effects:       ai.GOAPState{"hasTarget": true},
    })

    npc.Planner.AddAction(ai.GOAPAction{
        Name: MoveToTarget,
        CostFunc: func(state ai.GOAPState) float64 {
            return 3
        },
        Preconditions: ai.GOAPState{"hasTarget": true, "inAttackRange": false},
        Effects:       ai.GOAPState{"inAttackRange": true},
    })

    // Attack monster action
    npc.Planner.AddAction(ai.GOAPAction{
        Name: AttackMonster,
        CostFunc: func(state ai.GOAPState) float64 {
            return 4
        },
        Preconditions: ai.GOAPState{"hasTarget": true, "inAttackRange": true},
        Effects:       ai.GOAPState{"hasDefeatedMonster": true},
    })

    npc.Planner.AddAction(ai.GOAPAction{
        Name: Wander,
        CostFunc: func(state ai.GOAPState) float64 {
            return 6
        },
        Preconditions: ai.GOAPState{},
        Effects:       ai.GOAPState{"monstersArround": true, "seeMushroom": true},
    })

}

func (npc *Character) UpdateGOAPState() ai.GOAPState {
    if npc.TargetMonster != nil && npc.TargetMonster.Object.Space == nil {
        npc.TargetMonster = nil
    }

    state := ai.GOAPState{
        "lowHealth":       npc.Health < int(float32(npc.MaxHealth)*0.3),
        "hasFullHealth":   npc.Health == npc.MaxHealth,
        "inDanger":        npc.IsInDanger(),
        "hasTarget":       npc.HasTarget(),
        "inAttackRange":   npc.IsInAttackRange(),
        "monstersArround": npc.IsMonstersArround(),
        "mushroomNear":    npc.IsMushroomHere(),
        "seeMushroom":     npc.IsMushroomNear(),
    }
    return state
}

func (npc *Character) GenerateGOAPGoal(currentState ai.GOAPState) ai.GOAPState {
    if currentState["lowHealth"].(bool) && currentState["monstersArround"].(bool) {
        return ai.GOAPState{"inDanger": false}
    } else if currentState["inAttackRange"].(bool) {
        return ai.GOAPState{"hasDefeatedMonster": true}
    } else if npc.Health < npc.MaxHealth {
        return ai.GOAPState{"hasFullHealth": true}
    } else if currentState["hasTarget"].(bool) {
        return ai.GOAPState{"inAttackRange": true}
    } else if currentState["monstersArround"].(bool) {
        return ai.GOAPState{"hasTarget": true}
    } else {
        return ai.GOAPState{"monstersArround": true}
    }
}

func (npc *Character) ExecuteGOAPAction(action ai.GOAPAction) {
    switch action.Name {
    case "RunToSafety":
        npc.RunToSafety()
    case "LookForMushroom":
        npc.LookForMushroom()
    case "TakeMushroom":
        npc.Take()
    case "FindMonster":
        npc.FindMonster()
    case "MoveToTarget":
        npc.MoveTowards(npc.TargetMonster.Object.Center())
    case "AttackMonster":
        npc.AttackMonster()
    case "Wander":
        npc.Wander()
    }
}

func (npc *Character) IsMushroomHere() bool {
    _, distance := FindNearest(npc.Object, 32, "mushroom")
    fmt.Println("distance", distance)
    return distance < 16
}

func (npc *Character) IsMushroomNear() bool {
    nearbyMonsters := FindAll(npc.Object, 6*32, "mushroom")
    return len(nearbyMonsters) > 0
}

func (npc *Character) IsMonstersArround() bool {
    nearbyMonsters := FindAll(npc.Object, 6*32, "monster")
    return len(nearbyMonsters) > 0
}

func (npc *Character) IsInDanger() bool {
    nearbyMonsters := FindAll(npc.Object, 4*32, "monster")
    return len(nearbyMonsters) > 0
}

func (npc *Character) HasTarget() bool {
    // Check if the NPC has a target monster
    return npc.TargetMonster != nil
}

func (npc *Character) IsInAttackRange() bool {
    if npc.TargetMonster == nil {
        return false
    }
    distance := npc.Object.Center().Distance(npc.TargetMonster.Object.Center())
    return distance <= npc.Attack.Range
}

func (npc *Character) Wander() {
    canMove := false
    if npc.WanderTime.After(time.Now()) {
        canMove = npc.Move(npc.WanderTarget.X, npc.WanderTarget.Y)
    }

    for !canMove {
        angle := rand.Float64() * 2 * math.Pi
        direction := resolv.Vector{X: math.Cos(angle), Y: math.Sin(angle)}
        npc.WanderTime = time.Now().Add(time.Second * 5)
        npc.WanderTarget = direction
        canMove = npc.Move(direction.X, direction.Y)
    }

}

func (npc *Character) RunToSafety() {
    // Find the furthest point from all monsters and move towards it
    safePoint := FindSafePoint(npc.Object)
    npc.MoveTowards(safePoint)
}

func (npc *Character) LookForMushroom() {
    // Find the nearest mushroom and move towards it
    nearestMushroom, _ := FindNearest(npc.Object, 32*6, "mushroom")
    if nearestMushroom != nil {
        npc.MoveTowards(nearestMushroom.Center())
    }
}

func (npc *Character) FindMonster() {
    // Find the nearest monster and set it as the target
    nearestMonster, _ := FindNearest(npc.Object, 32*6, "monster")
    if nearestMonster != nil {
        npc.TargetMonster = nearestMonster.Data.(*Monster)
        npc.MoveTowards(nearestMonster.Center())
    }
}

func (npc *Character) AttackMonster() {
    if npc.TargetMonster != nil {
        npc.Attack.TriggerAttack()
        if npc.Attack.IsAttacking && !npc.Attack.HasDealtDamage {
            npc.TargetMonster.TakeDamage(npc.Attack.Damage)
            npc.Attack.HasDealtDamage = true
        }
    }
}

func (npc *Character) MoveTowards(target resolv.Vector) {
    direction := target.Sub(npc.Object.Center()).Unit()
    npc.Move(direction.X, direction.Y)
}
