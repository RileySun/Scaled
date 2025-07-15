package insertionsort

import(
	"golang.org/x/exp/constraints"
)

func InsertionSort[T constraints.Ordered](data []T) {
	for i := 1; i < len(data); i++ {
		key := data[i]
		j := i - 1

		for j >= 0 && data[j] > key {
			data[j+1] = data[j]
			j--
		}
		data[j+1] = key
	}
}