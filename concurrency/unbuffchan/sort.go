package unbuffchan

import(
	"log"
	"sort"
	"strings"
	
	"github.com/RileySun/Babylon"
)

func sortExample() {
	bab := babylon.NewBabylon(50)
	strRaw1 := bab.Babble()
	strRaw2 := bab.Babble()
	
	str1 := strings.Split(strRaw1, " ")
	str2 := strings.Split(strRaw2, " ")
	
	done := make(chan bool)
	
	go sortSlice(str1, done)
	go sortSlice(str2, done)
	
	for i := 0; i < 2; i++ {
		<-done
	}
	
	log.Println("100 strings sorted.")
}

func sortSlice(slice []string, done chan bool) {
	sort.Slice(slice, func(i, j int) bool {
		return slice[i] < slice[j]
	})
	
	done <- true
}