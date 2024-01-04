package poisson

import (
	"math/rand"
	"testing"
)

func TestSample2D(t *testing.T) {
	d := 10.0
	k := 10
	b := &Bounds{
		MinX: -50,
		MinY: -50,
		MaxX: 25,
		MaxY: 75,
	}
	s := rand.NewSource(943050168) // credit to https://www.random.org/

	points := Sample2D(d, k, b, s)

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
