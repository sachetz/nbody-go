package barnesHut

import (
	"fmt"
	"os"
	"proj3/particle"
	"proj3/tree"
	"proj3/utils"
	"time"
)

// Sequential code for Barnes-Hut nBody problem
func Sequential() {
	start := time.Now()
	datafile, err := os.Create("benchmarks/particles_s.dat") // Output file for particle positions
	utils.Check(err)
	defer datafile.Close()

	nParticles, nIters, _, logging := utils.GetParams()

	var p []*particle.Particle = make([]*particle.Particle, nParticles)          // Slices for randomly generated points
	particle.InitialiseParticlesInCircleSequential(p, 0, nParticles, nParticles) // Init position and velocity data
	//particle.InitialiseRandomPointsSequential(p, 0, nParticles)

	_, err = fmt.Fprintf(datafile, "%d %d %d\n", nParticles, nIters, 0)
	utils.Check(err)

	for iter := 0; iter < nIters; iter++ {
		if logging {
			fmt.Printf("Running iteration %d\n", iter+1)
		}

		max := particle.FindBoundsSequential(p, 0, nParticles)

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

		// Update the position
		particle.UpdatePosSequential(p, 0, nParticles)
	}
	dur := time.Since(start)
	fmt.Printf("%f\n", dur.Seconds())
}
