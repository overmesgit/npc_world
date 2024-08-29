package main

type World struct {
    characters []Character
    monsters   []Monster
    gameMap    *GameMap
}

func NewWorld() *World {
    return &World{
        characters: make([]Character, 0),
        monsters:   make([]Monster, 0),
        gameMap:    NewGameMap(),
    }
}

func (w *World) Update() {
    for i := range w.characters {
        w.characters[i].Update(w)
    }
    for i := range w.monsters {
        w.monsters[i].Update(w)
    }
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
