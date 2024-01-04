package poisson

import (
	"math"
	"math/rand"
)

type Bounds struct {
	MinX float64
	MinY float64
	MaxX float64
	MaxY float64
}

func (b *Bounds) Dx() float64 {
	return b.MaxX - b.MinX
}
func (b *Bounds) Dy() float64 {
	return b.MaxY - b.MinY
}

type Point2D struct {
	X float64
	Y float64
}

func (p *Point2D) dist(p1 *Point2D) float64 {
	dx := p1.X - p.X
	dy := p1.Y - p.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// Sample2D generates points with a minimum distance d between points and within bounds b.
// At each iteration, it will try k random points around a previous point, per Bridson's algorithm.
// The algorithm will use the given random source for repeatability.
func Sample2D(d float64, k int, b *Bounds, s rand.Source) []Point2D {
	r := rand.New(s)
	cellSize := math.Floor(d / math.Sqrt(2))
	cellsWidth := int(math.Ceil(b.Dx()/cellSize) + 1)
	cellsHeight := int(math.Ceil(b.Dy()/cellSize) + 1)
	points := make([]Point2D, 0, cellsWidth*cellsHeight)
	active := make([]Point2D, 0, int(math.Sqrt(float64(cellsWidth*cellsHeight))*4))
	grid := make([][]*Point2D, cellsWidth)
	for i := range grid {
		grid[i] = make([]*Point2D, cellsHeight)
	}

	p0 := Point2D{
		X: r.Float64()*b.Dx() + b.MinX,
		Y: r.Float64()*b.Dy() + b.MinY,
	}
	insertPoint(grid, cellSize, p0, b)
	points = append(points, p0)
	active = append(active, p0)

	for len(active) > 0 {
		i := r.Intn(len(active))
		p0 = active[i]
		found := false
		for t := 0; t < k; t++ {
			theta := r.Float64() * math.Pi * 2
			radius := r.Float64()*d + d
			p1 := Point2D{
				X: p0.X + radius*math.Cos(theta),
				Y: p0.Y + radius*math.Sin(theta),
			}
			if !isValidPoint(grid, &p1, cellSize, d, b) {
				continue
			}
			insertPoint(grid, cellSize, p1, b)
			points = append(points, p1)
			active = append(active, p1)
			found = true
			break
		}
		if !found {
			active[i] = active[len(active)-1]
			active = active[:len(active)-1]
		}
	}

	return points
}

func isValidPoint(grid [][]*Point2D, p *Point2D, cellSize float64, d float64, b *Bounds) bool {
	if p.X < b.MinX || p.X >= b.MaxX || p.Y < b.MinY || p.Y >= b.MaxY {
		return false
	}
	xindex := int(math.Floor((p.X - b.MinX) / cellSize))
	yindex := int(math.Floor((p.Y - b.MinY) / cellSize))
	x0 := max(xindex-2, 0)
	x1 := min(xindex+2, len(grid)-1)
	y0 := max(yindex-2, 0)
	y1 := min(yindex+2, len(grid[0])-1)
	for x := x0; x <= x1; x++ {
		for y := y0; y <= y1; y++ {
			if grid[x][y] != nil {
				if grid[x][y].dist(p) < d {
					return false
				}
			}
		}
	}
	return true
}

func insertPoint(grid [][]*Point2D, cellSize float64, p Point2D, b *Bounds) {
	xindex := int(math.Floor((p.X - b.MinX) / cellSize))
	yindex := int(math.Floor((p.Y - b.MinY) / cellSize))
	grid[xindex][yindex] = &p
}
