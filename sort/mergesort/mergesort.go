package mergesort

import(
	"golang.org/x/exp/constraints"
)

func MergeSort[T constraints.Ordered](left, right []T) []T {
	var result []T
	leftIndex, rightIndex := 0, 0

	for leftIndex < len(left) && rightIndex < len(right) {
		if left[leftIndex] < right[rightIndex] {
			result = append(result, left[leftIndex])
			leftIndex++
		} else {
			result = append(result, right[rightIndex])
			rightIndex++
		}
	}

	// Append any remaining elements
	result = append(result, left[leftIndex:]...)
	result = append(result, right[rightIndex:]...)

	return result
}