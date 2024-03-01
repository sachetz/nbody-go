package particle

import (
	"math"
	"sync"
)

// Randomly initialize particles in parallel mode
func InitialiseRandomPointsParallel(p []*Particle, nParticles int, nThreads int) {
	var wg sync.WaitGroup
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		lowerBound := i * nParticles / nThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/nThreads), float64(nParticles)))
		f := func() {
			defer wg.Done()
			InitialiseRandomPointsSequential(p, lowerBound, upperBound)
		}
		go f()
	}
	wg.Wait()
}

// Randomly initialize particles in a circle in parallel mode
func InitialiseParticlesInCircleParallel(p []*Particle, nParticles int, nThreads int) {
	var wg sync.WaitGroup
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		lowerBound := i * nParticles / nThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/nThreads), float64(nParticles)))
		f := func() {
			defer wg.Done()
			InitialiseParticlesInCircleSequential(p, lowerBound, upperBound, nParticles)
		}
		go f()
	}
	wg.Wait()
}

// Update the positions of the particles in parallel mode
func UpdatePosParallel(p []*Particle, nParticles int, nThreads int) {
	var wg sync.WaitGroup
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		lowerBound := i * nParticles / nThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/nThreads), float64(nParticles)))
		f := func() {
			defer wg.Done()
			UpdatePosSequential(p, lowerBound, upperBound)
		}
		go f()
	}
	wg.Wait()
}

// Find the bounds of the point coordinates in parallel mode
func FindBoundsParallel(p []*Particle, nParticles int, nThreads int) float64 {
	m := sync.Mutex{}
	var max float64
	var wg sync.WaitGroup
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		lowerBound := i * nParticles / nThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/nThreads), float64(nParticles)))
		f := func() {
			defer wg.Done()
			localMax := FindBoundsSequential(p, lowerBound, upperBound)
			m.Lock()
			max = math.Max(localMax, max)
			m.Unlock()
		}
		go f()
	}
	wg.Wait()
	return max
}
