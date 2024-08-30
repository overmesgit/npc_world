package main

import (
    "github.com/hajimehoshi/ebiten/v2"
    "math/rand"
    "time"
)

type World struct {
    characters       []Character
    monsters         []*Monster
    gameMap          *GameMap
    lastMonsterSpawn time.Time
    monsterSprite    *ebiten.Image
}

func NewWorld(monsterSprite *ebiten.Image) *World {
    return &World{
        characters:       make([]Character, 0),
        monsters:         make([]*Monster, 0),
        gameMap:          NewGameMap(),
        lastMonsterSpawn: time.Now(),
        monsterSprite:    monsterSprite,
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
    if time.Since(w.lastMonsterSpawn) > 10*time.Second {
        w.SpawnMonster()
        w.lastMonsterSpawn = time.Now()
    }
}

func (w *World) SpawnMonster() {
    centerX := float64((w.gameMap.Width / 2) * TileSize)
    centerY := float64((w.gameMap.Height / 2) * TileSize)

    // Ensure the spawn point is not on a mountain
    for i := 0; i < 100; i++ { // Limit attempts to prevent infinite loop
        x := centerX + (rand.Float64()*10-5)*TileSize
        y := centerY + (rand.Float64()*10-5)*TileSize
        if w.gameMap.IsTileWalkable(int(x/TileSize), int(y/TileSize)) {
            monster := NewMonster(x, y, w.monsterSprite)
            monster.NormalizeDirection() // Ensure the initial direction is normalized
            w.monsters = append(w.monsters, monster)
            return
        }
    }
    // If no suitable spawn point found, don't spawn a monster this time
}

func (w *World) AddCharacter(c Character) {
    w.characters = append(w.characters, c)
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
