package particle

import (
	"math"
	"math/rand"
)

const SOFTENING float64 = 1e-9

type Particle struct {
	X, Y   float64 // coordinates
	Vx, Vy float64 // momenta
}

// Calculate squared distance
func CalcDistance(p1 *Particle, p2 *Particle) float64 {
	return math.Sqrt((p2.X-p1.X)*(p2.X-p1.X) + (p2.Y-p1.Y)*(p2.Y-p1.Y) + SOFTENING)
}

// Calculate force between two particles
func CalcForce(p1 *Particle, p2 *Particle, dt float64) {
	var Fx float64 = 0.0
	var Fy float64 = 0.0
	if p1 != p2 {
		var invDist float64 = 1.0 / CalcDistance(p1, p2)
		var invDist3 float64 = invDist * invDist * invDist
		Fx += (p2.X - p1.X) * invDist3
		Fy += (p2.Y - p1.Y) * invDist3
		p1.Vx += dt * Fx
		p1.Vy += dt * Fy
	}
}

// Randomly initialize particle positions and momenta
func InitialiseRandomPoints(data []*Particle, n int) {
	for i := 0; i < n; i++ {
		data[i] = &Particle{}
		data[i].X = 2.0*rand.Float64() - 1.0
		data[i].Y = 2.0*rand.Float64() - 1.0
		data[i].Vx = 2.0*rand.Float64() - 1.0
		data[i].Vy = 2.0*rand.Float64() - 1.0
	}
}

// Randomly initialize particles in a circle - for testing
func InitialiseParticlesInCircle(p []*Particle, n int) {
	var r float64 = 1.0 // Radius of the circle
	for i := 0; i < n; i++ {
		p[i] = &Particle{}
		var a float64 = 2 * math.Pi * float64(i) / float64(n) // Angle of the new particle
		p[i].X = r * math.Cos(a)
		p[i].Y = r * math.Sin(a) // Circle centered at (0,0)
		p[i].Vx = 0
		p[i].Vy = 0
	}
}
