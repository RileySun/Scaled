package quicksort

import(
	"time"
	"math/rand"
	"golang.org/x/exp/constraints"
)

/* Set Random Seed */
func init() {
	rand.Seed(time.Now().UnixNano())
}

func QuickSort[T constraints.Ordered](data []T) {
	if len(data) < 2 {
		return
	}

	left, right := 0, len(data)-1
	pivotIndex := rand.Int() % len(data)
	data[pivotIndex], data[right] = data[right], data[pivotIndex]
	for i := range data {
		if data[i] < data[right] {
			data[i], data[left] = data[left], data[i]
			left++
		}
	}
	data[left], data[right] = data[right], data[left]
	QuickSort(data[:left])
	QuickSort(data[left+1:])
}