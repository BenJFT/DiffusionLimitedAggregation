package plot

import (
	"os"
	"bufio"

	"github.com/skratchdot/open-golang/open"
)

func ExamplePlot() {
	var (
		x = []float64{1, 2, 3, 4, 5}
		y = []float64{1, 4, 9, 16, 25}
		plot = Plot(NewScatter(x, y))
		name = "example.svg"
		file, err = os.Create(name)
		writer = bufio.NewWriter(file)
	)

	if err != nil {
		panic(err)
	}

	_, err = writer.Write([]byte(plot))
	if err != nil {
		panic(err)
	}
	open.Run(name)

	println("Hello")
	// Output:
}