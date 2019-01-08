package botutil

import (
	"github.com/chippydip/go-sc2ai/api"
)

func sortUnits(data *[]*api.Unit) {
	//sort.Sort((*sorter)(data))
	quickSort(*data, maxDepth(len(*data)))

	// check sort
	// prev := data[0].UnitType
	// for _, uu := range data {
	// 	if uu.UnitType < prev {
	// 		panic("not sorted!")
	// 	}
	// 	prev = uu.UnitType
	// }
}

// Generic sort.Interface wrapper

type sorter []*api.Unit

func (s *sorter) Len() int           { return len(*s) }
func (s *sorter) Swap(i, j int)      { (*s)[i], (*s)[j] = (*s)[j], (*s)[i] }
func (s *sorter) Less(i, j int) bool { return (*s)[i].UnitType < (*s)[j].UnitType }

// Type-specialized version of sort.Sort

func maxDepth(n int) int {
	var depth int
	for i := n; i > 0; i >>= 1 {
		depth++
	}
	return depth * 2
}

func quickSort(data []*api.Unit, maxDepth int) {
	for len(data) > 12 { // Use ShellSort for slices <= 12 elements
		if maxDepth == 0 {
			heapSort(data)
			return
		}
		maxDepth--
		mlo, mhi := doPivot(data)
		// Avoiding recursion on the larger subproblem guarantees
		// a stack depth of at most lg(b-a).
		if mlo < len(data)-mhi {
			quickSort(data[:mlo], maxDepth)
			//a = mhi // i.e., quickSort(data, mhi, b)
			data = data[mhi:]
		} else {
			quickSort(data[mhi:], maxDepth)
			//b = mlo // i.e., quickSort(data, a, mlo)
			data = data[:mlo]
		}
	}
	if len(data) > 1 {
		// Do ShellSort pass with gap 6
		// It could be written in this simplified form cause b-a <= 12
		for i := 6; i < len(data); i++ {
			if data[i].UnitType < data[i-6].UnitType {
				data[i], data[i-6] = data[i-6], data[i]
			}
		}
		insertionSort(data)
	}
}

func heapSort(data []*api.Unit) {
	first := 0
	lo := 0
	hi := len(data)

	// Build heap with greatest element at top.
	for i := (hi - 1) / 2; i >= 0; i-- {
		siftDown(data, i, hi, first)
	}

	// Pop elements, largest first, into end of data.
	for i := hi - 1; i >= 0; i-- {
		data[first], data[first+i] = data[first+i], data[first]
		siftDown(data, lo, i, first)
	}
}

func siftDown(data []*api.Unit, lo, hi, first int) {
	root := lo
	for {
		child := 2*root + 1
		if child >= hi {
			break
		}
		if child+1 < hi && data[first+child].UnitType < data[first+child+1].UnitType {
			child++
		}
		if data[first+root].UnitType >= data[first+child].UnitType {
			return
		}
		data[first+root], data[first+child] = data[first+child], data[first+root]
		root = child
	}
}

func insertionSort(data []*api.Unit) {
	for i := 0 + 1; i < len(data); i++ {
		for j := i; j > 0 && data[j].UnitType < data[j-1].UnitType; j-- {
			data[j], data[j-1] = data[j-1], data[j]
		}
	}
}

func doPivot(data []*api.Unit) (midlo, midhi int) {
	m := len(data) / 2
	if len(data) > 40 {
		// Tukey's ``Ninther,'' median of three medians of three.
		s := len(data) / 8
		medianOfThree(data, 0, 0+s, 0+2*s)
		medianOfThree(data, m, m-s, m+s)
		medianOfThree(data, len(data)-1, len(data)-1-s, len(data)-1-2*s)
	}
	medianOfThree(data, 0, m, len(data)-1)

	// Invariants are:
	//	data[lo] = pivot (set up by ChoosePivot)
	//	data[lo < i < a] < pivot
	//	data[a <= i < b] <= pivot
	//	data[b <= i < c] unexamined
	//	data[c <= i < hi-1] > pivot
	//	data[hi-1] >= pivot
	pivot := 0
	a, c := 0+1, len(data)-1

	for ; a < c && data[a].UnitType < data[pivot].UnitType; a++ {
	}
	b := a
	for {
		for ; b < c && data[pivot].UnitType >= data[b].UnitType; b++ { // data[b] <= pivot
		}
		for ; b < c && data[pivot].UnitType < data[c-1].UnitType; c-- { // data[c-1] > pivot
		}
		if b >= c {
			break
		}
		// data[b] > pivot; data[c-1] <= pivot
		data[b], data[c-1] = data[c-1], data[b]
		b++
		c--
	}
	// If hi-c<3 then there are duplicates (by property of median of nine).
	// Let be a bit more conservative, and set border to 5.
	protect := len(data)-c < 5
	if !protect && len(data)-c < len(data)/4 {
		// Lets test some points for equality to pivot
		dups := 0
		if data[pivot].UnitType >= data[len(data)-1].UnitType { // data[hi-1] = pivot
			data[c], data[len(data)-1] = data[len(data)-1], data[c]
			c++
			dups++
		}
		if data[b-1].UnitType >= data[pivot].UnitType { // data[b-1] = pivot
			b--
			dups++
		}
		// m-lo = (hi-lo)/2 > 6
		// b-lo > (hi-lo)*3/4-1 > 8
		// ==> m < b ==> data[m] <= pivot
		if data[m].UnitType >= data[pivot].UnitType { // data[m] = pivot
			data[m], data[b-1] = data[b-1], data[m]
			b--
			dups++
		}
		// if at least 2 points are equal to pivot, assume skewed distribution
		protect = dups > 1
	}
	if protect {
		// Protect against a lot of duplicates
		// Add invariant:
		//	data[a <= i < b] unexamined
		//	data[b <= i < c] = pivot
		for {
			for ; a < b && data[b-1].UnitType >= data[pivot].UnitType; b-- { // data[b] == pivot
			}
			for ; a < b && data[a].UnitType < data[pivot].UnitType; a++ { // data[a] < pivot
			}
			if a >= b {
				break
			}
			// data[a] == pivot; data[b-1] < pivot
			data[a], data[b-1] = data[b-1], data[a]
			a++
			b--
		}
	}
	// Swap pivot into middle
	data[pivot], data[b-1] = data[b-1], data[pivot]
	return b - 1, c
}

func medianOfThree(data []*api.Unit, m1, m0, m2 int) {
	// sort 3 elements
	if data[m1].UnitType < data[m0].UnitType {
		data[m1], data[m0] = data[m0], data[m1]
	}
	// data[m0] <= data[m1]
	if data[m2].UnitType < data[m1].UnitType {
		data[m2], data[m1] = data[m1], data[m2]
		// data[m0] <= data[m2] && data[m1] < data[m2]
		if data[m1].UnitType < data[m0].UnitType {
			data[m1], data[m0] = data[m0], data[m1]
		}
	}
	// now data[m0] <= data[m1] <= data[m2]
}
