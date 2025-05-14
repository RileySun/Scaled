package fanin

import(
	"log"
	"sync"
	"context"
)

func collector(ctx context.Context, wg *sync.WaitGroup, results <-chan *Result) {
	defer wg.Done()
	
	var finalResults []*Result
	for {
		select {
			case <-ctx.Done():
				return
			case result, ok := <- results:
				if !ok {
					log.Println("Healthy Results:", len(finalResults))
					return
				}
				
				if result.status == "OK" {
					finalResults = append(finalResults, result)
				}
		}
	}
}