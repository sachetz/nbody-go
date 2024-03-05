import subprocess
import matplotlib
matplotlib.use("Agg")
import matplotlib.pyplot as plt
from collections import defaultdict

if __name__ == "__main__":
    
    # Open benchmark file in write mode
    benchmarkFile = open("benchmark.txt", 'w', buffering=1)

    # Process for each test size
    for testSize in [10, 30, 100, 300, 1000, 3000, 10000]:
        
        benchmarkFile.write(f"--------------------------------------------------\n")
        benchmarkFile.write(f"Processing {testSize} test size.\n")
        
        # Run the process in sequential mode and get the times
        sequentialProcess = subprocess.Popen([
            "go", "run", "../nBody.go", "s", str(testSize), "200", "1", "false", "random"
        ], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        sequentialTime = float(sequentialProcess.stdout.readline().decode()[:-1])

        benchmarkFile.write(f"Total sequential time: {str(sequentialTime)}\n")
        benchmarkFile.write("\n")

        for mode in ["bsp", "ws"]:

            times = defaultdict(list)
            threads = [2, 4, 6, 8, 12]

            # For each no of threads
            for nThreads in threads:

                benchmarkFile.write(f"For nThreads {nThreads}\n")

                for experiment in range(1, 6):

                    # Run the process and get the time
                    parallelProcess = subprocess.Popen([
                        "go", "run", "../nBody.go", mode, str(testSize), "200", str(nThreads), "false", "random"
                    ], stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
                    parallelTime = float(parallelProcess.stdout.readline().decode()[:-1])
                    
                    # Calculate the speedup
                    benchmarkFile.write(f"In experiment {experiment}, time taken for execution in parallel: {parallelTime}\n")
                    times[nThreads] = times[nThreads] + [parallelTime]

            benchmarkFile.write("\n")

            speedups = [sequentialTime/(sum(times[key])/5) for key in sorted(times.keys())]

            # Plot the graph
            plt.figure(f"{mode}")
            plt.plot(threads, speedups, label=testSize)
            plt.legend(loc='upper left')
            plt.xlabel("Number of Threads")
            plt.ylabel("Speedup")
            plt.savefig(f"speedup-benchmark-{mode}")
    
    benchmarkFile.close()
