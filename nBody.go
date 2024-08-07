package main

import (
	"fmt"
	"nbody_go/barnesHut"
	"os"
)

const usage string = "main() takes 6 arguments:\n" +
	"\t" + "1. Mode - this can be:\n" +
	"\t\t" + "s (sequential) - default\n" +
	"\t\t" + "bsp (parallel with BSP)\n" +
	"\t\t" + "ws (parallel with work stealing)\n" +
	"\t" + "2. nPoints - number of points, default is 3000\n" +
	"\t" + "3. nIters - number of iterations, default is 200\n" +
	"\t" + "4. numThreads - number of threads/goroutines, default is 8\n" +
	"\t" + "5. logging - display logs/create output file, default is true\n" +
	"\t" + "6. initPoints - this can be:\n" +
	"\t\t" + "random - default\n" +
	"\t\t" + "circle\n" +
	"\t\t" + "skewed"

func main() {
	mode := "s"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}
	if mode == "s" {
		barnesHut.Sequential()
	} else if mode == "bsp" {
		barnesHut.Bsp()
	} else if mode == "ws" {
		barnesHut.Ws()
	} else {
		fmt.Println("Usage: " + usage)
	}
}
