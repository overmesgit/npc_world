package game

import (
    "example.com/maj/units"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/solarlune/resolv"
    "math/rand"
    "time"
)

// Add new struct
type Mushroom struct {
    Object *resolv.Object
}

type World struct {
    GameMap         *GameMap
    MonsterSprite   *ebiten.Image
    GoblinDenSprite *ebiten.Image
    Space           *resolv.Space
    Player          *units.Character
}

func NewWorld() *World {
    gameMap := NewGameMap()
    w := &World{
        GameMap: gameMap,
        Space:   resolv.NewSpace(gameMap.Width*TileSize, gameMap.Height*TileSize, TileSize, TileSize),
    }
    w.initializeCollisionSpace()
    w.spawnGoblinDens(10)
    w.spawnMushrooms(30)
    go w.mushroomSpawnRoutine()
    return w
}

func (w *World) Update() {
    for _, obj := range w.Space.Objects() {
        switch obj.Data.(type) {
        case *units.Character:
            character := obj.Data.(*units.Character)
            character.Update()
        case *units.Monster:
            monster := obj.Data.(*units.Monster)
            monster.Update()
            if monster.Health <= 0 {
                w.Space.Remove(monster.Object)
            }
        case *units.GoblinDen:
            den := obj.Data.(*units.GoblinDen)
            newMonsters := den.Update()
            for _, m := range newMonsters {
                w.Space.Add(m.Object)
            }
            if den.Health <= 0 {
                w.Space.Remove(den.Object)
            }
        }
    }
}

func (w *World) spawnMushrooms(count int) {
    for i := 0; i < count; i++ {
        x, y := w.findValidSpawnPoint()
        mushroom := &Mushroom{
            Object: resolv.NewObject(float64(x*TileSize), float64(y*TileSize), float64(TileSize), float64(TileSize)),
        }
        mushroom.Object.AddTags("mushroom")
        mushroom.Object.Data = mushroom
        w.Space.Add(mushroom.Object)
    }
}

func (w *World) mushroomSpawnRoutine() {
    ticker := time.NewTicker(3 * time.Second)
    for range ticker.C {
        w.spawnMushrooms(1)
    }
}

func (w *World) spawnGoblinDens(count int) {
    for i := 0; i < count; i++ {
        x, y := w.findValidSpawnPoint()
        den := units.NewGoblinDen(float64(x*TileSize), float64(y*TileSize))
        w.Space.Add(den.Object)
    }
}

func (w *World) findValidSpawnPoint() (int, int) {
    for {
        x := rand.Intn(w.GameMap.Width)
        y := rand.Intn(w.GameMap.Height)
        if w.IsSpawnPointValid(x, y) {
            return x, y
        }
    }
}

func (w *World) initializeCollisionSpace() {
    for y := 0; y < w.GameMap.Height; y++ {
        for x := 0; x < w.GameMap.Width; x++ {
            if w.GameMap.Tiles[y][x] == TileMountain {
                obj := resolv.NewObject(float64(x*TileSize), float64(y*TileSize), float64(TileSize), float64(TileSize))
                obj.SetShape(resolv.NewRectangle(0, 0, float64(TileSize), float64(TileSize)))
                obj.AddTags("mountain")
                w.Space.Add(obj)
            }
        }
    }
}

func (w *World) IsSpawnPointValid(xTile, yTile int) bool {
    if xTile < 0 || xTile >= w.GameMap.Width || yTile < 0 || yTile >= w.GameMap.Height {
        return false
    }
    collision := w.Space.CheckCells(xTile, yTile, 1, 1, "mountain", "character", "monster")
    if len(collision) == 0 {
        return true
    }
    return false
}

func (w *World) AddCharacter(c *units.Character) {
    if w.Player == nil {
        w.Player = c
    }
    w.Space.Add(c.Object)
}

func (w *World) GetPlayerCharacter() *units.Character {
    return w.Player
}
