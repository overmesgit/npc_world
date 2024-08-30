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
    if time.Since(w.lastMonsterSpawn) > 3*time.Second {
        w.SpawnMonster()
        w.lastMonsterSpawn = time.Now()
    }
}

func (w *World) SpawnMonster() bool {
	// Define the central area
	centerX := w.gameMap.Width
	centerY := w.gameMap.Height
	spawnRadius := 10 // Adjust this value to change the spawn area size

	// Try to find a suitable spawn point
	for attempts := 0; attempts < 100; attempts++ {
		x := centerX + (rand.Intn(spawnRadius*2+1) - spawnRadius)
		y := centerY + (rand.Intn(spawnRadius*2+1) - spawnRadius)

		if w.IsSpawnPointValid(x, y) {
			monster := NewMonster(float64(x*TileSize), float64(y*TileSize), w.monsterSprite)
			monster.NormalizeDirection()
			w.monsters = append(w.monsters, monster)
			return true
		}
	}
	return false
}

func (w *World) IsSpawnPointValid(x, y int) bool {
	// Check if the tile is walkable
	if !w.gameMap.IsTileWalkable(x, y) {
		return false
	}

	// Check for collision with other monsters
	for _, monster := range w.monsters {
		monsterTileX := int(monster.X / TileSize)
		monsterTileY := int(monster.Y / TileSize)
		if monsterTileX == x && monsterTileY == y {
			return false
		}
	}

	// Check for collision with characters
	for _, character := range w.characters {
		charTileX := int(character.X / TileSize)
		charTileY := int(character.Y / TileSize)
		if charTileX == x && charTileY == y {
			return false
		}
	}

	return true
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

func (w *World) IsTileWalkable(x, y int) bool {
    if x < 0 || x >= w.gameMap.Width || y < 0 || y >= w.gameMap.Height {
        return false
    }
    return w.gameMap.IsTileWalkable(x, y)
}
