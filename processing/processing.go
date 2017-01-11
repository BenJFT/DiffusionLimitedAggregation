/*
Author: Benedict Thompson
*/

package processing

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"math"
	"math/rand"
	"os"

	"github.com/Benjft/DiffusionLimitedAggregation/processing/aggregation"
	"github.com/Benjft/DiffusionLimitedAggregation/processing/analysis"
	"github.com/Benjft/DiffusionLimitedAggregation/util"
	"github.com/Benjft/DiffusionLimitedAggregation/util/drawing/svg"
	"github.com/Benjft/DiffusionLimitedAggregation/util/drawing/xyz"

	"github.com/gonum/plot"
	"github.com/gonum/plot/plotter"
	"github.com/gonum/plot/plotutil"
	"github.com/gonum/plot/vg"

	"github.com/skratchdot/open-golang/open"
)

func init() {
	os.Mkdir("out", os.ModeDir)
	os.Mkdir("out\\plot", os.ModeDir)
	os.Mkdir("out\\saves", os.ModeDir)

	gob.Register(RunState{})
}

// Structure used to store all info about a previous run, including it's parameters and results
type RunState struct {
	NPoints, NDimension, NRuns, Seed int64
	Sticking                         float64
	Points                           [][]aggregation.Point
}

// The most recent successful run or loaded run from save. All opps are done on the contents of this
var loadedRun = RunState{}

func Run(nPoints, nRuns, seed, nDimension int64, sticking float64) {
	var (
		channel chan []aggregation.Point = make(chan []aggregation.Point)
		points  [][]aggregation.Point    = make([][]aggregation.Point, nRuns)
	)

	rand.Seed(seed)
	for i := int64(0); i < nRuns; i++ {
		// sets each run to go concurrently, sending their results over 'channel'. This improves the run time
		// when there are a large number of runs
		go func (seed int64) {
			channel <- aggregation.RunNew(nPoints, seed, nDimension, sticking)
		} (rand.Int63())
	}

	// waits for each output
	for i := int64(0); i < nRuns; i++ {
		points[i] = <-channel
	}

	// none of the runs failed update the loaded run
	loadedRun = RunState{
		NPoints:    nPoints,
		NDimension: nDimension,
		NRuns:      nRuns,
		Seed:       seed,
		Sticking:   sticking,
		Points:     points,
	}
}

