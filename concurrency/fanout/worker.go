package fanout

import(
	"log"
	"sync"
	"time"
	"context"
)

func worker(ctx context.Context, wg *sync.WaitGroup, urls <-chan string) {
	defer wg.Done()
	
	for {
		select {
			case <-ctx.Done(): //Parent ctx done
				return
			case url, ok := <-urls:
				if !ok {
					return
				} //Closed Channel
				
				workerCtx, cancel := context.WithTimeout(ctx, time.Second * 2)
				defer cancel()
				
				result := status(workerCtx, url)
				
				str := "OK"
				if result.err != nil {
					str = "Error"
				}
				log.Println(str + ":", result.url, " - ", result.status)
		}
	}
}