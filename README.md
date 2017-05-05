# Diffusion Limited Aggregation 

Simulates diffusion limited aggregation on a lattice. Once a simulation has been run/loaded it can be drawn, or analyze to find it's Hausdorff dimensions. This can be compared under regeims of varying numbers of spatial dimensions (can go arbitrarilly high, though increases in time compelxity) or with a varying chance for adhesion at any particular valid site.

This application was written as part of my scientific computing module at the university of Nottingham.

### Use
the application will launch a command line interface. Here you will have the option of several commands.
___
`run -npoints=5000 -dimensions=2 -seed=1 -runs=32 -pstick=1.0 -load=false`  
This runs the application with the default parameters shown here. Runs are kept in memory, and are overwritten by any future run commnnd. The flags effect it as follows:

`-npoints` changes the number points that will be included in the final aggregate.  
`-dimensions` changes the number of spatial dimensions in which the simulation will run. Can be set arbitrarily high.  
`-seed` changes the random number seed used to generate new runs.  
`-runs` causes the simulation to run multiple times with different seeds (generated from the passed seed).  
`-pstick` changes the probability of a particle sticking at any valid site.  
`-load` If a simulation has been run with the same parameters it will try to load it from the disk to save time.  
___
`save`
Saves all stored runs to the disk so they can be loaded in later.
___
`varydimension -npoints=5000 -start=2 -stop=6 -step=1 -seed=1 -runs=32 -pstick=1 -load=false`  
Runs the simulation multiple times, varying the number of spatial dimensions as it does so. Runs are kept in memory, and are overwritten by any future run commnnd. The flags effect it as follows:

`-npoints` changes the number points that will be included in the final aggregate.  
`-start` The lowest number of dimensions to simulate.  
`-stop` The highest number of dimensions to simulate.  
`step` The number of dimensions to increase by in the range.  
`-seed` changes the random number seed used to generate new runs.  
`-runs` causes the simulation to run multiple times with different seeds (generated from the passed seed).  
`-pstick` changes the probability of a particle sticking at any valid site.  
`-load` If a simulation has been run with the same parameters it will try to load it from the disk to save time.  
___
`varysticking -npoints=5000 -dimensions=2 -seed=1 -runs=32 -load=false -step=0.1 0.1 1.0`  
Runs the simulation multiple times, varying the sticking probability as it does so. Runs are kept in memory, and are overwritten by any future run commnnd. The flags effect it as follows:

`-npoints` changes the number points that will be included in the final aggregate.  
`-dimensions` changes the number of spatial dimensions in which the simulation will run. Can be set arbitrarily high.  
`-seed` changes the random number seed used to generate new runs.  
`-runs` causes the simulation to run multiple times with different seeds (generated from the passed seed).  
`-load` If a simulation has been run with the same parameters it will try to load it from the disk to save time.  
`-step` Sets the gap between sticking probabilities specified in the tail. If this is *>0* it will use point between the two point specified in the tail.
If instead it is *0* it will use a whole list of numbers from the tail. (i.e _-step=0 0.1 0.4 0.5 0.7 1_ would run the simulation for sticking probabilities of 0.1, 0.4, 0.5, 0.7, and 1)
___
`growth -hidecurves=false -raw=false -hidetrend=false`  
Gathers the data from the last set of runs and uses the information from similar runs to find and their Hausdorff dimensions based of the growth rate.

`-hidecurves` Sets if the growth rate of each set of similar simulations should be drawn.  
`-raw` Sets if the raw data for each plot should be output to a csv file for use externally.
`-hidetrend` Sets if the relationship between the growth rate based Hausdorff dimension should be plotted
___
`draw`  
Draws any loaded 2D aggregates as svg files with hue showing the order in which particles were added to the aggregate.
___
### Example Uses

`varysticking -npoints=10000 -runs=1 -load=true -step=0 0.001 0.01 0.1 1.0 save draw`  
Runs the simulation once for each of sticking values 0.001 0.01 0.1 1.0 then draws each to an svg file.

`varydimensions -stop=7 -load=true growth -hidecurves=true`  
Runs the simulation 32 times for each of D=2,3,4,5,6,7 and then plots the overall trend in the growth rate.
