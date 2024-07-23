package poisson

import (
	"math/rand/v2"
	"testing"
)

func check(t *testing.T, points []Point2D, b *Bounds, d float64) {
	if len(points) == 0 {
		t.Fatalf("No points generated")
	}

	for i, p1 := range points {
		if p1.X < b.MinX || p1.X >= b.MaxX {
			t.Fatalf("Invalid X value %f", p1.X)
		}
		if p1.Y < b.MinY || p1.Y >= b.MaxY {
			t.Fatalf("Invalid Y value %f", p1.Y)
		}
		for j := i + 1; j < len(points); j++ {
			p2 := points[j]
			dist := p1.dist(&p2)
			if dist < d {
				t.Fatalf("Points too close %+v %+v (%f)", p1, p2, dist)
			}
		}
	}
}

func TestSample2D(t *testing.T) {
	d := 10.0
	k := 10
	b := &Bounds{
		MinX: -50,
		MinY: -50,
		MaxX: 25,
		MaxY: 75,
	}
	s := rand.NewPCG(8568094394964690136, 11528959135235502846)

	points := Sample2D(d, k, b, s)

	check(t, points, b, d)
}

func TestSample2DSmall(t *testing.T) {
	d := 0.01
	k := 10
	b := &Bounds{
		MinX: -0.05,
		MinY: -0.05,
		MaxX: 0.025,
		MaxY: 0.075,
	}
	s := rand.NewPCG(1086012171205092109, 14631601773360610953)

	points := Sample2D(d, k, b, s)

	check(t, points, b, d)
}
