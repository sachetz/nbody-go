package particle

import (
	"math"
	"proj3/barrier"
	"sync"
)

// Randomly initialize particles in parallel mode
func InitialiseRandomPointsParallel(p []*Particle, nParticles int, nThreads int) {
	b := barrier.NewBarrier(nThreads + 1)
	for i := 0; i < nThreads; i++ {
		lowerBound := i * nParticles / nThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/nThreads), float64(nParticles)))
		f := func(tid int, bar *barrier.Barrier) {
			InitialiseRandomPointsSequential(p, lowerBound, upperBound)
			bar.Wait()
		}
		go f(i, b)
	}
	b.Wait()
}

// Randomly initialize particles in a circle in parallel mode
func InitialiseParticlesInCircleParallel(p []*Particle, nParticles int, nThreads int) {
	b := barrier.NewBarrier(nThreads + 1)
	for i := 0; i < nThreads; i++ {
		lowerBound := i * nParticles / nThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/nThreads), float64(nParticles)))
		f := func(tid int, bar *barrier.Barrier) {
			InitialiseParticlesInCircleSequential(p, lowerBound, upperBound, nParticles)
			bar.Wait()
		}
		go f(i, b)
	}
	b.Wait()
}

// Update the positions of the particles in parallel mode
func UpdatePosParallel(p []*Particle, nParticles int, nThreads int) {
	b := barrier.NewBarrier(nThreads + 1)
	for i := 0; i < nThreads; i++ {
		lowerBound := i * nParticles / nThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/nThreads), float64(nParticles)))
		f := func(tid int, bar *barrier.Barrier) {
			UpdatePosSequential(p, lowerBound, upperBound)
			bar.Wait()
		}
		go f(i, b)
	}
	b.Wait()
}

// Find the bounds of the point coordinates in parallel mode
func FindBoundsParallel(p []*Particle, nParticles int, nThreads int) float64 {
	m := sync.Mutex{}
	var max float64
	b := barrier.NewBarrier(nThreads + 1)
	for i := 0; i < nThreads; i++ {
		lowerBound := i * nParticles / nThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/nThreads), float64(nParticles)))
		f := func(tid int, bar *barrier.Barrier) {
			localMax := FindBoundsSequential(p, lowerBound, upperBound)
			m.Lock()
			max = math.Max(localMax, max)
			m.Unlock()
			bar.Wait()
		}
		go f(i, b)
	}
	b.Wait()
	return max
}
