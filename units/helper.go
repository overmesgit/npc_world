package units

import (
    "github.com/solarlune/resolv"
    "math"
)

func FindAll(source *resolv.Object, distance float64, tags ...string) []*resolv.Object {
    var nearestObjects []*resolv.Object

    checkX := source.Center().X - distance
    checkY := source.Center().Y - distance
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
