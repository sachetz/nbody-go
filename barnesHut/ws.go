package barnesHut

import (
	"fmt"
	"math"
	"os"
	"nbody_go/particle"
	"nbody_go/queue"
	"nbody_go/tree"
	"nbody_go/utils"
	"sync/atomic"
	"time"
)

func resetQueues(numThreads int, nParticles int) ([]*queue.BoundedDeque, *atomic.Int32) {
	workQueues := make([]*queue.BoundedDeque, numThreads) // Create work stealing queues for each thread
	for i := 0; i < numThreads; i++ {
		lowerBound := i * nParticles / numThreads
		upperBound := int(math.Min(float64((i+1)*nParticles/numThreads), float64(nParticles)))
		workQueues[i] = queue.NewBoundedDeque(lowerBound, upperBound)
	}
	var completed atomic.Int32
	completed.Store(0)
	return workQueues, &completed
}

func Ws() {
	start := time.Now()
	datafile, err := os.Create("particles_ws.dat") // Output file for particle positions
	utils.Check(err)
	defer datafile.Close()

	config := utils.GetParams()
	nParticles := config.NParticles
	nIters := config.NIters
	numThreads := config.NumThreads
	logging := config.Logging
	initPoints := config.InitPoints

	workQueues, completed := resetQueues(numThreads, nParticles)

	var p []*particle.Particle = make([]*particle.Particle, nParticles) // Slices for randomly generated points

	particle.GeneratePoints(p, initPoints, nParticles)

	_, err = fmt.Fprintf(datafile, "%d %d %d\n", nParticles, nIters, 0)
	utils.Check(err)

	for iter := 0; iter < nIters; iter++ {
		if logging {
			fmt.Printf("Running iteration %d\n", iter+1)
		}

		max := particle.FindBoundsWorkStealing(p, nParticles, numThreads, workQueues, completed)
		workQueues, completed = resetQueues(numThreads, nParticles)

		var root *tree.QuadTree = tree.CreateNode(nil, -1*max, -1*max, max, max) // Create root of the tree

		/* Add points to tree */
		tree.AddParticlesWorkStealing(p, root, nParticles, numThreads, workQueues, completed)
		workQueues, completed = resetQueues(numThreads, nParticles)
		if logging {
			for i := 0; i < nParticles; i++ {
				fmt.Fprintf(datafile, "%f %f \n", p[i].X, p[i].Y)
			}
		}

		// Compute center of mass for the tree
		tree.ComputeCenterOfMass(root)

		tree.CalcTreeForceAndUpdatePosWorkStealing(p, root, numThreads, nParticles, workQueues, completed)
		workQueues, completed = resetQueues(numThreads, nParticles)
	}
	dur := time.Since(start)
	fmt.Printf("%f\n", dur.Seconds())
}
