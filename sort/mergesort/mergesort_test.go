package mergesort

import(
	"log"
	"time"
	"testing"
	
	"github.com/RileySun/Babylon"
)

func TestInt(t *testing.T) {
	testTime := time.Now()
	log.Println("Start Int Sort")
	slice1 := []int{1, 54, 63, 26, 356, 74, 223}
	slice2 := []int{63, 7, 235, 65, 57, 29}
	log.Println(slice1, " | ", slice2)
	slice := MergeSort(slice1, slice2)
	log.Println(slice)
	log.Println("End Int Sort: ", time.Since(testTime))
}

func TestFloat(t *testing.T) {
	testTime := time.Now()
	log.Println("Start Float Sort")
	slice1 := []float64{0.2323, 0.454, 0.7345, 0.46763, 0.45575}
	slice2 := []float64{0.5757, 0.46576, 0.243245, 0.35575}
	log.Println(slice1, " | ", slice2)
	slice := MergeSort(slice1, slice2)
	log.Println(slice)
	log.Println("End Float Sort: ", time.Since(testTime))
}

func TestString(t *testing.T) {
	testTime := time.Now()
	log.Println("Start String Sort")
	slice1 := babylon.NewBabylon().BabbleSlice(6)
	slice2 := babylon.NewBabylon().BabbleSlice(6)
	log.Println(slice1, " | ", slice2)
	slice := MergeSort(slice1, slice2)
	log.Println(slice)
	log.Println("End String Sort: ", time.Since(testTime))
}