package units

import (
    "fmt"
    "github.com/solarlune/resolv"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNPCBehaviors(t *testing.T) {
    space := resolv.NewSpace(1000, 1000, 32, 32)
    npc := NewCharacter(100, 100, "TestNPC")
    space.Add(npc.Object)
    InitNPCGOAP(npc)

    t.Run("Test RunToSafety", func(t *testing.T) {
        initialPos := npc.Object.Position
        monster := NewMonster(120, 120, nil)
        space.Add(monster.Object)

        npc.RunToSafety()

        if npc.Object.Position == initialPos {
            t.Errorf("Expected NPC to move away from danger")
        }
    })

    t.Run("Test LookForMushroom", func(t *testing.T) {
        mushroom := NewMushroom(space, 96, 96)

        before := npc.Object.Center().Distance(mushroom.Object.Center())
        npc.LookForMushroom()
        after := npc.Object.Center().Distance(mushroom.Object.Center())

        if after > before {
            t.Errorf("Expected NPC to move towards mushroom before: %v after: %v ", before, after)
        }

    })

    t.Run("Test LookForMushroom", func(t *testing.T) {
        mushroom := NewMushroom(space, 128, 128)

        for i := 0; i < 30; i++ {
            before := npc.Object.Center().Distance(mushroom.Object.Center())
            npc.LookForMushroom()
            after := npc.Object.Center().Distance(mushroom.Object.Center())

            if npc.Object.Center() != mushroom.Object.Center() && after > before {
                t.Errorf("Iter %v Expected NPC to move towards mushroom before: %v after: %v ", i, before, after)
            }
        }

    })

}

func TestGOAP_NPCBehaviors(t *testing.T) {
    space := resolv.NewSpace(1000, 1000, 32, 32)
    npc := NewCharacter(100, 100, "TestNPC")
    space.Add(npc.Object)
    InitNPCGOAP(npc)

    t.Run("Test LookForMushroom", func(t *testing.T) {
        mushroom := NewMushroom(space, 120, 120)
        npc.Health = 80
        before := npc.Object.Center().Distance(mushroom.Object.Center())

        currentState := npc.UpdateGOAPState()
        goalState := npc.GenerateGOAPGoal(currentState)

        fmt.Println(currentState, goalState)
        npc.CurrentPlan = npc.Planner.Plan(currentState, goalState)
        if npc.CurrentPlan == nil {
            return
        }

        action := npc.CurrentPlan[0]
        fmt.Println(npc.Name, action)
        npc.ExecuteGOAPAction(action)

        after := npc.Object.Center().Distance(mushroom.Object.Center())

        assert.Equal(t, LookForMushroom, action.Name)
        assert.Truef(t, after < before, "Expected NPC to move towards mushroom before: %v after: %v ", before, after)
    })

}
