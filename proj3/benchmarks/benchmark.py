import subprocess
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from collections import defaultdict

if __name__ == "__main__":
    
    # Open benchmark file in write mode
    benchmarkFile = open("benchmark.txt", 'w', buffering=1)

    # Process for each input type
    for inputType in ["random", "circle", "skewed"]:

        benchmarkFile.write(f"--------------------------------------------------\n")
        benchmarkFile.write(f"Processing {inputType} input type.\n")

        for nIters in [1, 10, 100, 1000]:

            benchmarkFile.write(f"--------------------------------------------------\n")
            benchmarkFile.write(f"Processing {nIters} iterations.\n")

            # Process for each test size
            for testSize in [10, 30, 100, 300, 1000, 3000, 10000]:
                
                benchmarkFile.write(f"--------------------------------------------------\n")
                benchmarkFile.write(f"Processing {testSize} test size.\n")
                
                # Run the process in sequential mode and get the times
                sequentialProcess = subprocess.Popen([
                    "go", "run", "../nBody.go", "s", str(testSize), str(nIters), "1", "false", inputType
                ], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
                sequentialTime = float(sequentialProcess.stdout.readline().decode()[:-1])

                benchmarkFile.write(f"Total sequential time: {str(sequentialTime)}\n")
                benchmarkFile.write("\n")

                for mode in ["bsp", "ws"]:

                    benchmarkFile.write(f"Running for mode {mode}\n\n")

                    times = defaultdict(list)
                    threads = [2, 4, 6, 8, 12]

                    # For each no of threads
                    for nThreads in threads:

                        benchmarkFile.write(f"For nThreads {nThreads}\n")

                        for experiment in range(1, 4):

                            # Run the process and get the time
                            parallelProcess = subprocess.Popen([
                                "go", "run", "../nBody.go", mode, str(testSize), str(nIters), str(nThreads), "false", inputType
                            ], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
                            parallelTime = float(parallelProcess.stdout.readline().decode()[:-1])
                            
                            # Calculate the speedup
                            #benchmarkFile.write(f"In experiment {experiment}, time taken for execution in parallel: {parallelTime}\n")
                            benchmarkFile.write(f"Time taken for execution in parallel: {parallelTime}\n")
                            times[nThreads] += [parallelTime]

                    benchmarkFile.write("\n")

                    speedups = [sequentialTime/min(times[key]) for key in sorted(times.keys())]

                    # Plot the graph
                    plt.figure(f"{inputType}-{mode}-{nIters}")
                    plt.plot(threads, speedups, label=testSize)
                    plt.legend(loc='upper left')
                    plt.xlabel("Number of Threads")
                    plt.ylabel("Speedup")
                    plt.savefig(f"speedup-benchmark-{inputType}-{mode}-{nIters}")
            plt.close(f"{inputType}-{mode}-{nIters}")        
    
    benchmarkFile.close()