// Saves the 3D aggregate as an 'xyz' file. These can be opened using many 3D graphics editors, I know that a blender 3D
// extension exists for importing these. 'xyz' was the chosen file type as i don't have to deal with geometry this way
func draw3D(state []aggregation.Point, title string) {
	var name string = fmt.Sprintf("out\\plot\\%s.xyz", title)

	// attempts to create a new file with the provided name
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// ensures the file is closed when the function returns or program stops
	defer file.Close()

	//handles writing the data to the opened file
	writer := bufio.NewWriter(file)
	// gets the plaintext formatting of the state as an xyz file
	str := xyz.DrawAggregate(state)
	// writes the bytes representing the string to the file
	_, err = writer.Write([]byte(str))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// If nothing went wrong, open the file using the system default
	err = open.Run(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Saves the 2D aggregate as an 'svg' file. This can be opened in most web browsers or vector graphics editors.
func draw2D(state []aggregation.Point, title string) {
	name := fmt.Sprintf("out\\plot\\%s.svg", title)

	// attempts to create a new file with the provided name
	file, err := os.Create(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// ensures the file is closed when the function returns or program stops
	defer file.Close()

	//handles writing the data to the opened file
	writer := bufio.NewWriter(file)
	// gets the plaintext formatting of the state as an xyz file
	str := svg.DrawAggregate(state)
	// writes the bytes representing the string to the file
	_, err = writer.Write([]byte(str))
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// If nothing went wrong, open the file using the system default
	err = open.Run(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}

// Called to render the latest run (given that the run is in only 2 or 3 dimensions)
func Draw(title string) {
	// Create a default formatted title using the run parameters if no title was provided
	if len(title) == 0 {
		title = fmt.Sprintf("aggregate-n%d-seed%d-dims%d-stick%f",
			loadedRun.NPoints,
			loadedRun.Seed,
			loadedRun.NDimension,
			loadedRun.Sticking)
	}

	// draw using the appropreate function for the number of dimensions
	if n := loadedRun.NDimension; n == 2 {
		for run, state := range loadedRun.Points {
			runtitle := fmt.Sprintf("%s-run%d", title, run)
			go draw2D(state, runtitle)
		}
	} else if n == 3 {
		for run, state := range loadedRun.Points {
			runtitle := fmt.Sprintf("%s-run%d", title, run)
			go draw3D(state, runtitle)
		}
	} else {
		fmt.Println("Can only draw 2D and 3D lattices")
	}
}

// Saves #loadedRun using gob formatting so that it can bea easily loaded at a later time
func Save(title string) {
	// Create a default formatted title using the run parameters if no title was provided
	if title == "" {
		title = fmt.Sprintf("save-n%d-seed%d-dims%d-stick%f-runs%d",
			loadedRun.NPoints,
			loadedRun.Seed,
			loadedRun.NDimension,
			loadedRun.Sticking,
			loadedRun.NRuns)
	}
	path := fmt.Sprintf("out\\saves\\%s.save", title)

	// Create the file at the given location
	file, err := os.Create(path)
	if err != nil {
		fmt.Println(err.Error())

	}
	// Ensures the file is closed when the function quits
	defer file.Close()

	// write the state to a gob in the given file
	encoder := gob.NewEncoder(file)
	err = encoder.Encode(loadedRun)
	if err != nil {
		fmt.Println(err.Error())
	}
}

// Loads the named run into loadedRun
func Load(title string) {
	path := fmt.Sprintf("out\\saves\\%s.save", title)

	// attempt to open the file fr reading
	file, err := os.Open(path)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// Ensures the file is closed when the function quits
	defer file.Close()

	// read the gob in the file into a new Run struct
	decoder := gob.NewDecoder(file)
	var tmpRun RunState
	err = decoder.Decode(&tmpRun)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// If nothing went wrong set the make this the active Run
	loadedRun = tmpRun
}

// MeanAndConf95 returns the mean
// and the magnitude of the 95% confidence
// interval on the mean as low and high
// error values.
//
// MeanAndConf95 may be used as
// the f argument to NewErrorPoints.
func meanAndConf95(vls []float64) (mean, lowerr, higherr float64) {
	n := float64(len(vls))

	sum := 0.0
	for _, v := range vls {
		sum += v
	}
	mean = sum / n

	var stdev float64 = 0.0
	if n > 1 {
		sum = 0.0
		for _, v := range vls {
			diff := v - mean
			sum += diff * diff
		}
		stdev = math.Sqrt(sum / (n - 1.5))
	} else {
		stdev = 0
	}

	conf := 1.96 * stdev / math.Sqrt(n)
	return mean, conf, conf
}

// uses the approximate radius of the aggregate to plot log(N)/log(R) and find the fractal dimensions
func Radii(title string) {
	// open a channel and wait to receive the radius calculations for each run in the current loaded state
	var channel chan []float64 = make(chan []float64)
	for _, run := range loadedRun.Points {
		go func(run []aggregation.Point) {
			channel <- analysis.GyrationRadii(run)
		}(run)
	}

	radii := make([][]float64, loadedRun.NRuns)
	for i := range loadedRun.Points {
		radii[i] = <-channel
	}

	// flip the arrays so that each row contains the radius for the same number of particles (makes plotting
	// scatters and calculating error easier)
	radii = util.Transpose(radii)

	pts := make([]plotter.XYer, len(radii)-1)

	// make arrays of points marking X=log(R) Y=log(N). Skips the first point as R(N=1) = 0
	for i, r := range radii[1:] {
		xys := make(plotter.XYs, len(r))
		N := math.Log10(float64(i + 2))
		for j, y := range r {
			xys[j].X = N
			xys[j].Y = math.Log10(y)
		}
		pts[i] = xys
	}

	// create a new plot of the data
	plt, err := plot.New()
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	plt.X.Label.Text = "Log(N)"
	plt.Y.Label.Text = "Log(R)"



	// plot the average after each point is added and find the error with 95% confidence.
	mean95, err := plotutil.NewErrorPoints(meanAndConf95, pts...)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// add these to the plot
	err = plotutil.AddScatters(plt, mean95)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	err = plotutil.AddYErrorBars(plt, mean95)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// calculate the regression coefficients and their respective errors
	regressionPoints := plotutil.ErrorPoints{
		XYs:mean95.XYs[len(mean95.XYs)/20:],
		XErrors:mean95.XErrors[len(mean95.XYs)/20:],
		YErrors:mean95.YErrors[len(mean95.XYs)/20:],
	}
	var a, b, eb float64
	if loadedRun.NRuns > 1 {
		a, b, _, eb = util.WeightedLeastSquares(regressionPoints)
	} else {
		a, b, _, eb = util.LeastSquares(regressionPoints.XYs)
	}

	eb = eb/(b*b - eb*eb)
	// Print the approximation of the fractal dimensions from the regression
	fmt.Printf("D = %.3f \u00B1 %.3f\n", 1/b, eb*1.96) // multiply by 1.96 for 95% confidence interval
	// add and label the Least Squares fit
	label := fmt.Sprintf("y = %.3f + %.3fx", a, b)
	fit := plotter.NewFunction(func(x float64) float64 { return a + b*x })
	plt.Add(fit)
	plt.Legend.Add(label, fit)

	// Save and show the plot
	if title == "" {
		title = fmt.Sprintf("radii-n%d-seed%d-dims%d-stick%f-runs%d",
			loadedRun.NPoints,
			loadedRun.Seed,
			loadedRun.NDimension,
			loadedRun.Sticking,
			loadedRun.NRuns)
	}

	name := fmt.Sprintf("out\\plot\\%s.svg", title)
	plt.Save(8*vg.Inch, 5.66*vg.Inch, name)
	err = open.Run(name)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
}
