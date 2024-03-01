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
		f := func(tid int) {
			defer wg.Done()
			InitialiseRandomPointsSequential(p, lowerBound, upperBound)
		}
		go f(i)
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
		f := func(tid int) {
			defer wg.Done()
			InitialiseParticlesInCircleSequential(p, lowerBound, upperBound, nParticles)
		}
		go f(i)
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
		f := func(tid int) {
			defer wg.Done()
			localMax := FindBoundsSequential(p, lowerBound, upperBound)
			m.Lock()
			max = math.Max(localMax, max)
			m.Unlock()
		}
		go f(i)
	}
	wg.Wait()
	return max
}
