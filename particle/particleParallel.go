package particle

import (
	"math"
	"math/rand"
	"nbody_go/queue"
	"sync"
	"sync/atomic"
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

func FindBoundsWorkStealing(p []*Particle, nParticles int, nThreads int, workQueues []*queue.BoundedDeque, completed *atomic.Int32) float64 {
	m := sync.Mutex{}
	var max float64
	var wg sync.WaitGroup
	for i := 0; i < nThreads; i++ {
		wg.Add(1)
		f := func(tid int, workQueues []*queue.BoundedDeque) {
			defer wg.Done()

			var maxLocal float64 = math.Abs(p[0].X)

			// Work on own queue
			for {
				idx := workQueues[tid].PopBottom()
				if idx == -1 {
					completed.Add(1)
					break
				}
				FindBounds(p, idx, &maxLocal)
			}

			// If any randomly selected queue is non empty, steal from it and work on that job
			for completed.Load() < int32(nThreads) {
				target := rand.Int31n(int32(nThreads))
				stolenWork := workQueues[target].PopTop()
				if stolenWork == -1 {
					continue
				}
				FindBounds(p, stolenWork, &maxLocal)
			}

			m.Lock()
			max = math.Max(maxLocal, max)
			m.Unlock()
		}
		go f(i, workQueues)
	}
	wg.Wait()
	return max
}
