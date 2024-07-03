# Introduction

This project aims to implement an efficient algorithm to simulate the N-Body problem using the Barnes-Hut algorithm. The project is designed and implemented in Go, with a major focus on parallelism and improving the performance of the simulation - mainly using 2 parallelisation techniques, BSP and work-stealing queues. The performance of the algorithm is tested across different sizes and input types, with the initial placement of the points generated to be random, in a circle and skewed to form a cluster with sparse outliers, the number of points tested in the range of 10 to 10000, and the number of iterations set between 10 and 1000. The project saves the location of each point after every time step, so as to visualize the movement of the points, the results of which are added to the project benchmarks.

# Parallelized Implementations

For the parallel implementations, each stage of the algorithm is executed concurrently across 𝑇 threads, with each thread working on 𝑁/𝑇 particles.

Two parallel implementations are compared. First, a Bulk Synchronous Pattern (BSP) based approach is implemented, which is a great fit for this algorithm due to the nature of distinct supersteps. Barriers (built using condition variables) were used to separate logical supersteps to synchronize the execution of the different threads. No thread is allowed to proceed with further computation until all threads have reached that barrier, thus ensuring data consistency between the supersteps. This is an ideal solution for this system, as the Barnes-Hut algorithm is inherently divided into clear supersteps.

The second implementation involved the use of work-stealing queues. The skewness and clustering of the initial data points hinted that distributing the work in terms of number of points might not be an ideal approach, as in cases of clusters, there would be an imbalance in the work assigned to each thread. Work-stealing queues would effectively solve this problem by ensuring that threads have an opportunity to steal other threads’ work, resulting in an effective increase in the overall parallelism. The work-stealing queue was implemented using a double-ended queue, where the bottom of the queue is accessible to the owner thread (and supports both the pop and push operations), whereas other steal can only pop any pending jobs from the top of the queue.

# Instructions to Run

### Replicating the experiment and graph generation

To replicate the experiment, the user must `cd` to the `/<path-to-local>/benchmarks` directory, update the graph.sh file with the correct parameters for slurm, and execute the command: sbatch graph.sh.

The process then calls a python script that handles the execution of the process to replicate the experiment. The python script performs the following operations:
1. Calls the process in sequential mode and records the time
2. Calls the process with combinations of various parameters - pattern of initializing points,
number of iterations/time steps, number of points, number of threads and the parallel
mode (BSP or work stealing)
3. Calculates the observed speedup from the execution times
4. Generates graphs between the number of threads and speedup comparing different
parameter combinations.
The slurm logs are generated in the slurm/out folder, and the benchmarks are recorded in the benchmarks.txt file, both located in the `/<path-to-local>/benchmarks` folder.

### Running specific versions of the program

To run specific versions of the program, the user must cd to the `/<path-to-local>` directory and run the nBody.go file as `go run nBody.go <args>`. The process supports the following arguments, in order:
1. Mode - this takes values s (sequential), bsp (barrier based parallelisation), ws (work stealing based parallelisation). The default is sequential.
2. nPoints - This is the number of generated points. The default is 3000.
3. nIters - This is the number of iterations/time steps. The default is 200.
4. numThreads - Number of threads for the parallel implementations. The default is 8.
5. Logging - This indicates whether logging or file write is enabled. The default is true.
6. initPoints - This is the initial setup of points. This takes values random (which initializes points randomly in a subspace), circle (which initializes points in a circle), and skewed (which initializes points in a cluster and adds some sparse outliers).

If logging is set to true, the process writes the generated location of the points after every time step into a file `particles_<mode>.py`. These change of locations can be converted into a gif using the plot_particles.py file in the benchmarks folder, by providing `<mode>` as an argument. The result of two such sample runs with 16 points and 200 iterations initialized in a circle is included in the repository.