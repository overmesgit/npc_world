package game

import (
    "example.com/maj/units"
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/solarlune/resolv"
    "math/rand"
)

type World struct {
    Characters      []*units.Character
    Monsters        []*units.Monster
    GoblinDens      []*units.GoblinDen
    GameMap         *GameMap
    MonsterSprite   *ebiten.Image
    GoblinDenSprite *ebiten.Image
    Space           *resolv.Space
}

func NewWorld(monsterSprite, goblinDenSprite *ebiten.Image) *World {
    gameMap := NewGameMap()
    w := &World{
        Characters:      make([]*units.Character, 0),
        Monsters:        make([]*units.Monster, 0),
        GoblinDens:      make([]*units.GoblinDen, 0),
        GameMap:         gameMap,
        MonsterSprite:   monsterSprite,
        GoblinDenSprite: goblinDenSprite,
        Space:           resolv.NewSpace(gameMap.Width*TileSize, gameMap.Height*TileSize, TileSize, TileSize),
    }
    w.initializeCollisionSpace()
    w.spawnGoblinDens(3) // Spawn 3 goblin dens
    return w
}

func (w *World) Update() {
    for i := range w.Characters {
        w.Characters[i].Update(w.Monsters)
    }

    for _, den := range w.GoblinDens {
        newMonsters := den.Update()
        w.Monsters = append(w.Monsters, newMonsters...)
    }

    // Update Monsters and remove dead ones
    aliveMonsters := make([]*units.Monster, 0)
    for _, monster := range w.Monsters {
        monster.Update(w.Characters)
        if monster.Health > 0 {
            aliveMonsters = append(aliveMonsters, monster)
        }
    }
    w.Monsters = aliveMonsters
}

func (w *World) spawnGoblinDens(count int) {
    for i := 0; i < count; i++ {
        x, y := w.findValidSpawnPoint()
        den := units.NewGoblinDen(float64(x*TileSize), float64(y*TileSize))
        w.GoblinDens = append(w.GoblinDens, den)
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
    w.Characters = append(w.Characters, c)
    w.Space.Add(c.Object)
}

func (w *World) GetCharacters() []*units.Character {
    return w.Characters
}

func (w *World) GetPlayerCharacter() *units.Character {
    // Assuming the first character is always the player
    if len(w.Characters) > 0 {
        return w.Characters[0]
    }
    return nil
}
