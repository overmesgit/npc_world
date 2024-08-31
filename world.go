package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "github.com/solarlune/resolv"
    "math/rand"
    "time"
)

type World struct {
    characters       []Character
    monsters         []*Monster
    gameMap          *GameMap
    lastMonsterSpawn time.Time
    monsterSprite    *ebiten.Image
    maxMonsters      int
    space            *resolv.Space
}

func NewWorld(monsterSprite *ebiten.Image) *World {
    gameMap := NewGameMap()
    w := &World{
        characters:       make([]Character, 0),
        monsters:         make([]*Monster, 0),
        gameMap:          gameMap,
        lastMonsterSpawn: time.Now(),
        monsterSprite:    monsterSprite,
        maxMonsters:      10,
        space:            resolv.NewSpace(gameMap.Width*TileSize, gameMap.Height*TileSize, TileSize, TileSize),
    }
    w.initializeCollisionSpace()
    return w
}

func (w *World) initializeCollisionSpace() {
    for y := 0; y < w.gameMap.Height; y++ {
        for x := 0; x < w.gameMap.Width; x++ {
            if w.gameMap.Tiles[y][x] == TileMountain {
                obj := resolv.NewObject(float64(x*TileSize), float64(y*TileSize), float64(TileSize), float64(TileSize))
                obj.SetShape(resolv.NewRectangle(0, 0, float64(TileSize), float64(TileSize)))
                obj.AddTags("mountain")
                w.space.Add(obj)
            }
        }
    }
}

func (w *World) Update() {
    for i := range w.characters {
        w.characters[i].Update(w)
    }

    // Update monsters and remove dead ones
    aliveMonsters := make([]*Monster, 0)
    for _, monster := range w.monsters {
        monster.Update(w)
        if monster.Health > 0 {
            aliveMonsters = append(aliveMonsters, monster)
        }
    }
    w.monsters = aliveMonsters

    // Spawn new monster every 10 seconds
    if time.Since(w.lastMonsterSpawn) > 3*time.Second {
        w.SpawnMonster()
        w.lastMonsterSpawn = time.Now()
    }
}

func (w *World) SpawnMonster() bool {
    // Define the central area
    centerX := 18
    centerY := 18
    spawnRadius := 4 // Adjust this value to change the spawn area size

    // Try to find a suitable spawn point
    for attempts := 0; attempts < 100; attempts++ {
        xTile := centerX + (rand.Intn(spawnRadius*2+1) - spawnRadius)
        yTile := centerY + (rand.Intn(spawnRadius*2+1) - spawnRadius)

        if w.IsSpawnPointValid(xTile, yTile) {
            monster := NewMonster(float64(xTile*TileSize), float64(yTile*TileSize), w.monsterSprite)
            monster.NormalizeDirection()
            w.monsters = append(w.monsters, monster)
            w.space.Add(monster.Object)
            return true
        }
    }
    return false
}

func (w *World) IsSpawnPointValid(xTile, yTile int) bool {
    if xTile < 0 || xTile >= w.gameMap.Width || yTile < 0 || yTile >= w.gameMap.Height {
        return false
    }
    collision := w.space.CheckCells(xTile, yTile, 1, 1, "mountain", "character", "monster")
    if len(collision) == 0 {
        return true
    }
    return false
}

func (w *World) AddCharacter(c Character) {
    w.characters = append(w.characters, c)
    w.space.Add(c.Object)
}

func (w *World) GetCharacters() []Character {
    return w.characters
}

func (w *World) GetPlayerCharacter() *Character {
    // Assuming the first character is always the player
    if len(w.characters) > 0 {
        return &w.characters[0]
    }
    return nil
}
