package waitgroup

import(
	"log"
	"sort"
	"sync"
	"strings"
	
	"github.com/RileySun/Babylon"
)

func sortExample() {
	bab := babylon.NewBabylon()
	strRaw1 := bab.Babble(50)
	strRaw2 := bab.Babble(50)
	
	str1 := strings.Split(strRaw1, " ")
	str2 := strings.Split(strRaw2, " ")
	
	var wg sync.WaitGroup
	wg.Add(2)
	
	go sortSlice(str1, &wg)
	go sortSlice(str2, &wg)
	
	wg.Wait()
	
	log.Println("100 strings sorted.")
}

func sortSlice(slice []string, wg *sync.WaitGroup) {
	defer wg.Done()

	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
}