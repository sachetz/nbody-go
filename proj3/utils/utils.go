package utils

import (
	"fmt"
	"os"
	"strconv"
)

const Dt float64 = 0.01      // Time step
const Theta float64 = 0.7071 // Approximation constant for using center of mass

func Check(e error) {
	if e != nil {
		panic(e)
	}
}

func GetParams() (int, int, int) {
	var nParticles int = 3000 // Default number of particles
	var nIters int = 200      // Number of steps in simulation
	var numThreads int = 8    // Default number of threads
	var err error

	if len(os.Args) > 2 {
		nParticles, err = strconv.Atoi(os.Args[2])
		Check(err)
		fmt.Printf("Number of particles set to %d\n", nParticles)
	}
	if len(os.Args) > 3 {
		nIters, err = strconv.Atoi(os.Args[3])
		Check(err)
		fmt.Printf("Number of iterations set to %d\n", nIters)
	}
	if len(os.Args) > 4 {
		numThreads, err = strconv.Atoi(os.Args[4])
		Check(err)
		fmt.Printf("Number of threads set to %d\n", 1)
	}
	return nParticles, nIters, numThreads
}
