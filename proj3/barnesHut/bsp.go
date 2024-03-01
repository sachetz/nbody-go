package barnesHut

import (
	"fmt"
	"os"
	"proj3/particle"
	"proj3/tree"
	"proj3/utils"
)

// Code for Barnes-Hut nBody problem using BSP pattern
func Bsp() {
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
		for i := 0; i < nParticles; i++ {
			tree.AddParticleToTree(p[i], root)
			fmt.Fprintf(datafile, "%f %f \n", p[i].X, p[i].Y)
		}

		// Compute center of mass for the tree
		tree.ComputeCenterOfMass(root)

		// Calculate force on each particle
		for i := 0; i < nParticles; i++ {
			tree.CalcTreeForce(p[i], root, utils.Theta, utils.Dt)
		}

		particle.UpdatePosParallel(p, nParticles, numThreads)
	}
}
