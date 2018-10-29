package main

import (
	"math"

	"github.com/biogo/biogo/seq/multi"
	"gonum.org/v1/gonum/stat"
)

func minFloat64(vs ...float64) float64 {
	min := math.MaxFloat64
	for _, v := range vs {
		if v < min {
			min = v
		}
	}
	return min
}

// TODO: Should use math.Big to prevent overflow
func factorial(v int) float64 {
	fact := 1
	for i := 1; i <= v; i++ {
		fact *= i
	}
	return float64(fact)
}

func minInCountsMap(counts map[byte]int) int {
	min := math.MaxInt16
	for _, val := range counts {
		if val < min {
			min = val
		}
	}
	return min
}

func maxInFreqMap(freqs map[byte]float32) float32 {
	max := float32(math.SmallestNonzeroFloat32)
	for _, val := range freqs {
		if max < val {
			max = val
		}
	}
	return max
}

func getMinVarWindow(windows []window, alnLength int) window {
	best := float64(math.MaxInt16)
	bestWindow := windows[0]

	for _, w := range windows {
		l1 := float64(w[0])
		l2 := float64(w[1] - w[0])
		l3 := float64(alnLength - w[0])
		variance := stat.Variance([]float64{l1, l2, l3}, nil)
		if variance < best {
			best = variance
			bestWindow = w
		}
	}
	return bestWindow
}

// anyUndeterminedBlocks checks if any blocks are only undetermined/ambiguous characters
// Not the same as anyBlocksWoAllSites()
func anyUndeterminedBlocks(bestWindow window, uceAln *multi.Multi) bool {
	leftAln, _ := uceAln.Subseq(0, bestWindow[0])
	coreAln, _ := uceAln.Subseq(bestWindow[0], bestWindow[1])
	rightAln, _ := uceAln.Subseq(bestWindow[1], uceAln.Len())

	leftFreq := bpFreqCalc(leftAln)
	coreFreq := bpFreqCalc(coreAln)
	rightFreq := bpFreqCalc(rightAln)

	// If any frequency is NaN
	// TODO: Likely better with bpFreqCalc returning an error value
	if maxInFreqMap(leftFreq) == 0 || maxInFreqMap(coreFreq) == 0 || maxInFreqMap(rightFreq) == 0 {
		return true
	}
	return false
}

// anyBlocksWoAllSites checks for blocks with only undetermined/ambiguous characters
// Not the same as anyUndeterminedBlocks()
func anyBlocksWoAllSites(bestWindow window, uceAln *multi.Multi) bool {
	leftAln, _ := uceAln.Subseq(0, bestWindow[0])
	coreAln, _ := uceAln.Subseq(bestWindow[0], bestWindow[1])
	rightAln, _ := uceAln.Subseq(bestWindow[1], uceAln.Len())

	leftCounts := countBases(leftAln)
	coreCounts := countBases(coreAln)
	rightCounts := countBases(rightAln)

	if minInCountsMap(leftCounts) == 0 || minInCountsMap(coreCounts) == 0 || minInCountsMap(rightCounts) == 0 {
		return true
	}
	return false
}

func csvColToPlotMatrix(best window, n int) []int8 {
	matrix := make([]int8, n)
	for i := range matrix {
		switch {
		case i < best[0]:
			matrix[i] = -1
		case best[0] < i && i < best[1]:
			matrix[i] = 0
		case best[1] < i:
			matrix[i] = 1
		}
	}
	return matrix
}