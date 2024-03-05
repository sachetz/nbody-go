package particle

import (
	"math"
	"math/rand"
	"proj3/utils"
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

// Randomly initialize particles in parallel mode
func InitialiseRandomPointsSequential(data []*Particle, lowerBound int, upperBound int) {
	for i := lowerBound; i < upperBound; i++ {
		data[i] = &Particle{}
		data[i].X = 2.0*rand.Float64() - 1.0
		data[i].Y = 2.0*rand.Float64() - 1.0
		data[i].Vx = 2.0*rand.Float64() - 1.0
		data[i].Vy = 2.0*rand.Float64() - 1.0
	}
}

// Randomly initialize particles in a circle in sequential mode
func InitialiseParticlesInCircleSequential(p []*Particle, lowerBound int, upperBound int, nParticles int) {
	var r float64 = 1.0 // Radius of the circle
	for i := lowerBound; i < upperBound; i++ {
		p[i] = &Particle{}
		var a float64 = 2 * math.Pi * float64(i) / float64(nParticles) // Angle of the new particle
		p[i].X = r * math.Cos(a)
		p[i].Y = r * math.Sin(a) // Circle centered at (0,0)
		p[i].Vx = 0
		p[i].Vy = 0
	}
}

// Update the positions of the particles in sequential mode
func UpdatePosSequential(p []*Particle, lowerBound int, upperBound int) {
	for i := lowerBound; i < upperBound; i++ {
		p[i].X += p[i].Vx * utils.Dt
		p[i].Y += p[i].Vy * utils.Dt
	}
}

// Find the bounds of the point coordinates in sequential mode
func FindBoundsSequential(p []*Particle, lowerBound int, upperBound int) float64 {
	var max float64 = math.Abs(p[lowerBound].X)
	for i := lowerBound; i < upperBound; i++ {
		FindBounds(p, i, &max)
	}
	return max
}

func FindBounds(p []*Particle, idx int, max *float64) {
	*max = math.Max(*max, math.Abs(p[idx].X))
	*max = math.Max(*max, math.Abs(p[idx].Y))
}