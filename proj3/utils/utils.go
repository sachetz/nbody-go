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

func GetParams() (int, int, int, bool) {
	var nParticles int = 3000 // Default number of particles
	var nIters int = 200      // Number of steps in simulation
	var numThreads int = 8    // Default number of threads
	var logging bool = true
	var err error

	if len(os.Args) > 2 {
		nParticles, err = strconv.Atoi(os.Args[2])
		Check(err)
	}
	if len(os.Args) > 3 {
		nIters, err = strconv.Atoi(os.Args[3])
		Check(err)
	}
	if len(os.Args) > 4 {
		numThreads, err = strconv.Atoi(os.Args[4])
		Check(err)
	}
	if len(os.Args) > 5 {
		logging, err = strconv.ParseBool(os.Args[5])
		Check(err)
	}
	if logging {
		fmt.Printf("Number of particles set to %d\n", nParticles)
		fmt.Printf("Number of iterations set to %d\n", nIters)
		fmt.Printf("Number of threads set to %d\n", numThreads)
		fmt.Printf("Logging set to %t\n", logging)
	}
	return nParticles, nIters, numThreads, logging
}
