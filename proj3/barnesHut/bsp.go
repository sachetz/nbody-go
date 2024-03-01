package barnesHut

import (
	"fmt"
	"os"
	"proj3/particle"
	"proj3/tree"
	"proj3/utils"
	"time"
)

// Code for Barnes-Hut nBody problem using BSP pattern
func Bsp() {
	start := time.Now()
	datafile, err := os.Create("benchmarks/particles_bsp.dat") // Output file for particle positions
	utils.Check(err)
	defer datafile.Close()

	nParticles, nIters, numThreads := utils.GetParams()

	var p []*particle.Particle = make([]*particle.Particle, nParticles)     // Slices for randomly generated points
	particle.InitialiseParticlesInCircleParallel(p, nParticles, numThreads) // Init position and velocity data

	_, err = fmt.Fprintf(datafile, "%d %d %d\n", nParticles, nIters, 0)
	utils.Check(err)

	for iter := 0; iter < nIters; iter++ {
		fmt.Printf("Running iteration %d\n", iter+1)

		max := particle.FindBoundsParallel(p, nParticles, numThreads)

		var root *tree.QuadTree = tree.CreateNode(nil, -1*max, -1*max, max, max) // Create root of the tree

		/* Add points to tree */
		tree.AddParticlesParallel(p, root, nParticles, numThreads)

		// Compute center of mass for the tree
		// TODO - parallelize
		tree.ComputeCenterOfMass(root)

		tree.CalcTreeForceParallel(p, root, numThreads, nParticles)

		particle.UpdatePosParallel(p, nParticles, numThreads)
	}
	dur := time.Since(start)
	fmt.Printf("Time taken %f", dur.Seconds())
}
