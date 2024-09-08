package units

import (
    gamemap "example.com/maj/map"
    "github.com/solarlune/resolv"
)

type Mushroom struct {
    Object *resolv.Object
}

func NewMushroom(space *resolv.Space, x, y float64) *Mushroom {
    mushroom := &Mushroom{
        Object: resolv.NewObject(x, y, float64(gamemap.TileSize), float64(gamemap.TileSize)),
    }
    mushroom.Object.AddTags("mushroom")
    mushroom.Object.Data = mushroom
    space.Add(mushroom.Object)
    return mushroom
}
