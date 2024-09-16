package units

import (
    "example.com/maj/ai"
    gamemap "example.com/maj/map"
    "example.com/maj/pathfinding"
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
    MoveToDen       = "MoveToDen"
    AttackDen       = "AttackDen"
)

func InitNPCGOAP(npc *Character) {
    npc.Planner = ai.NewGOAPPlanner()

    // Run to safety action
    npc.Planner.AddAction(ai.GOAPAction{
        Name: RunToSafety,
        CostFunc: func(state ai.GOAPState) float64 {
            return 1
        },
        Preconditions: ai.GOAPState{"lowHealth": true, "monstersArround": true},
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
        Name: MoveToDen,
        CostFunc: func(state ai.GOAPState) float64 {
            return 4
        },
        Preconditions: ai.GOAPState{"seeGoblinDen": true},
        Effects:       ai.GOAPState{"denInAttackRange": true},
    })

    npc.Planner.AddAction(ai.GOAPAction{
        Name: AttackDen,
        CostFunc: func(state ai.GOAPState) float64 {
            return 4
        },
        Preconditions: ai.GOAPState{"denInAttackRange": true},
        Effects:       ai.GOAPState{"hasDefeatedDen": true},
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
        "lowHealth":        npc.Health < int(float32(npc.MaxHealth)*0.3),
        "hasFullHealth":    npc.Health == npc.MaxHealth,
        "inDanger":         npc.IsInDanger(),
        "hasTarget":        npc.HasTarget(),
        "inAttackRange":    npc.IsInAttackRange(),
        "monstersArround":  npc.IsMonstersArround(),
        "mushroomNear":     npc.IsMushroomHere(),
        "seeMushroom":      npc.IsMushroomNear(),
        "denInAttackRange": npc.DenInAttackRange(),
        "seeGoblinDen":     npc.seeGoblinDen(),
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
    } else if currentState["denInAttackRange"].(bool) {
        return ai.GOAPState{"hasDefeatedDen": true}
    } else if currentState["seeGoblinDen"].(bool) {
        return ai.GOAPState{"denInAttackRange": true}
    } else {
        return ai.GOAPState{"monstersArround": true}
    }
}

func (npc *Character) ExecuteGOAPAction(action ai.GOAPAction) {
    switch action.Name {
    case RunToSafety:
        npc.RunToSafety()
    case LookForMushroom:
        npc.LookForMushroom()
    case TakeMushroom:
        npc.Take()
    case FindMonster:
        npc.FindMonster()
    case MoveToTarget:
        npc.MoveTowards(npc.TargetMonster.Object.Center())
    case AttackMonster:
        npc.AttackMonster()
    case Wander:
        npc.Wander()
    case MoveToDen:
        npc.MoveTowardsDen()
    case AttackDen:
        npc.AttackDen()
    }
}

func (npc *Character) IsMushroomHere() bool {
    _, distance := FindNearest(npc.Object, 32, "mushroom")
    return distance < 16
}

func (npc *Character) IsMushroomNear() bool {
    nearbyMonsters := FindAll(npc.Object, 6*32, "mushroom")
    return len(nearbyMonsters) > 0
}

func (npc *Character) seeGoblinDen() bool {
    nearbyMonsters := FindAll(npc.Object, 6*32, "goblin_den")
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

func (npc *Character) DenInAttackRange() bool {
    _, distance := FindNearest(npc.Object, npc.Attack.Range, "goblin_den")
    return distance <= npc.Attack.Range
}

func (npc *Character) Wander() {
    canMove := false
    if npc.WanderTime.After(time.Now()) {
        canMove = npc.Move(npc.WanderTarget)
    }

    for !canMove {
        angle := rand.Float64() * 2 * math.Pi
        direction := resolv.Vector{X: math.Cos(angle), Y: math.Sin(angle)}
        npc.WanderTime = time.Now().Add(time.Second * 5)
        npc.WanderTarget = direction
        canMove = npc.Move(direction)
    }

}

func (npc *Character) RunToSafety() {
    // Find the furthest point from all monsters and move towards it
    npc.TargetMonster = nil
    safePoint := FindSafePoint(npc.Object)
    if safePoint == npc.Object.Center() {
        npc.Wander()
    } else {
        npc.MoveTowards(safePoint)
    }
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

func (npc *Character) AttackDen() {
    denObj, distance := FindNearest(npc.Object, npc.Attack.Range, "goblin_den")
    if denObj != nil && distance <= npc.Attack.Range {
        npc.Attack.TriggerAttack()
        if npc.Attack.IsAttacking && !npc.Attack.HasDealtDamage {
            denObj.Data.(*GoblinDen).TakeDamage(npc.Attack.Damage)
            npc.Attack.HasDealtDamage = true
        }
    }
}

func (npc *Character) MoveTowards(target resolv.Vector) {
    npcCenter := npc.Object.Center()
    startX, startY := npc.Object.Space.WorldToSpace(npcCenter.X, npcCenter.Y)

    endX, endY := npc.Object.Space.WorldToSpace(target.X, target.Y)

    path, _, _ := pathfinding.FindPath(npc.Object.Space, startX, startY, endX, endY)

    if len(path) == 0 {
        return
    }

    lastStep := path[len(path)-1].(pathfinding.PathNode)
    lastStepWorld := npc.Object.Space.SpaceToWorldVec(lastStep.X, lastStep.Y)
    nextTarget := lastStepWorld
    if len(path) > 1 {
        secondLastStep := path[len(path)-2].(pathfinding.PathNode)
        secondLastStepWorld := npc.Object.Space.SpaceToWorldVec(secondLastStep.X, secondLastStep.Y)

        // Check if NPC has surpassed the last target
        if npc.Object.Position.Distance(secondLastStepWorld) <= 32 {
            nextTarget = secondLastStepWorld
        }

    }

    halfCell := resolv.Vector{X: float64(gamemap.TileSize / 2), Y: float64(gamemap.TileSize / 2)}
    nextTarget = nextTarget.Add(halfCell)
    npc.MoveToPoint(nextTarget)
}

func (npc *Character) MoveTowardsDen() {
    denObj, _ := FindNearest(npc.Object, 6*32, "goblin_den")
    if denObj != nil {
        direction := denObj.Center().Sub(npc.Object.Center()).Unit()
        npc.Move(direction)
    }
}
