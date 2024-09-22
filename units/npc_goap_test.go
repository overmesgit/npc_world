package units

import (
    gamemap "example.com/maj/map"
    "fmt"
    "github.com/solarlune/resolv"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestNPCBehaviors(t *testing.T) {

    t.Run("Test RunToSafety", func(t *testing.T) {
        space, npc := InitSpace(100, 100)
        initialPos := npc.Object.Position
        monster := NewMonster(120, 120, nil)
        space.Add(monster.Object)

        npc.RunToSafety()

        assert.NotEqual(t, npc.Object.Position, initialPos, "Expected NPC to move away from danger")
    })

    t.Run("Test LookForMushroom Single step", func(t *testing.T) {
        space, npc := InitSpace(100, 100)
        mushroom := NewMushroom(space, 96, 96)

        before := npc.Object.Center().Distance(mushroom.Object.Center())
        npc.LookForMushroom()
        after := npc.Object.Center().Distance(mushroom.Object.Center())

        if after > before {
            t.Errorf("Expected NPC to move towards mushroom before: %v after: %v ", before, after)
        }

    })

    t.Run("Test LookForMushroom Multiple steps", func(t *testing.T) {
        space, npc := InitSpace(96, 96)
        mushroom := NewMushroom(space, 128, 128)

        for i := 0; i < 30; i++ {
            before := npc.Object.Center().Distance(mushroom.Object.Center())
            npc.LookForMushroom()
            after := npc.Object.Center().Distance(mushroom.Object.Center())

            fmt.Println(npc.Object.Center(), mushroom.Object.Center())
            if npc.Object.Center() == mushroom.Object.Center() {
                break
            }
            assert.Truef(t, after < before, "Iter %v Expected NPC to move towards mushroom before: %v after: %v ", i, before, after)
        }
        assert.Equal(t, npc.Object.Center(), mushroom.Object.Center(), "NPS didn't come to the mushroom")

    })

    t.Run("Test Go arround obsticles", func(t *testing.T) {
        // ##.##
        // ###M#
        // ####G
        // 2, 2 cell
        space, npc := InitSpace(112, 64)
        // 4, 4 cell
        mushroom := NewMushroom(space, 128, 128)
        // 3, 3 cell
        NewMountain(space, 96, 96)

        res := CheckWorld(space, 100, 100, 32, 32)
        assert.True(t, len(res) > 0, "Can't find mountain as obsticle")
        for i := 0; i < 70; i++ {
            before := npc.Object.Center().Distance(mushroom.Object.Center())
            npc.LookForMushroom()
            after := npc.Object.Center().Distance(mushroom.Object.Center())

            fmt.Println(npc.Object.Center(), mushroom.Object.Center())
            if npc.Object.Center() == mushroom.Object.Center() {
                break
            }
            assert.Truef(t, after < before, "Iter %v Expected NPC to move towards mushroom before: %v after: %v ", i, before, after)
        }
        assert.Equal(t, npc.Object.Center(), mushroom.Object.Center(), "NPS didn't come to the mushroom")

    })

}

func NewMountain(space *resolv.Space, x, y float64) {
    size := float64(gamemap.TileSize)
    obj := resolv.NewObject(x, y, size, size)
    obj.SetShape(resolv.NewRectangle(0, 0, size, size))
    obj.AddTags("mountain")
    space.Add(obj)
}

func InitSpace(x, y float64) (*resolv.Space, *Character) {
    space := resolv.NewSpace(1000, 1000, 32, 32)
    npc := NewCharacter(x, y, "TestNPC")
    space.Add(npc.Object)
    InitNPCGOAP(npc)
    return space, npc
}

func Test_GOAPBehaviors(t *testing.T) {
    t.Run("Test Plan LookForMushroom", func(t *testing.T) {
        space, npc := InitSpace(100, 100)
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

    t.Run("Test run away from enemy", func(t *testing.T) {
        space, npc := InitSpace(100, 160)
        monster := NewMonster(100, 100, nil)
        space.Add(monster.Object)

        npc.Health = 20

        for i := range 80 {
            before := npc.Object.Center().Distance(monster.Object.Center())

            currentState := npc.UpdateGOAPState()
            goalState := npc.GenerateGOAPGoal(currentState)

            fmt.Println(currentState, goalState)
            npc.CurrentPlan = npc.Planner.Plan(currentState, goalState)
            assert.NotNil(t, npc.CurrentPlan, "Current plan is nil")

            action := npc.CurrentPlan[0]
            fmt.Println(npc.Name, action, npc.Object.Center())
            npc.ExecuteGOAPAction(action)

            after := npc.Object.Center().Distance(monster.Object.Center())

            assert.Equal(t, RunToSafety, action.Name)
            assert.Truef(t, after > before, "Iter %v Expected NPC to run from monster: %v after: %v ", i, before, after)
            if after >= gamemap.TileSize*3 {
                break
            }
        }
    })

}
