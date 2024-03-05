#!/bin/bash
#
#SBATCH --mail-user=sachetz@cs.uchicago.edu
#SBATCH --mail-type=ALL
#SBATCH --job-name=proj3_graph
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/sachetz/ParallelProgramming/project-3-sachetz/proj3/benchmarks
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=3:00:00


module load golang/1.19
python benchmark.py
