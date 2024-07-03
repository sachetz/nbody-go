#!/bin/bash
#
#SBATCH --mail-user=zodesachet@gmail.com
#SBATCH --mail-type=ALL
#SBATCH --job-name=nbody_go_graph
#SBATCH --output=./slurm/out/%j.%N.stdout
#SBATCH --error=./slurm/out/%j.%N.stderr
#SBATCH --chdir=/home/sachetz/nbody_go/benchmarks
#SBATCH --partition=debug 
#SBATCH --nodes=1
#SBATCH --ntasks=1
#SBATCH --cpus-per-task=16
#SBATCH --mem-per-cpu=900
#SBATCH --exclusive
#SBATCH --time=3:59:59


module load golang/1.19
python benchmark.py
