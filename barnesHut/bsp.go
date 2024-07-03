package barnesHut

import (
	"fmt"
	"os"
	"nbody_go/particle"
	"nbody_go/tree"
	"nbody_go/utils"
	"time"
)

// Code for Barnes-Hut nBody problem using BSP pattern
func Bsp() {
	start := time.Now()
	datafile, err := os.Create("particles_bsp.dat") // Output file for particle positions
	utils.Check(err)
	defer datafile.Close()

	config := utils.GetParams()
	nParticles := config.NParticles
	nIters := config.NIters
	numThreads := config.NumThreads
	logging := config.Logging
	initPoints := config.InitPoints

	var p []*particle.Particle = make([]*particle.Particle, nParticles) // Slices for randomly generated points

	particle.GeneratePoints(p, initPoints, nParticles)

	_, err = fmt.Fprintf(datafile, "%d %d %d\n", nParticles, nIters, 0)
	utils.Check(err)

	for iter := 0; iter < nIters; iter++ {
		if logging {
			fmt.Printf("Running iteration %d\n", iter+1)
		}

		max := particle.FindBoundsParallel(p, nParticles, numThreads)

		var root *tree.QuadTree = tree.CreateNode(nil, -1*max, -1*max, max, max) // Create root of the tree

		/* Add points to tree */
		tree.AddParticlesParallel(p, root, nParticles, numThreads)
		if logging {
			for i := 0; i < nParticles; i++ {
				fmt.Fprintf(datafile, "%f %f \n", p[i].X, p[i].Y)
			}
		}

		// Compute center of mass for the tree
		tree.ComputeCenterOfMass(root)

		tree.CalcTreeForceAndUpdatePosParallel(p, root, numThreads, nParticles)
	}
	dur := time.Since(start)
	fmt.Printf("%f\n", dur.Seconds())
}
