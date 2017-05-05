# Diffusion Limited Aggregation 

Simulates diffusion limited aggregation on a lattice. Once a simulation has been run/loaded it can be drawn, or analysed to find it's Hausdorff dimensions. This can be compared under regeims of varying numbers of spatial dimensions (can go arbitrarilly high, though increases in time compelxity) or with a varying chance for adhesion at any particular valid site.

This application was written as part of my scientific computing module at the university of nottingham.

### Use
the application will launch a command line interface. Here you will have the option of several commands.

*run* _-npoints=5000 -dimensions=2 -seed=1 -runs=32 -pstick=1.0 -load=false_
This runs the application with the default parameters shown here. Runs are kept in memory, and are overwritten by any future run commnnd. The flags effect it as follows:
*_-npoints_* changes the number points that will be included in the final aggregate.
*_-dimensions_* changes the number of spatial dimensions in which the simulation will run. Can be set arbitrarally high.
*_-seed_* changes the random number seed used to generate new runs.
*_-runs_* causes the simulation to run multiple times with different seeds (genrerated from the passed seed).
*_-pstick_* changes the probability of a particle sticking at any valid site.
*_-load_* If a simulation has been run with the same parameters it will try to load it from the disk to save time.

*save*
Saves all stored runs to the disk so they can be loaded in later.

*varydimension* _-npoints=5000 -start=2 -stop=6 -step=1 -seed=1 -runs=32 -pstick=1 -load=false_
Runs the simulation multiple times, varying the number of spatial dimensions as it does so. Runs are kept in memory, and are overwritten by any future run commnnd. The flags effect it as follows:
*_-npoints_* changes the number points that will be included in the final aggregate.
*_-start_* The lowest number of dimensions to simulate.
*_-stop_* The highest number of dimensions to simulate.
*_step_* The number of dimensions to increase by in the range.
*_-seed_* changes the random number seed used to generate new runs.
*_-runs_* causes the simulation to run multiple times with different seeds (genrerated from the passed seed).
*_-pstick_* changes the probability of a particle sticking at any valid site.
*_-load_* If a simulation has been run with the same parameters it will try to load it from the disk to save time.

*varysticking* _-npoints=5000 -dimensions=2 -seed=1 -runs=32 -load=false -step=0.1 0.1 1.0_
Runs the simulation multiple times, varying the sticking probability as it does so. Runs are kept in memory, and are overwritten by any future run commnnd. The flags effect it as follows:
*_-npoints_* changes the number points that will be included in the final aggregate.
*_-dimensions_* changes the number of spatial dimensions in which the simulation will run. Can be set arbitrarally high.
*_-seed_* changes the random number seed used to generate new runs.
*_-runs_* causes the simulation to run multiple times with different seeds (genrerated from the passed seed).
*_-load_* If a simulation has been run with the same parameters it will try to load it from the disk to save time.
*_-step_* Sets the gap between sticking probabilities specified in the tail. If this is *>0* it will use point between the two point specified in the tail. If instead it is *0* it will use a whole list of numbers from the tail. (i.e _-step=0 0.1 0.4 0.5 0.7 1_ would run the simulation for sticking probabilities of 0.1, 0.4, 0.5, 0.7, and 1)
