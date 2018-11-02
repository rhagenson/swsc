package main

import (
	"math"

	"bitbucket.org/rhagenson/swsc/nexus"
)

type window [2]int

func (w *window) Start() int {
	return w[0]
}

func (w *window) Stop() int {
	return w[1]
}

func getAllWindows(uceAln nexus.Alignment, minWin int) []window {
	windows := generateWindows(uceAln.Len(), minWin)
	return windows
}

// generateWindows produces windows of at least a minimum size given a total length
// Windows must be:
//   1) at least minimum window from the start of the UCE (ie, first start at minimum+1)
//   2) at least minimum window from the end of the UCE (ie, last end at length-minimum+1)
//   3) at least minimum window in length (ie, window{start, end)})
func generateWindows(len, min int) []window {
	windows := make([]window, 0)
	// TODO: Should preallocate windows array and fill via start+end-minWin indexing
	for start := min + 1; start+min < len; start++ {
		for end := start + min + 1; end+min <= len; end++ {
			windows = append(windows, window{start, end})
		}
	}
	return windows
}

func getBestWindows(metrics map[Metric][]float64, windows []window, alnLen int, inVarSites []bool) map[Metric]window {
	// 1) Make an empty array
	// rows = number of metrics
	// columns = number of windows
	// data = nil, allocate new backing slice
	// Each "cell" of the matrix created by {metric}x{window} is the position-wise SSE for that combination
	sses := make(map[Metric]map[window]float64, 3)

	// 2) Get SSE for each cell in array
	for _, win := range windows {
		// Get SSEs for a given window
		temp := getSses(metrics, win, inVarSites)
		for m := range temp {
			if _, ok := sses[m][win]; !ok {
				sses[m] = make(map[window]float64, 0)
				sses[m][win] = sse(temp[m])
			} else {
				sses[m][win] = sse(temp[m])
			}
		}
	}

	// Find minimum values and record the window(s) they occur in
	minMetricWindows := make(map[Metric][]window)
	for m, windows := range sses {
		bestVal := math.MaxFloat64
		for w, val := range windows {
			if val < bestVal {
				bestVal = val
				minMetricWindows[m] = []window{w}
			} else if val == bestVal {
				minMetricWindows[m] = append(minMetricWindows[m], w)
			}
		}

	}
	absMinWindow := make(map[Metric]window)
	for m := range minMetricWindows {
		absMinWindow[m] = getMinVarWindow(minMetricWindows[m], alnLen)
	}

	return absMinWindow
}
