package units

import (
    gamemap "example.com/maj/map"
    "github.com/solarlune/resolv"
    "math"
)

func FindAll(source *resolv.Object, distance float64, tags ...string) []*resolv.Object {
    var nearestObjects []*resolv.Object

    checkX := source.Center().X - distance/2
    checkY := source.Center().Y - distance/2
    checkSize := distance * 2
    nearbyObjects := source.Space.CheckWorld(
        checkX, checkY,
        checkSize, checkSize, tags...)

    for _, obj := range nearbyObjects {
        if obj == source {
            continue
        }
        nearestObjects = append(nearestObjects, obj)
    }

    return nearestObjects
}

func FindNearest(source *resolv.Object, distance float64, tags ...string) (*resolv.Object, float64) {
    var nearestChar *resolv.Object
    minDistance := math.Inf(1)

    nearbyObjects := FindAll(source, distance, tags...)

    for _, obj := range nearbyObjects {
        distance := obj.Center().Distance(source.Center())
        if distance < minDistance {
            minDistance = distance
            nearestChar = obj
        }
    }

    return nearestChar, minDistance
}

func FindSafePoint(source *resolv.Object) resolv.Vector {
    monsters := FindAll(source, 6*32, "monster")
    if len(monsters) == 0 {
        return source.Center()
    }

    // Find the average position of all monsters
    var avgX, avgY float64
    for _, monster := range monsters {
        avgX += monster.Position.X
        avgY += monster.Position.Y
    }
    avgX /= float64(len(monsters))
    avgY /= float64(len(monsters))

    // Move in the opposite direction of the average monster position
    safeDirection := source.Center().Sub(resolv.NewVector(avgX, avgY)).Unit()
    safePoint := source.Center().Add(safeDirection.Scale(5 * 32)) // Move 5 tiles away

    // Ensure the safe point is within the world bounds
    //    safePoint.X = math.Max(0, math.Min(safePoint.X, float64(w.GameMap.Width*32)))
    //    safePoint.Y = math.Max(0, math.Min(safePoint.Y, float64(w.GameMap.Height*32)))

    return safePoint
}

func CheckWorld(space *resolv.Space, x, y float64, w, h float64) []*resolv.Object {
    if int(x)/gamemap.TileSize != int(x+w)/gamemap.TileSize && int(x+w)%gamemap.TileSize != 0 {
        w += gamemap.TileSize
    }
    if int(y)/gamemap.TileSize != int(y+h)/gamemap.TileSize && int(y+h)%gamemap.TileSize != 0 {
        h += gamemap.TileSize
    }
    return space.CheckWorld(x, y, w, h)
}
