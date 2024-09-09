package poisson

import (
	"math"
	"math/rand/v2"
)

// Bounds is a rectangular area.
type Bounds struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

func (b *Bounds) width() float64 {
	return b.MaxX - b.MinX
}
func (b *Bounds) height() float64 {
	return b.MaxY - b.MinY
}

// Point is a 2D point.
type Point struct {
	X float64
	Y float64
}

func (p *Point) distanceTo(p1 *Point) float64 {
	dx := p1.X - p.X
	dy := p1.Y - p.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Sample2D generates a set of 2D points using Poisson-disc sampling, which ensures
// that no two points are closer than a specified minimum distance.
//
// Parameters:
//   - distance: The minimum distance between any two points.
//   - kTries: The number of attempts to generate a new point around an active point before removing the active point from consideration.
//   - bounds: The rectangular bounds within which points should be generated.
//   - start: An optional starting point. If nil, a random point within the bounds is used.
//   - source: A random source for repeatability.
//
// Returns:
//   - A slice of points that satisfy the Poisson-disc sampling criteria.
func Sample2D(distance float64, kTries int, bounds Bounds, start *Point, source rand.Source) []Point {
	rnd := rand.New(source)

	// Calculate cell size and grid dimensions
	cellSize := distance / math.Sqrt(2)
	cellsWidth := int(math.Ceil(bounds.width()/cellSize) + 1)
	cellsHeight := int(math.Ceil(bounds.height()/cellSize) + 1)

	// Result points
	points := make([]Point, 0, cellsWidth*cellsHeight/3)
	// Points to consider adding another point around
	active := make([]Point, 0, int(math.Sqrt(float64(cellsWidth*cellsHeight))*4))

	// Grid to hold points to efficiently check for nearby points
	grid := make([][]*Point, cellsWidth)
	for i := range grid {
		grid[i] = make([]*Point, cellsHeight)
	}

	// Initial point, either at the specified start location or randomly within bounds
	var p0 Point
	if start == nil {
		p0 = Point{
			X: rnd.Float64()*bounds.width() + bounds.MinX,
			Y: rnd.Float64()*bounds.height() + bounds.MinY,
		}
	} else {
		p0 = *start
	}
	if insertPoint(grid, p0, cellSize, distance, &bounds) {
		points = append(points, p0)
		active = append(active, p0)
	}

	for len(active) > 0 {
		// Choose a random active point
		ai := rnd.IntN(len(active))
		p0 = active[ai]

		found := false
		for k := 0; k < kTries; k++ {
			// Generate random points around the active point
			theta := rnd.Float64() * math.Pi * 2
			radius := rnd.Float64()*distance + distance
			p1 := Point{
				X: p0.X + radius*math.Cos(theta),
				Y: p0.Y + radius*math.Sin(theta),
			}
			if insertPoint(grid, p1, cellSize, distance, &bounds) {
				points = append(points, p1)
				active = append(active, p1)
				found = true
				break
			}
		}

		// If no valid point was found, remove the active point from the list
		if !found {
			active[ai] = active[len(active)-1]
			active = active[:len(active)-1]
		}
	}

	return points
}

// insertPoint tries to insert a point into the grid while ensuring that it is within bounds and not too close to any existing points.
//
// Returns whether the point was inserted into the grid.
func insertPoint(grid [][]*Point, p Point, cellSize float64, distance float64, bounds *Bounds) bool {
	if p.X < bounds.MinX || p.X > bounds.MaxX || p.Y < bounds.MinY || p.Y > bounds.MaxY {
		return false
	}

	// Grid cell indices for the point
	xindex := int(math.Floor((p.X - bounds.MinX) / cellSize))
	yindex := int(math.Floor((p.Y - bounds.MinY) / cellSize))

	// Check nearby cells for nearby points
	x0 := max(xindex-2, 0)
	x1 := min(xindex+2, len(grid)-1)
	y0 := max(yindex-2, 0)
	y1 := min(yindex+2, len(grid[0])-1)
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			if grid[x][y] != nil {
				if grid[x][y].distanceTo(&p) < distance {
					return false
				}
			}
		}
	}

	// No nearby points, insert the point into the grid
	grid[xindex][yindex] = &p
	return true
}
