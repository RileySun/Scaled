package fanin

import(
	"sync"
	"time"
	
	"github.com/RileySun/Scaled/utils"
)

type Result struct {
	value int
	status string
}

func mock() <-chan *Result {
	out := make(chan *Result)
	
	var wg sync.WaitGroup

	go func() {
		for i:=0; i<25; i++ {
			wg.Add(1)
			
			go func() {
				defer wg.Done()
				sec := utils.GetRandomInt(0, 4)
				time.Sleep(time.Second * time.Duration(sec))
			
				stat := "OK"
				if sec > 2 {
					stat = "ERROR"
				}
					
				out <- &Result{value:sec, status:stat}
			}()
		}
		
		//Wait to close
		wg.Wait()
		close(out)
	}()
	
	return out
}