package main

import (
	"fmt"
	"math"
	"os"
	"proj3/particle"
	"proj3/tree"
	"strconv"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// Pass number of files as an argument, default is one
func main() {
	datafile, err := os.Create("particles.dat") // Output file for particle positions
	check(err)
	defer datafile.Close()

	var nParticles int = 3000    // Default number of particles
	const dt float64 = 0.01      // Time step
	const nIters int = 200       // Number of steps in simulation
	const theta float64 = 0.7071 // Approximation constant for using center of mass

	if len(os.Args) > 1 {
		nParticles, err = strconv.Atoi(os.Args[1])
		check(err)
	}
	fmt.Printf("Number of particles set to %d\n", nParticles)

	var p []*particle.Particle = make([]*particle.Particle, nParticles) // Arrays for randomly generated points
	particle.InitialiseRandomPoints(p, nParticles)                      // Init position and velocity data

	_, err = fmt.Fprintf(datafile, "%d %d %d\n", nParticles, nIters, 0)
	check(err)

	for iter := 0; iter < nIters; iter++ {
		fmt.Printf("Running iteration %d\n", iter+1)

		var maxX float64 = p[0].X
		var minX float64 = p[0].X
		var maxY float64 = p[0].Y
		var minY float64 = p[0].Y
		for i := 1; i < nParticles; i++ {
			maxX = math.Max(p[i].X, maxX)
			minX = math.Min(p[i].X, minX)
			maxY = math.Max(p[i].Y, maxY)
			minY = math.Min(p[i].Y, minY)
		}
		var max float64 = math.Max(math.Abs(maxX), math.Abs(minX))
		max = math.Max(math.Abs(minY), math.Abs(max))
		max = math.Max(math.Abs(maxY), math.Abs(max))

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
			tree.CalcTreeForce(p[i], root, theta, dt)
		}

		// Update the position
		for i := 0; i < nParticles; i++ {
			p[i].X += p[i].Vx * dt
			p[i].Y += p[i].Vy * dt
		}
	}
}
