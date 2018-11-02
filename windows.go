package main

import (
	"math"

	"bitbucket.org/rhagenson/swsc/nexus"
	"github.com/pkg/errors"
)

type Window [2]int

func (w *Window) Start() int {
	return w[0]
}

func (w *Window) Stop() int {
	return w[1]
}

func getAllWindows(uceAln nexus.Alignment, minWin int) []Window {
	windows, _ := generateWindows(uceAln.Len(), minWin)
	return windows
}

// generateWindows produces windows of at least a minimum size given a total length
// Windows must be:
//   1) at least minimum window from the start of the UCE (ie, first start at minimum+1)
//   2) at least minimum window from the end of the UCE (ie, last end at length-minimum+1)
//   3) at least minimum window in length (ie, window{start, end)})
// TODO: Preallocate windows array. Tried my hand at this once and hit a mental block.
func generateWindows(length, min int) ([]Window, error) {
	if min > length/3 {
		return nil, errors.New("minimum window size is too large")
	}
	windows := make([]Window, 0)
	for start := min; start <= length-min-min; start++ {
		for end := start + min; end+min <= length; end++ {
			windows = append(windows, Window{start, end})
		}
	}
	return windows, nil
}

func getBestWindows(metrics map[Metric][]float64, windows []Window, alnLen int, inVarSites []bool) map[Metric]Window {
	// 1) Make an empty array
	// rows = number of metrics
	// columns = number of windows
	// data = nil, allocate new backing slice
	// Each "cell" of the matrix created by {metric}x{window} is the position-wise SSE for that combination
	sses := make(map[Metric]map[Window]float64, 3)

	// 2) Get SSE for each cell in array
	for _, win := range windows {
		// Get SSEs for a given Window
		temp := getSses(metrics, win, inVarSites)
		for m := range temp {
			if _, ok := sses[m][win]; !ok {
				sses[m] = make(map[Window]float64, 0)
				sses[m][win] = sse(temp[m])
			} else {
				sses[m][win] = sse(temp[m])
			}
		}
	}

	// Find minimum values and record the window(s) they occur in
	minMetricWindows := make(map[Metric][]Window)
	for m, windows := range sses {
		bestVal := math.MaxFloat64
		for w, val := range windows {
			if val < bestVal {
				bestVal = val
				minMetricWindows[m] = []Window{w}
			} else if val == bestVal {
				minMetricWindows[m] = append(minMetricWindows[m], w)
			}
		}

	}
	absMinWindow := make(map[Metric]Window)
	for m := range minMetricWindows {
		absMinWindow[m] = getMinVarWindow(minMetricWindows[m], alnLen)
	}

	return absMinWindow
}
