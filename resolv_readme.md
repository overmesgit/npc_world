![](https://i.imgur.com/BDQ2wWJ.gif)

# Resolv v0.6.0

[pkg.go.dev](https://pkg.go.dev/github.com/solarlune/resolv)

## What is Resolv?

Resolv is a 2D collision detection and resolution library, specifically created for simpler, arcade-y (non-realistic) video games. Resolv is written in Go, but the core concepts are fairly straightforward and could be easily adapted for use with other languages or game development frameworks.

Basically: It allows you to do simple physics easier, without actually doing the physics part - that's still on you.

## Why is it called that?

Because it's like... You know, collision resolution? To **resolve** a collision? So... That's the name. I juste took an e off because I seem to have misplaced it.

## Why did you create Resolv?

Because I was making games in Go and found that existing frameworks tend to omit collision testing and resolution code. Collision testing isn't too hard, but it's done frequently enough, and most games need simple enough physics that it makes sense to make a library to handle collision testing and resolution for simple, "arcade-y" games; if you need realistic physics, you have other options like Box2D or something.

____

As an aside, this actually used to be quite different; I decided to rework it when I was less than satisfied with my previous efforts, and after a few attempts over several months, I got it to this state (which, I think, is largely better). That said, there are breaking changes between the previous version, v0.4, and the current one (v0.5). These changes were necessary, in my opinion, to improve the library.

In comparison to the previous version of Resolv, v0.5 includes, among other things:

- A redesigned API from scratch. 
- The usage of floats instead of ints for position and movement, simplifying real usage of the library dramatically.
- Broadphase grid-based collision testing and querying for simple collisions, which means solid performance gains.
- ConvexPolygons for SAT intersection tests.

It's still a work-in-progress, but should be solid enough for usage in the field.

## How do I get it?

`go get github.com/solarlune/resolv`

## How do I use it?

There's a couple of ways to use Resolv. 

Firstly, you can create a Space, create Objects, add them to the Space, check the Space for collisions / intersections, and finally update the Objects as they move.

In Resolv v0.5, a Space represents a limited, bounded area which is separated into even Cells of predetermined size. Any Objects added to the Space fill at least one Cell (as long as the Object is within the Space). By checking a position in the Space, you can tell which Cells are occupied and so, where objects generally are. This is the broadphase, simpler portion of Resolv. 

Here's an example:

```go

var space *resolv.Space
var playerObj *resolv.Object

// As an example, in a game's initialization function that runs once
// when the game or level starts...

func Init() {

    // First, we want to create a Space. This represents the areas in our game world 
    // where objects can move about and check for collisions.

    // The first two arguments represent our Space's width and height, while the 
    // following two represent the individual Cells' sizes. The smaller the Cells' 
    // sizes, the finer the collision detection. Generally, each Cell should be 
    // reasonably be the size of a "unit", whatever that may be for the game.
    // For example, the player character, enemies, and collectibles could fit 
    // into one or more of these Cells.

    space = resolv.NewSpace(640, 480, 16, 16)

    // Next, we can start creating things and adding it to the Space.

    // Here's some level geometry. resolv.NewObject() takes the X and Y
    // position, and width and height to create a new *resolv.Object.
    // You can also specify tags when creating an Object. Tags can be used 
    // to filter down objects when checking the Space for a collision.

    space.Add(
        resolv.NewObject(0, 0, 640, 16),
        resolv.NewObject(0, 480-16, 640, 16),
        resolv.NewObject(0, 16, 16, 480-32),
        resolv.NewObject(640-16, 16, 16, 480-32),
    )

    // We'll keep a reference to the player's Object to move it later.
    playerObj = resolv.NewObject(32, 32, 16, 16)

    // Finally, we add the Object to the Space, and we're good to go!
    space.Add(playerObj)

}

// Later on, in the game's update loop, which runs once per game frame...

func Update() {

    // Let's say we are attempting to move the player to the right by 2 
    // pixels. Here's how we could do it.
    dx := 2.0

    // To start, we check to see if there would be a collision if the 
    // playerObj were to move to the right by 2 pixels. The Check function 
    // returns a Collision object if so.

    if collision := playerObj.Check(dx, 0); collision != nil {
        
        // If there was a collision, the "playerObj" Object can't move fully 
        // to the right by 2, and Object.Check() would return a *Collision object.
        // A *Collision object contains the Objects and Cells that the calling 
        // *resolv.Object ran into when it called Check().

        // To resolve (haha) this collision, we probably want to move the player into
        // contact with that Object. So, we call Collision.ContactWithObject() on the 
        // first Object that we came into contact with (which is stored in the Collision).

        // Collision.ContactWithObject() will return a Vector, indicating how much
        // distance to move to come into contact with the specified Object.

        // We could also come into contact with the cell to the right using 
        // Collision.ContactWithCell(collision.Cells[0]).
        dx = collision.ContactWithObject(collision.Objects[0]).X()

    }

    // If there wasn't a collision, then dx will just be 2, as set above, and the 
    // movement will go through unimpeded.
    playerObj.X += dx

    // Lastly, when we move an Object, we need to call Object.Update() so it can be 
    // updated within the Space as well. For static / unmoving Objects, this is
    // unnecessary, as Object.Update() is called once when an Object is first added to a Space.
    playerObj.Update()

    // If we were making a platformer, you could then check for the Y-axis as well. 
    // Conceptually, this is decomposing movement into two separate axes, and is a familiar 
    // and well-used approach for handling movement in a standard tile-based platformer. 
    // See this fantastic post on the subject:
    // http://higherorderfun.com/blog/2012/05/20/the-guide-to-implementing-2d-platformers/

    // If you want to filter out types of Objects to check for, add tags on the objects 
    // you want to filter using Object.AddTags(), or when the Object is created 
    // with resolv.NewObject(), and specify them in Object.Check().

    onlySolidOrHazardous := playerObj.Check(dx, 0, "hazard", "solid")

}

// That's it!

```

The second way to use Resolv is to check for a more accurate shape intersection test by assigning two Objects Shapes, and then checking for an intersection between them. Checking for an intersection between Shapes internally performs separating axis theorum (SAT) collision testing (when checking against ConvexPolygons), and represents the more inefficient narrow-phase portion of Resolv. If you can get by without doing Shape-based intersection testing, it would be most performant to do so.

```go

playerObj *resolv.Object
stairs *resolv.Object
space *resolv.Space

func Init() {
    
    space = resolv.NewSpace(640, 480, 16, 16)

    // Create the Object as usual, but then...
    playerObj = resolv.NewObject(32, 128, 16, 16)
    // Assign the Object a Shape. A Rectangle is, for now, a ConvexPolygon that's simply 
    // rectangular, rather than a specific, separate Shape.
    playerObj.SetShape(resolv.NewRectangle(0, 0, 16, 16))
    // Then we add the Object to the Space.
    space.Add(playerObj)

    // Note that we can just use the shapes directly as well.

    stairs = resolv.NewObject(96, 128, 16, 16)

    // Here, we use resolv.NewConvexPolygon() to create a new ConvexPolygon Shape. It takes 
    // a series of float64 values indicating the X and Y positions of each vertex; the call 
    // below, for example, creates a triangle.

    stairs.SetShape(resolv.NewConvexPolygon(
        0, 0, // Position of the polygon

        16, 0, // (x, y) pair for the first vertex
        16, 16, // (x, y) pair for the second vertex
        0, 16, // (x, y) pair for the third and last vertex
    ))

    //     0
    //    /|
    //   / |
    //  /  |
    // 2---1

    // Note that the vertices are in clockwise order. They can be in either clockwise or 
    // counter-clockwise order as long as it's consistent throughout your application. 
    // As an aside, resolv.NewRectangle() defines the vertices in clockwise order.
    space.Add(stairs)

}

func Update() {

    dx := 1.0

    // Shape.Intersection() returns a ContactSet, representing information 
    // regarding the intersection between two Shapes (i.e. the point(s) of
    // collision, the distance to move to get out, etc).
    if intersection := playerObj.Shape.Intersection(dx, 0, stairs.Shape); intersection != nil {
        
        // We are colliding with the stairs shape, so we can move according
        // to the delta (MTV) to get out of it.
        dx = intersection.MTV.X()

        // You might want to move a bit less (say, 0.1) than the delta to
        // avoid "bouncing", depending on your application.

    }

    playerObj.X += dx

    // When Object.Update() is called, the Object's Shape is also moved
    // accordingly.
    playerObj.Update()

}

```

Index ¶
Variables
func ToDegrees(radians float64) float64
func ToRadians(degrees float64) float64
type Cell
func (cell *Cell) Contains(obj *Object) bool
func (cell *Cell) ContainsTags(tags ...string) bool
func (cell *Cell) Occupied() bool
type Circle
func NewCircle(x, y, radius float64) *Circle
func (circle *Circle) Bounds() (Vector, Vector)
func (circle *Circle) Clone() IShape
func (circle *Circle) Intersection(dx, dy float64, other IShape) *ContactSet
func (c *Circle) IntersectionForEach(dx, dy float64, f func(c *ContactSet) bool, others ...IShape)
func (circle *Circle) IntersectionPointsCircle(other *Circle) []Vector
func (circle *Circle) Move(x, y float64)
func (circle *Circle) MoveVec(vec Vector)
func (circle *Circle) PointInside(point Vector) bool
func (circle *Circle) Position() Vector
func (circle *Circle) Radius() float64
func (circle *Circle) Rotate(radians float64)
func (circle *Circle) Rotation() float64
func (circle *Circle) Scale() Vector
func (circle *Circle) SetPosition(x, y float64)
func (circle *Circle) SetPositionVec(vec Vector)
func (circle *Circle) SetRadius(radius float64)
func (circle *Circle) SetRotation(rotation float64)
func (circle *Circle) SetScale(w, h float64)
func (circle *Circle) SetScaleVec(scale Vector)
type Collision
func (cc *Collision) ContactWithCell(cell *Cell) Vector
func (cc *Collision) ContactWithObject(object *Object) Vector
func (cc *Collision) HasTags(tags ...string) bool
func (cc *Collision) ObjectsByTags(tags ...string) []*Object
func (cc *Collision) SlideAgainstCell(cell *Cell, avoidTags ...string) (Vector, bool)
type ContactSet
func NewContactSet() *ContactSet
func (cs *ContactSet) BottommostPoint() Vector
func (cs *ContactSet) LeftmostPoint() Vector
func (cs *ContactSet) RightmostPoint() Vector
func (cs *ContactSet) TopmostPoint() Vector
type ConvexPolygon
func NewConvexPolygon(x, y float64, points ...float64) *ConvexPolygon
func NewConvexPolygonVec(position Vector, points ...Vector) *ConvexPolygon
func NewLine(x1, y1, x2, y2 float64) *ConvexPolygon
func NewRectangle(x, y, w, h float64) *ConvexPolygon
func (cp *ConvexPolygon) AddPoints(vertexPositions ...float64)
func (cp *ConvexPolygon) AddPointsVec(points ...Vector)
func (cp *ConvexPolygon) Bounds() (Vector, Vector)
func (cp *ConvexPolygon) Center() Vector
func (cp *ConvexPolygon) Clone() IShape
func (cp *ConvexPolygon) ContainedBy(otherShape IShape) bool
func (cp *ConvexPolygon) FlipH()
func (cp *ConvexPolygon) FlipV()
func (cp *ConvexPolygon) Intersection(dx, dy float64, other IShape) *ContactSet
func (p *ConvexPolygon) IntersectionForEach(dx, dy float64, f func(c *ContactSet) bool, others ...IShape)
func (cp *ConvexPolygon) Lines() []*collidingLine
func (cp *ConvexPolygon) Move(x, y float64)
func (cp *ConvexPolygon) MoveVec(vec Vector)
func (polygon *ConvexPolygon) PointInside(point Vector) bool
func (cp *ConvexPolygon) Position() Vector
func (cp *ConvexPolygon) Project(axis Vector) Projection
func (cp *ConvexPolygon) RecenterPoints()
func (cp *ConvexPolygon) ReverseVertexOrder()
func (polygon *ConvexPolygon) Rotate(radians float64)
func (polygon *ConvexPolygon) Rotation() float64
func (cp *ConvexPolygon) SATAxes() []Vector
func (polygon *ConvexPolygon) Scale() Vector
func (cp *ConvexPolygon) SetPosition(x, y float64)
func (cp *ConvexPolygon) SetPositionVec(vec Vector)
func (polygon *ConvexPolygon) SetRotation(radians float64)
func (polygon *ConvexPolygon) SetScale(x, y float64)
func (polygon *ConvexPolygon) SetScaleVec(vec Vector)
func (cp *ConvexPolygon) Transformed() []Vector
type IShape
type ModVector
func (ip ModVector) Add(other Vector) ModVector
func (ip ModVector) ClampAngle(baselineVector Vector, maxAngle float64) ModVector
func (ip ModVector) ClampMagnitude(maxMag float64) ModVector
func (ip ModVector) Clone() ModVector
func (ip ModVector) Divide(scalar float64) ModVector
func (ip ModVector) Expand(margin, min float64) ModVector
func (ip ModVector) Invert() ModVector
func (ip ModVector) Lerp(other Vector, percentage float64) ModVector
func (ip ModVector) Mult(other Vector) ModVector
func (ip ModVector) Rotate(angle float64) ModVector
func (ip ModVector) Round(snapToUnits float64) ModVector
func (ip ModVector) Scale(scalar float64) ModVector
func (ip ModVector) SetZero() ModVector
func (ip ModVector) Slerp(targetDirection Vector, percentage float64) ModVector
func (ip ModVector) String() string
func (ip ModVector) Sub(other Vector) ModVector
func (ip ModVector) SubMagnitude(mag float64) ModVector
func (ip ModVector) ToVector() Vector
func (ip ModVector) Unit() ModVector
type Object
func NewObject(x, y, w, h float64, tags ...string) *Object
func (obj *Object) AddTags(tags ...string)
func (obj *Object) AddToIgnoreList(ignoreObj *Object)
func (obj *Object) Bottom() float64
func (obj *Object) BoundsToSpace(dx, dy float64) (int, int, int, int)
func (obj *Object) CellPosition() (int, int)
func (obj *Object) Center() Vector
func (obj *Object) Check(dx, dy float64, tags ...string) *Collision
func (obj *Object) Clone() *Object
func (obj *Object) HasTags(tags ...string) bool
func (obj *Object) Overlaps(other *Object) bool
func (obj *Object) RemoveFromIgnoreList(ignoreObj *Object)
func (obj *Object) RemoveTags(tags ...string)
func (obj *Object) Right() float64
func (obj *Object) SetBottom(y float64)
func (obj *Object) SetBounds(topLeft, bottomRight Vector)
func (obj *Object) SetCenter(x, y float64)
func (obj *Object) SetRight(x float64)
func (obj *Object) SetShape(shape IShape)
func (obj *Object) SharesCells(other *Object) bool
func (obj *Object) SharesCellsTags(tags ...string) bool
func (obj *Object) Tags() []string
func (obj *Object) Update()
type Projection
func (projection Projection) IsInside(other Projection) bool
func (projection Projection) Overlap(other Projection) float64
func (projection Projection) Overlapping(other Projection) bool
type Space
func NewSpace(spaceWidth, spaceHeight, cellWidth, cellHeight int) *Space
func (sp *Space) Add(objects ...*Object)
func (sp *Space) Cell(x, y int) *Cell
func (sp *Space) CellsInLine(startX, startY, endX, endY int) []*Cell
func (sp *Space) CheckCells(x, y, w, h int, tags ...string) []*Object
func (sp *Space) CheckWorld(x, y, w, h float64, tags ...string) []*Object
func (sp *Space) CheckWorldVec(pos, size Vector, tags ...string) []*Object
func (sp *Space) Height() int
func (sp *Space) Objects() []*Object
func (sp *Space) Remove(objects ...*Object)
func (sp *Space) Resize(width, height int)
func (sp *Space) SpaceToWorld(x, y int) (float64, float64)
func (sp *Space) SpaceToWorldVec(x, y int) Vector
func (sp *Space) UnregisterAllObjects()
func (sp *Space) Width() int
func (sp *Space) WorldToSpace(x, y float64) (int, int)
func (sp *Space) WorldToSpaceVec(position Vector) (int, int)
type Vector
func NewVector(x, y float64) Vector
func NewVectorZero() Vector
func (vec Vector) Add(other Vector) Vector
func (vec Vector) Angle(other Vector) float64
func (vec Vector) AngleRotation() float64
func (vec Vector) ClampAngle(baselineVec Vector, maxAngle float64) Vector
func (vec Vector) ClampMagnitude(maxMag float64) Vector
func (vec Vector) Distance(other Vector) float64
func (vec Vector) DistanceSquared(other Vector) float64
func (vec Vector) Divide(scalar float64) Vector
func (vec Vector) Dot(other Vector) float64
func (vec Vector) Equals(other Vector) bool
func (vec Vector) Expand(margin, min float64) Vector
func (vec Vector) Floats() [2]float64
func (vec Vector) Invert() Vector
func (vec Vector) IsZero() bool
func (vec Vector) Lerp(other Vector, percentage float64) Vector
func (vec Vector) Magnitude() float64
func (vec Vector) MagnitudeSquared() float64
func (vec *Vector) Modify() ModVector
func (vec Vector) Mult(other Vector) Vector
func (vec Vector) Rotate(angle float64) Vector
func (vec Vector) Round(snapToUnits float64) Vector
func (vec Vector) Scale(scalar float64) Vector
func (vec Vector) Set(x, y float64) Vector
func (vec Vector) SetX(x float64) Vector
func (vec Vector) SetY(y float64) Vector
func (vec Vector) Slerp(targetDirection Vector, percentage float64) Vector
func (vec Vector) String() string
func (vec Vector) Sub(other Vector) Vector
func (vec Vector) SubMagnitude(mag float64) Vector
func (vec Vector) Unit() Vector
Constants ¶
This section is empty.

Variables ¶
View Source
var WorldDown = WorldUp.Invert()
WorldDown represents a unit vector in the global direction of WorldDown on the right-handed OpenGL / Tetra3D's coordinate system (+Y).

View Source
var WorldLeft = WorldRight.Invert()
WorldLeft represents a unit vector in the global direction of WorldLeft on the right-handed OpenGL / Tetra3D's coordinate system (-X).

View Source
var WorldRight = NewVector(1, 0)
WorldRight represents a unit vector in the global direction of WorldRight on the right-handed OpenGL / Tetra3D's coordinate system (+X).

View Source
var WorldUp = NewVector(0, 1)
WorldUp represents a unit vector in the global direction of WorldUp on the right-handed OpenGL / Tetra3D's coordinate system (+Y).

Functions ¶
func ToDegrees ¶
added in v0.6.0
func ToDegrees(radians float64) float64
ToDegrees is a helper function to easily convert radians to degrees for human readability.

func ToRadians ¶
added in v0.6.0
func ToRadians(degrees float64) float64
ToRadians is a helper function to easily convert degrees to radians (which is what the rotation-oriented functions in Tetra3D use).

Types ¶
type Cell ¶
type Cell struct {
	X, Y    int       // The X and Y position of the cell in the Space - note that this is in Grid position, not World position.
	Objects []*Object // The Objects that a Cell contains.
}
Cell is used to contain and organize Object information.

func (*Cell) Contains ¶
func (cell *Cell) Contains(obj *Object) bool
Contains returns whether a Cell contains the specified Object at its position.

func (*Cell) ContainsTags ¶
func (cell *Cell) ContainsTags(tags ...string) bool
ContainsTags returns whether a Cell contains an Object that has the specified tag at its position.

func (*Cell) Occupied ¶
func (cell *Cell) Occupied() bool
Occupied returns whether a Cell contains any Objects at all.

type Circle ¶
type Circle struct {
	// contains filtered or unexported fields
}
func NewCircle ¶
func NewCircle(x, y, radius float64) *Circle
NewCircle returns a new Circle, with its center at the X and Y position given, and with the defined radius.

func (*Circle) Bounds ¶
func (circle *Circle) Bounds() (Vector, Vector)
Bounds returns the top-left and bottom-right corners of the Circle.

func (*Circle) Clone ¶
func (circle *Circle) Clone() IShape
func (*Circle) Intersection ¶
func (circle *Circle) Intersection(dx, dy float64, other IShape) *ContactSet
Intersection tests to see if a Circle intersects with the other given Shape. dx and dy are delta movement variables indicating movement to be applied before the intersection check (thereby allowing you to see if a Shape would collide with another if it were in a different relative location). If an Intersection is found, a ContactSet will be returned, giving information regarding the intersection.

func (*Circle) IntersectionForEach ¶
added in v0.7.0
func (c *Circle) IntersectionForEach(dx, dy float64, f func(c *ContactSet) bool, others ...IShape)
IntersectionForEach runs a specified function for each contact set caused by contact with any of the shapes passed. If the custom function returns false, then the intersection testing stops iterating through further objects.

func (*Circle) IntersectionPointsCircle ¶
func (circle *Circle) IntersectionPointsCircle(other *Circle) []Vector
IntersectionPointsCircle returns the intersection points of the two circles provided.

func (*Circle) Move ¶
func (circle *Circle) Move(x, y float64)
Move translates the Circle by the designated X and Y values.

func (*Circle) MoveVec ¶
func (circle *Circle) MoveVec(vec Vector)
MoveVec translates the Circle by the designated Vector.

func (*Circle) PointInside ¶
func (circle *Circle) PointInside(point Vector) bool
PointInside returns if the given Vector is inside of the circle.

func (*Circle) Position ¶
func (circle *Circle) Position() Vector
Position() returns the X and Y position of the Circle.

func (*Circle) Radius ¶
func (circle *Circle) Radius() float64
Radius returns the radius of the Circle.

func (*Circle) Rotate ¶
added in v0.6.0
func (circle *Circle) Rotate(radians float64)
Circles can't rotate, of course. This function is just a stub to make them acceptable as IShapes.

func (*Circle) Rotation ¶
added in v0.6.0
func (circle *Circle) Rotation() float64
Circles can't rotate, of course. This function is just a stub to make them acceptable as IShapes.

func (*Circle) Scale ¶
added in v0.6.0
func (circle *Circle) Scale() Vector
Scale returns the scale multiplier of the Circle, twice; this is to have it adhere to the Shape interface.

func (*Circle) SetPosition ¶
func (circle *Circle) SetPosition(x, y float64)
SetPosition sets the center position of the Circle using the X and Y values given.

func (*Circle) SetPositionVec ¶
func (circle *Circle) SetPositionVec(vec Vector)
SetPosition sets the center position of the Circle using the Vector given.

func (*Circle) SetRadius ¶
added in v0.6.0
func (circle *Circle) SetRadius(radius float64)
SetRadius sets the radius of the Circle, updating the scale multiplier to reflect this change.

func (*Circle) SetRotation ¶
added in v0.6.0
func (circle *Circle) SetRotation(rotation float64)
Circles can't rotate, of course. This function is just a stub to make them acceptable as IShapes.

func (*Circle) SetScale ¶
added in v0.6.0
func (circle *Circle) SetScale(w, h float64)
SetScale sets the scale multiplier of the Circle (this is W and H to have it adhere to IShape as a contract; in truth, the Circle's radius will be set to 0.5 * the maximum out of the width and height height values given).

func (*Circle) SetScaleVec ¶
added in v0.7.0
func (circle *Circle) SetScaleVec(scale Vector)
SetScaleVec sets the scale multiplier of the Circle (this is W and H to have it adhere to IShape as a contract; in truth, the Circle's radius will be set to 0.5 * the maximum out of the width and height height values given).

type Collision ¶
type Collision struct {
	Objects []*Object // Slice of objects that were collided with; sorted according to distance to calling Object.
	Cells   []*Cell   // Slice of cells that were collided with; sorted according to distance to calling Object.
	// contains filtered or unexported fields
}
Collision contains the results of an Object.Check() call, and represents a collision between an Object and cells that contain other Objects. The Objects array indicate the Objects collided with.

func (*Collision) ContactWithCell ¶
func (cc *Collision) ContactWithCell(cell *Cell) Vector
ContactWithCell returns the delta to move to have the checking object come into contact with the specified Cell.

func (*Collision) ContactWithObject ¶
func (cc *Collision) ContactWithObject(object *Object) Vector
ContactWithObject returns the delta to move to have the checking object come into contact with the specified Object.

func (*Collision) HasTags ¶
func (cc *Collision) HasTags(tags ...string) bool
HasTags returns whether any objects within the Collision have all of the specified tags. This slice does not contain the Object that called Check().

func (*Collision) ObjectsByTags ¶
func (cc *Collision) ObjectsByTags(tags ...string) []*Object
ObjectsByTags returns a slice of Objects from the cells reported by a Collision object by searching for Objects with a specific set of tags. This slice does not contain the Object that called Check().

func (*Collision) SlideAgainstCell ¶
func (cc *Collision) SlideAgainstCell(cell *Cell, avoidTags ...string) (Vector, bool)
SlideAgainstCell returns how much distance the calling Object can slide to avoid a collision with the targetObject, and a boolean indicating if such a slide was possible. This only works on vertical and horizontal axes (x and y directly), primarily for platformers / top-down games. avoidTags is a sequence of tags (as strings) to indicate when sliding is valid (i.e. if a Cell contains an Object that has the tag given in the avoidTags slice, then sliding CANNOT happen).

type ContactSet ¶
type ContactSet struct {
	Points []Vector // Slice of points indicating contact between the two Shapes.
	MTV    Vector   // Minimum Translation Vector; this is the vector to move a Shape on to move it outside of its contacting Shape.
	Center Vector   // Center of the Contact set; this is the average of all Points contained within the Contact Set.
}
func NewContactSet ¶
func NewContactSet() *ContactSet
func (*ContactSet) BottommostPoint ¶
func (cs *ContactSet) BottommostPoint() Vector
BottommostPoint returns the bottom-most point out of the ContactSet's Points slice. If the Points slice is empty somehow, this returns nil.

func (*ContactSet) LeftmostPoint ¶
func (cs *ContactSet) LeftmostPoint() Vector
LeftmostPoint returns the left-most point out of the ContactSet's Points slice. If the Points slice is empty somehow, this returns nil.

func (*ContactSet) RightmostPoint ¶
func (cs *ContactSet) RightmostPoint() Vector
RightmostPoint returns the right-most point out of the ContactSet's Points slice. If the Points slice is empty somehow, this returns nil.

func (*ContactSet) TopmostPoint ¶
func (cs *ContactSet) TopmostPoint() Vector
TopmostPoint returns the top-most point out of the ContactSet's Points slice. If the Points slice is empty somehow, this returns nil.

type ConvexPolygon ¶
type ConvexPolygon struct {
	Points []Vector // Points represents the points constructing the ConvexPolygon.

	// X, Y           float64  // X and Y are the position of the ConvexPolygon.
	// ScaleW, ScaleH float64 // The width and height for scaling
	Closed bool // Closed is whether the ConvexPolygon is closed or not; only takes effect if there are more than 2 points.
	// contains filtered or unexported fields
}
ConvexPolygon represents a series of points, connected by lines, constructing a convex shape. The polygon has a position, a scale, a rotation, and may or may not be closed.

func NewConvexPolygon ¶
func NewConvexPolygon(x, y float64, points ...float64) *ConvexPolygon
NewConvexPolygon creates a new convex polygon at the position given, from the provided set of X and Y positions of 2D points (or vertices). You don't need to pass any points at this stage, but if you do, you should pass whole pairs. The points should generally be ordered clockwise, from X and Y of the first, to X and Y of the last. For example: NewConvexPolygon(30, 20, 0, 0, 10, 0, 10, 10, 0, 10) would create a 10x10 convex polygon square, with the vertices at {0,0}, {10,0}, {10, 10}, and {0, 10}, with the polygon itself occupying a position of 30, 20. You can also pass the points using vectors with ConvexPolygon.AddPointsVec().

func NewConvexPolygonVec ¶
added in v0.7.0
func NewConvexPolygonVec(position Vector, points ...Vector) *ConvexPolygon
func NewLine ¶
func NewLine(x1, y1, x2, y2 float64) *ConvexPolygon
NewLine is a helper function that returns a ConvexPolygon composed of a single line. The Polygon has a position of x1, y1, and the line stretches to x2-x1 and y2-y1.

func NewRectangle ¶
func NewRectangle(x, y, w, h float64) *ConvexPolygon
NewRectangle returns a rectangular ConvexPolygon at the {x, y} position given with the vertices ordered in clockwise order, positioned at {0, 0}, {w, 0}, {w, h}, {0, h}. TODO: In actuality, an AABBRectangle should be its own "thing" with its own optimized Intersection code check.

func (*ConvexPolygon) AddPoints ¶
func (cp *ConvexPolygon) AddPoints(vertexPositions ...float64)
AddPoints allows you to add points to the ConvexPolygon with a slice or selection of float64s, with each pair indicating an X or Y value for a point / vertex (i.e. AddPoints(0, 1, 2, 3) would add two points - one at {0, 1}, and another at {2, 3}).

func (*ConvexPolygon) AddPointsVec ¶
func (cp *ConvexPolygon) AddPointsVec(points ...Vector)
AddPointsVec allows you to add points to the ConvexPolygon with a slice of Vectors, each indicating a point / vertex.

func (*ConvexPolygon) Bounds ¶
func (cp *ConvexPolygon) Bounds() (Vector, Vector)
Bounds returns two Vectors, comprising the top-left and bottom-right positions of the bounds of the ConvexPolygon, post-transformation.

func (*ConvexPolygon) Center ¶
func (cp *ConvexPolygon) Center() Vector
Center returns the transformed Center of the ConvexPolygon.

func (*ConvexPolygon) Clone ¶
func (cp *ConvexPolygon) Clone() IShape
Clone returns a clone of the ConvexPolygon as an IShape.

func (*ConvexPolygon) ContainedBy ¶
func (cp *ConvexPolygon) ContainedBy(otherShape IShape) bool
ContainedBy returns if the ConvexPolygon is wholly contained by the other shape provided.

func (*ConvexPolygon) FlipH ¶
func (cp *ConvexPolygon) FlipH()
FlipH flips the ConvexPolygon's vertices horizontally, across the polygon's width, according to their initial offset when adding the points.

func (*ConvexPolygon) FlipV ¶
func (cp *ConvexPolygon) FlipV()
FlipV flips the ConvexPolygon's vertices vertically according to their initial offset when adding the points.

func (*ConvexPolygon) Intersection ¶
func (cp *ConvexPolygon) Intersection(dx, dy float64, other IShape) *ContactSet
Intersection tests to see if a ConvexPolygon intersects with the other given Shape. dx and dy are the delta movement to be applied before the intersection check (thereby allowing you to see if a Shape would collide with another if it were in a different relative location). If an Intersection is found, a ContactSet will be returned, giving information regarding the intersection.

func (*ConvexPolygon) IntersectionForEach ¶
added in v0.7.0
func (p *ConvexPolygon) IntersectionForEach(dx, dy float64, f func(c *ContactSet) bool, others ...IShape)
IntersectionForEach runs a specified function for each contact set caused by contact with any of the shapes passed. If the custom function returns false, then the intersection testing stops iterating through further objects.

func (*ConvexPolygon) Lines ¶
func (cp *ConvexPolygon) Lines() []*collidingLine
Lines returns a slice of transformed internalLines composing the ConvexPolygon.

func (*ConvexPolygon) Move ¶
func (cp *ConvexPolygon) Move(x, y float64)
Move translates the ConvexPolygon by the designated X and Y values.

func (*ConvexPolygon) MoveVec ¶
func (cp *ConvexPolygon) MoveVec(vec Vector)
MoveVec translates the ConvexPolygon by the designated Vector.

func (*ConvexPolygon) PointInside ¶
func (polygon *ConvexPolygon) PointInside(point Vector) bool
PointInside returns if a Point (a Vector) is inside the ConvexPolygon.

func (*ConvexPolygon) Position ¶
func (cp *ConvexPolygon) Position() Vector
Position returns the position of the ConvexPolygon.

func (*ConvexPolygon) Project ¶
func (cp *ConvexPolygon) Project(axis Vector) Projection
Project projects (i.e. flattens) the ConvexPolygon onto the provided axis.

func (*ConvexPolygon) RecenterPoints ¶
added in v0.6.1
func (cp *ConvexPolygon) RecenterPoints()
RecenterPoints recenters the vertices in the polygon, such that they are all equidistant from the center. For example, say you had a polygon with the following three points: {0, 0}, {10, 0}, {0, 16}. After calling cp.RecenterPoints(), the polygon's points would be at {-5, -8}, {5, -8}, {-5, 8}.

func (*ConvexPolygon) ReverseVertexOrder ¶
func (cp *ConvexPolygon) ReverseVertexOrder()
ReverseVertexOrder reverses the vertex ordering of the ConvexPolygon.

func (*ConvexPolygon) Rotate ¶
added in v0.6.0
func (polygon *ConvexPolygon) Rotate(radians float64)
Rotate is a helper function to rotate a ConvexPolygon by the radians given.

func (*ConvexPolygon) Rotation ¶
added in v0.6.0
func (polygon *ConvexPolygon) Rotation() float64
Rotation returns the rotation (in radians) of the ConvexPolygon.

func (*ConvexPolygon) SATAxes ¶
func (cp *ConvexPolygon) SATAxes() []Vector
SATAxes returns the axes of the ConvexPolygon for SAT intersection testing.

func (*ConvexPolygon) Scale ¶
added in v0.6.0
func (polygon *ConvexPolygon) Scale() Vector
Scale returns the scale multipliers of the ConvexPolygon.

func (*ConvexPolygon) SetPosition ¶
func (cp *ConvexPolygon) SetPosition(x, y float64)
SetPosition sets the position of the ConvexPolygon. The offset of the vertices compared to the X and Y position is relative to however you initially defined the polygon and added the vertices.

func (*ConvexPolygon) SetPositionVec ¶
func (cp *ConvexPolygon) SetPositionVec(vec Vector)
SetPositionVec allows you to set the position of the ConvexPolygon using a Vector. The offset of the vertices compared to the X and Y position is relative to however you initially defined the polygon and added the vertices.

func (*ConvexPolygon) SetRotation ¶
added in v0.6.0
func (polygon *ConvexPolygon) SetRotation(radians float64)
SetRotation sets the rotation for the ConvexPolygon; note that the rotation goes counter-clockwise from 0 to pi, and then from -pi at 180 down, back to 0. This rotation scheme follows the way math.Atan2() works.

func (*ConvexPolygon) SetScale ¶
added in v0.6.0
func (polygon *ConvexPolygon) SetScale(x, y float64)
SetScale sets the scale multipliers of the ConvexPolygon.

func (*ConvexPolygon) SetScaleVec ¶
added in v0.7.0
func (polygon *ConvexPolygon) SetScaleVec(vec Vector)
SetScaleVec sets the scale multipliers of the ConvexPolygon using the provided Vector.

func (*ConvexPolygon) Transformed ¶
func (cp *ConvexPolygon) Transformed() []Vector
Transformed returns the ConvexPolygon's points / vertices, transformed according to the ConvexPolygon's position.

type IShape ¶
added in v0.6.0
type IShape interface {
	// Intersection tests to see if a Shape intersects with the other given Shape. dx and dy are delta movement variables indicating
	// movement to be applied before the intersection check (thereby allowing you to see if a Shape would collide with another if it
	// were in a different relative location). If an Intersection is found, a ContactSet will be returned, giving information regarding
	// the intersection.
	Intersection(dx, dy float64, other IShape) *ContactSet
	// IntersectionForEach runs a specified function for each contact set caused by contact with any of
	// the shapes passed. If the custom function returns false, then the intersection testing stops
	// iterating through further objects.
	IntersectionForEach(dx, dy float64, f func(c *ContactSet) bool, others ...IShape)
	// Bounds returns the top-left and bottom-right points of the Shape.
	Bounds() (Vector, Vector)
	// Position returns the X and Y position of the Shape.
	Position() Vector
	// SetPosition allows you to place a Shape at another location.
	SetPosition(x, y float64)
	// SetPositionVec allows you to place a Shape at another location using a Vector.
	SetPositionVec(position Vector)

	// Rotation returns the current rotation value for the Shape.
	Rotation() float64

	// SetRotation sets the rotation value for the Shape.
	// Note that the rotation goes counter-clockwise from 0 at right to pi/2 in the upwards direction,
	// pi or -pi at left, -pi/2 in the downwards direction, and finally back to 0.
	// This can be visualized as follows:
	//
	//   U
	// L   R
	//   D
	//
	// R: 0
	// U: pi/2
	// L: pi / -pi
	// D: -pi/2
	SetRotation(radians float64)

	// Rotate rotates the IShape by the radians provided.
	// Note that the rotation goes counter-clockwise from 0 at right to pi/2 in the upwards direction,
	// pi or -pi at left, -pi/2 in the downwards direction, and finally back to 0.
	// This can be visualized as follows:
	//
	//   U
	// L   R
	//   D
	//
	// R: 0
	// U: pi/2
	// L: pi / -pi
	// D: -pi/2
	Rotate(radians float64)

	Scale() Vector // Returns the scale of the IShape (the radius for Circles).

	// Sets the overall scale of the IShape; 1.0 is 100% scale, 2.0 is 200%, and so on.
	// The greater of these values is used for the radius for Circles.
	SetScale(w, h float64)

	// Sets the overall scale of the IShape using the provided Vector; 1.0 is 100% scale, 2.0 is 200%, and so on.
	// The greater of these values is used for the radius for Circles.
	SetScaleVec(vec Vector)

	// Move moves the IShape by the x and y values provided.
	Move(x, y float64)
	// MoveVec moves the IShape by the movement values given in the vector provided.
	MoveVec(vec Vector)

	// Clone duplicates the IShape.
	Clone() IShape
}
type ModVector ¶
added in v0.7.0
type ModVector struct {
	*Vector
}
ModVector represents a reference to a Vector, made to facilitate easy method-chaining and modifications on that Vector (as you don't need to re-assign the results of a chain of operations to the original variable to "save" the results). Note that a ModVector is not meant to be used to chain methods on a vector to pass directly into a function; you can just use the normal vector functions for that purpose. ModVectors are pointers, which are allocated to the heap. This being the case, they should be slower relative to normal Vectors, so use them only in non-performance-critical parts of your application.

func (ModVector) Add ¶
added in v0.7.0
func (ip ModVector) Add(other Vector) ModVector
Add adds the other Vector provided to the ModVector. This function returns the calling ModVector for method chaining.

func (ModVector) ClampAngle ¶
added in v0.7.0
func (ip ModVector) ClampAngle(baselineVector Vector, maxAngle float64) ModVector
ClampAngle clamps the Vector such that it doesn't exceed the angle specified (in radians). This function returns the calling ModVector for method chaining.

func (ModVector) ClampMagnitude ¶
added in v0.7.0
func (ip ModVector) ClampMagnitude(maxMag float64) ModVector
ClampMagnitude clamps the overall magnitude of the Vector to the maximum magnitude specified. This function returns the calling ModVector for method chaining.

func (ModVector) Clone ¶
added in v0.7.0
func (ip ModVector) Clone() ModVector
Clone returns a ModVector of a clone of its backing Vector. This function returns the calling ModVector for method chaining.

func (ModVector) Divide ¶
added in v0.7.0
func (ip ModVector) Divide(scalar float64) ModVector
Divide divides a Vector by the given scalar (ignoring the W component). This function returns the calling ModVector for method chaining.

func (ModVector) Expand ¶
added in v0.7.0
func (ip ModVector) Expand(margin, min float64) ModVector
Expand expands the ModVector by the margin specified, in absolute units, if each component is over the minimum argument. To illustrate: Given a ModVector of {1, 0.1, -0.3}, ModVector.Expand(0.5, 0.2) would give you a ModVector of {1.5, 0.1, -0.8}. This function returns the calling ModVector for method chaining.

func (ModVector) Invert ¶
added in v0.7.0
func (ip ModVector) Invert() ModVector
Invert inverts all components of the calling Vector. This function returns the calling ModVector for method chaining.

func (ModVector) Lerp ¶
added in v0.7.0
func (ip ModVector) Lerp(other Vector, percentage float64) ModVector
Lerp performs a linear interpolation between the starting Vector and the provided other Vector, to the given percentage (ranging from 0 to 1). This function returns the calling ModVector for method chaining.

func (ModVector) Mult ¶
added in v0.7.0
func (ip ModVector) Mult(other Vector) ModVector
Mult performs Hadamard (component-wise) multiplication with the Vector on the other Vector provided. This function returns the calling ModVector for method chaining.

func (ModVector) Rotate ¶
added in v0.7.0
func (ip ModVector) Rotate(angle float64) ModVector
Rotate rotates the calling Vector by the angle provided (in radians). This function returns the calling ModVector for method chaining.

func (ModVector) Round ¶
added in v0.7.0
func (ip ModVector) Round(snapToUnits float64) ModVector
Round snaps the Vector's components to the given space in world units, returning a clone (e.g. Vector{0.1, 1.27, 3.33}.Snap(0.25) will return Vector{0, 1.25, 3.25}). This function returns the calling ModVector for method chaining.

func (ModVector) Scale ¶
added in v0.7.0
func (ip ModVector) Scale(scalar float64) ModVector
Scale scales the Vector by the scalar provided. This function returns the calling ModVector for method chaining.

func (ModVector) SetZero ¶
added in v0.7.0
func (ip ModVector) SetZero() ModVector
func (ModVector) Slerp ¶
added in v0.7.0
func (ip ModVector) Slerp(targetDirection Vector, percentage float64) ModVector
Slerp performs a linear interpolation between the starting Vector and the provided other Vector, to the given percentage (ranging from 0 to 1). This function returns the calling ModVector for method chaining.

func (ModVector) String ¶
added in v0.7.0
func (ip ModVector) String() string
String converts the ModVector to a string. Because it's a ModVector, it's represented with a *.

func (ModVector) Sub ¶
added in v0.7.0
func (ip ModVector) Sub(other Vector) ModVector
Sub subtracts the other Vector from the calling ModVector. This function returns the calling ModVector for method chaining.

func (ModVector) SubMagnitude ¶
added in v0.7.0
func (ip ModVector) SubMagnitude(mag float64) ModVector
SubMagnitude subtacts the given magnitude from the Vector's. If the vector's magnitude is less than the given magnitude to subtract, a zero-length Vector will be returned. This function returns the calling ModVector for method chaining.

func (ModVector) ToVector ¶
added in v0.7.0
func (ip ModVector) ToVector() Vector
func (ModVector) Unit ¶
added in v0.7.0
func (ip ModVector) Unit() ModVector
Unit normalizes the ModVector (sets it to be of unit length). It does not alter the W component of the Vector. This function returns the calling ModVector for method chaining.

type Object ¶
type Object struct {
	Shape         IShape      // A shape for more specific collision-checking.
	Space         *Space      // Reference to the Space the Object exists within
	Position      Vector      // The position of the Object in the Space
	Size          Vector      // The size of the Object in the Space
	TouchingCells []*Cell     // An array of Cells the Object is touching
	Data          interface{} // A pointer to a user-definable object
	// contains filtered or unexported fields
}
Object represents an object that can be spread across one or more Cells in a Space. An Object is essentially an AABB (Axis-Aligned Bounding Box) Rectangle.

func NewObject ¶
func NewObject(x, y, w, h float64, tags ...string) *Object
NewObject returns a new Object of the specified position and size.

func (*Object) AddTags ¶
func (obj *Object) AddTags(tags ...string)
AddTags adds tags to the Object.

func (*Object) AddToIgnoreList ¶
func (obj *Object) AddToIgnoreList(ignoreObj *Object)
AddToIgnoreList adds the specified Object to the Object's internal collision ignoral list. Cells that contain the specified Object will not be counted when calling Check().

func (*Object) Bottom ¶
func (obj *Object) Bottom() float64
Bottom returns the bottom Y coordinate of the Object (i.e. object.Y + object.H).

func (*Object) BoundsToSpace ¶
func (obj *Object) BoundsToSpace(dx, dy float64) (int, int, int, int)
BoundsToSpace returns the Space coordinates of the shape (x, y, w, and h), given its world position and size, and a supposed movement of dx and dy.

func (*Object) CellPosition ¶
func (obj *Object) CellPosition() (int, int)
CellPosition returns the cellular position of the Object's center in the Space.

func (*Object) Center ¶
func (obj *Object) Center() Vector
Center returns the center position of the Object.

func (*Object) Check ¶
func (obj *Object) Check(dx, dy float64, tags ...string) *Collision
Check checks the space around the object using the designated delta movement (dx and dy). This is done by querying the containing Space's Cells so that it can see if moving it would coincide with a cell that houses another Object (filtered using the given selection of tag strings). If so, Check returns a Collision. If no objects are found or the Object does not exist within a Space, this function returns nil.

func (*Object) Clone ¶
func (obj *Object) Clone() *Object
Clone clones the Object with its properties into another Object. It also clones the Object's Shape (if it has one).

func (*Object) HasTags ¶
func (obj *Object) HasTags(tags ...string) bool
HasTags indicates if an Object has any of the tags indicated.

func (*Object) Overlaps ¶
func (obj *Object) Overlaps(other *Object) bool
Overlaps returns if an Object overlaps another Object.

func (*Object) RemoveFromIgnoreList ¶
func (obj *Object) RemoveFromIgnoreList(ignoreObj *Object)
RemoveFromIgnoreList removes the specified Object from the Object's internal collision ignoral list. Objects removed from this list will once again be counted for Check().

func (*Object) RemoveTags ¶
func (obj *Object) RemoveTags(tags ...string)
RemoveTags removes tags from the Object.

func (*Object) Right ¶
func (obj *Object) Right() float64
Right returns the right X coordinate of the Object (i.e. object.X + object.W).

func (*Object) SetBottom ¶
func (obj *Object) SetBottom(y float64)
SetBottom sets the Y position of the Object so that the bottom edge is at the Y position given.

func (*Object) SetBounds ¶
func (obj *Object) SetBounds(topLeft, bottomRight Vector)
func (*Object) SetCenter ¶
func (obj *Object) SetCenter(x, y float64)
SetCenter sets the Object such that its center is at the X and Y position given.

func (*Object) SetRight ¶
func (obj *Object) SetRight(x float64)
SetRight sets the X position of the Object so the right edge is at the X position given.

func (*Object) SetShape ¶
added in v0.5.1
func (obj *Object) SetShape(shape IShape)
SetShape sets the Shape on the Object, in case you need to use precise per-Shape intersection detection. SetShape calls Object.Update() as well, so that it's able to update the Shape's position to match its Object as necessary. (If you don't use this, the Shape's position might not match the Object's, depending on if you set the Shape after you added the Object to a Space and if you don't call Object.Update() yourself afterwards.)

func (*Object) SharesCells ¶
func (obj *Object) SharesCells(other *Object) bool
SharesCells returns whether the Object occupies a cell shared by the specified other Object.

func (*Object) SharesCellsTags ¶
func (obj *Object) SharesCellsTags(tags ...string) bool
SharesCellsTags returns if the Cells the Object occupies have an object with the specified tags.

func (*Object) Tags ¶
func (obj *Object) Tags() []string
Tags returns the tags an Object has.

func (*Object) Update ¶
func (obj *Object) Update()
Update updates the object's association to the Cells in the Space. This should be called whenever an Object is moved. This is automatically called once when creating the Object, so you don't have to call it for static objects.

type Projection ¶
type Projection struct {
	Min, Max float64
}
}

func (Projection) IsInside ¶
func (projection Projection) IsInside(other Projection) bool
IsInside returns whether the Projection is wholly inside of the other, provided Projection.

func (Projection) Overlap ¶
func (projection Projection) Overlap(other Projection) float64
Overlap returns the amount that a Projection is overlapping with the other, provided Projection. Credit to https://dyn4j.org/2010/01/sat/#sat-nointer

func (Projection) Overlapping ¶
func (projection Projection) Overlapping(other Projection) bool
Overlapping returns whether a Projection is overlapping with the other, provided Projection. Credit to https://www.sevenson.com.au/programming/sat/

type Space ¶
type Space struct {
	Cells                 [][]*Cell
	CellWidth, CellHeight int // Width and Height of each Cell in "world-space" / pixels / whatever
}
Space represents a collision space. Internally, each Space contains a 2D array of Cells, with each Cell being the same size. Cells contain information on which Objects occupy those spaces.

func NewSpace ¶
func NewSpace(spaceWidth, spaceHeight, cellWidth, cellHeight int) *Space
NewSpace creates a new Space. spaceWidth and spaceHeight is the width and height of the Space (usually in pixels), which is then populated with cells of size cellWidth by cellHeight. Generally, you want cells to be the size of the smallest collide-able objects in your game, and you want to move Objects at a maximum speed of one cell size per collision check to avoid missing any possible collisions.

func (*Space) Add ¶
func (sp *Space) Add(objects ...*Object)
Add adds the specified Objects to the Space, updating the Space's cells to refer to the Object.

func (*Space) Cell ¶
func (sp *Space) Cell(x, y int) *Cell
Cell returns the Cell at the given cellular / spatial (not world) X and Y position in the Space. If the X and Y position are out of bounds, Cell() will return nil.

func (*Space) CellsInLine ¶
func (sp *Space) CellsInLine(startX, startY, endX, endY int) []*Cell
func (*Space) CheckCells ¶
func (sp *Space) CheckCells(x, y, w, h int, tags ...string) []*Object
CheckCells checks a set of cells (from x,y to x + w, y + h in cellular coordinates) and returns a slice of the objects found within those Cells. The objects must have any of the tags provided (if any are provided).

func (*Space) CheckWorld ¶
added in v0.7.0
func (sp *Space) CheckWorld(x, y, w, h float64, tags ...string) []*Object
CheckWorld checks the cells of the Grid with the given world coordinates. Internally, this is just syntactic sugar for calling Space.WorldToSpace() on the position and size given.

func (*Space) CheckWorldVec ¶
added in v0.7.0
func (sp *Space) CheckWorldVec(pos, size Vector, tags ...string) []*Object
CheckWorldVec checks the cells of the Grid with the given world coordinates. This function takes vectors for the position and size of the checked area. Internally, this is just syntactic sugar for calling Space.WorldToSpace() on the position and size given.

func (*Space) Height ¶
func (sp *Space) Height() int
Height returns the height of the Space grid in Cells (so a 320x240 Space with 16x16 cells would have a height of 15).

func (*Space) Objects ¶
func (sp *Space) Objects() []*Object
Objects loops through all Cells in the Space (from top to bottom, and from left to right) to return all Objects that exist in the Space. Of course, each Object is counted only once.

func (*Space) Remove ¶
func (sp *Space) Remove(objects ...*Object)
Remove removes the specified Objects from being associated with the Space. This should be done whenever an Object is removed from the game.

func (*Space) Resize ¶
func (sp *Space) Resize(width, height int)
Resize resizes the internal Cells array.

func (*Space) SpaceToWorld ¶
func (sp *Space) SpaceToWorld(x, y int) (float64, float64)
SpaceToWorld converts from a position in the Space (on a grid) to a world-based position, given the size of the Space when first created.

func (*Space) SpaceToWorldVec ¶
added in v0.7.0
func (sp *Space) SpaceToWorldVec(x, y int) Vector
func (*Space) UnregisterAllObjects ¶
func (sp *Space) UnregisterAllObjects()
UnregisterAllObjects unregisters all Objects registered to Cells in the Space.

func (*Space) Width ¶
func (sp *Space) Width() int
Width returns the width of the Space grid in Cells (so a 320x240 Space with 16x16 cells would have a width of 20).

func (*Space) WorldToSpace ¶
func (sp *Space) WorldToSpace(x, y float64) (int, int)
WorldToSpace converts from a world position (x, y) to a position in the Space (a grid-based position).

func (*Space) WorldToSpaceVec ¶
added in v0.7.0
func (sp *Space) WorldToSpaceVec(position Vector) (int, int)
WorldToSpaceVec converts from a world position Vector to a position in the Space (a grid-based position).

type Vector ¶
added in v0.7.0
type Vector struct {
	X float64 // The X (1st) component of the Vector
	Y float64 // The Y (2nd) component of the Vector
}
Vector represents a 2D Vector, which can be used for usual 2D applications (position, direction, velocity, etc). Any Vector functions that modify the calling Vector return copies of the modified Vector, meaning you can do method-chaining easily. Vectors seem to be most efficient when copied (so try not to store pointers to them if possible, as dereferencing pointers can be more inefficient than directly acting on data, and storing pointers moves variables to heap).

func NewVector ¶
added in v0.7.0
func NewVector(x, y float64) Vector
NewVector creates a new Vector with the specified x, y, and z components. The W component is generally ignored for most purposes.

func NewVectorZero ¶
added in v0.7.0
func NewVectorZero() Vector
NewVectorZero creates a new "zero-ed out" Vector, with the values of 0, 0, 0, and 0 (for W).

func (Vector) Add ¶
added in v0.7.0
func (vec Vector) Add(other Vector) Vector
Plus returns a copy of the calling vector, added together with the other Vector provided (ignoring the W component).

func (Vector) Angle ¶
added in v0.7.0
func (vec Vector) Angle(other Vector) float64
Angle returns the angle between the calling Vector and the provided other Vector (ignoring the W component).

func (Vector) AngleRotation ¶
added in v0.7.0
func (vec Vector) AngleRotation() float64
func (Vector) ClampAngle ¶
added in v0.7.0
func (vec Vector) ClampAngle(baselineVec Vector, maxAngle float64) Vector
ClampAngle clamps the Vector such that it doesn't exceed the angle specified (in radians). This function returns a normalized (unit) Vector.

func (Vector) ClampMagnitude ¶
added in v0.7.0
func (vec Vector) ClampMagnitude(maxMag float64) Vector
ClampMagnitude clamps the overall magnitude of the Vector to the maximum magnitude specified, returning a copy with the result.

func (Vector) Distance ¶
added in v0.7.0
func (vec Vector) Distance(other Vector) float64
Distance returns the distance from the calling Vector to the other Vector provided.

func (Vector) DistanceSquared ¶
added in v0.7.0
func (vec Vector) DistanceSquared(other Vector) float64
Distance returns the squared distance from the calling Vector to the other Vector provided. This is faster than Distance(), as it avoids using math.Sqrt().

func (Vector) Divide ¶
added in v0.7.0
func (vec Vector) Divide(scalar float64) Vector
Divide divides a Vector by the given scalar (ignoring the W component), returning a copy with the result.

func (Vector) Dot ¶
added in v0.7.0
func (vec Vector) Dot(other Vector) float64
Dot returns the dot product of a Vector and another Vector (ignoring the W component).

func (Vector) Equals ¶
added in v0.7.0
func (vec Vector) Equals(other Vector) bool
Equals returns true if the two Vectors are close enough in all values (excluding W).

func (Vector) Expand ¶
added in v0.7.0
func (vec Vector) Expand(margin, min float64) Vector
Expand expands the Vector by the margin specified, in absolute units, if each component is over the minimum argument. To illustrate: Given a Vector of {1, 0.1, -0.3}, Vector.Expand(0.5, 0.2) would give you a Vector of {1.5, 0.1, -0.8}. This function returns a copy of the Vector with the result.

func (Vector) Floats ¶
added in v0.7.0
func (vec Vector) Floats() [2]float64
Floats returns a [2]float64 array consisting of the Vector's contents.

func (Vector) Invert ¶
added in v0.7.0
func (vec Vector) Invert() Vector
Invert returns a copy of the Vector with all components inverted.

func (Vector) IsZero ¶
added in v0.7.0
func (vec Vector) IsZero() bool
IsZero returns true if the values in the Vector are extremely close to 0 (excluding W).

func (Vector) Lerp ¶
added in v0.7.0
func (vec Vector) Lerp(other Vector, percentage float64) Vector
Lerp performs a linear interpolation between the starting Vector and the provided other Vector, to the given percentage (ranging from 0 to 1).

func (Vector) Magnitude ¶
added in v0.7.0
func (vec Vector) Magnitude() float64
Magnitude returns the length of the Vector.

func (Vector) MagnitudeSquared ¶
added in v0.7.0
func (vec Vector) MagnitudeSquared() float64
MagnitudeSquared returns the squared length of the Vector; this is faster than Length() as it avoids using math.Sqrt().

func (*Vector) Modify ¶
added in v0.7.0
func (vec *Vector) Modify() ModVector
Modify returns a ModVector object (a pointer to the original vector).

func (Vector) Mult ¶
added in v0.7.0
func (vec Vector) Mult(other Vector) Vector
Mult performs Hadamard (component-wise) multiplication on the calling Vector with the other Vector provided, returning a copy with the result (and ignoring the Vector's W component).

func (Vector) Rotate ¶
added in v0.7.0
func (vec Vector) Rotate(angle float64) Vector
Rotate returns a copy of the Vector, rotated around an axis Vector with the x, y, and z components provided, by the angle provided (in radians), counter-clockwise. The function is most efficient if passed an orthogonal, normalized axis (i.e. the X, Y, or Z constants). Note that this function ignores the W component of both Vectors.

func (Vector) Round ¶
added in v0.7.0
func (vec Vector) Round(snapToUnits float64) Vector
Round rounds off the Vector's components to the given space in world unit increments, returning a clone (e.g. Vector{0.1, 1.27, 3.33}.Snap(0.25) will return Vector{0, 1.25, 3.25}).

func (Vector) Scale ¶
added in v0.7.0
func (vec Vector) Scale(scalar float64) Vector
Scale scales a Vector by the given scalar (ignoring the W component), returning a copy with the result.

func (Vector) Set ¶
added in v0.7.0
func (vec Vector) Set(x, y float64) Vector
Set sets the values in the Vector to the x, y, and z values provided.

func (Vector) SetX ¶
added in v0.7.0
func (vec Vector) SetX(x float64) Vector
SetX sets the X component in the vector to the value provided.

func (Vector) SetY ¶
added in v0.7.0
func (vec Vector) SetY(y float64) Vector
SetY sets the Y component in the vector to the value provided.

func (Vector) Slerp ¶
added in v0.7.0
func (vec Vector) Slerp(targetDirection Vector, percentage float64) Vector
Slerp performs a spherical linear interpolation between the starting Vector and the provided ending Vector, to the given percentage (ranging from 0 to 1). This should be done with directions, usually, rather than positions. This being the case, this normalizes both Vectors.

func (Vector) String ¶
added in v0.7.0
func (vec Vector) String() string
String returns a string representation of the Vector, excluding its W component (which is primarily used for internal purposes).

func (Vector) Sub ¶
added in v0.7.0
func (vec Vector) Sub(other Vector) Vector
Sub returns a copy of the calling Vector, with the other Vector subtracted from it (ignoring the W component).

func (Vector) SubMagnitude ¶
added in v0.7.0
func (vec Vector) SubMagnitude(mag float64) Vector
SubMagnitude subtracts the given magnitude from the Vector's existing magnitude. If the vector's magnitude is less than the given magnitude to subtract, a zero-length Vector will be returned.

func (Vector) Unit ¶
added in v0.7.0
func (vec Vector) Unit() Vector
Unit returns a copy of the Vector, normalized (set to be of unit length). It does not alter the W component of the Vector.