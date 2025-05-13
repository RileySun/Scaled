package fanout

import(
	"os"
	"os/signal"
	
	"sync"
	"context"
)

func fanoutExample(urls ...string) {
	//Ctx & Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	graceful(cancel)
	
	//Channels
	urlChan := make(chan string)
	
	//Wg
	var wg sync.WaitGroup
	
	//Workers
	for i := 0; i < len(urls); i++ {
		wg.Add(1)
		go worker(ctx, &wg, urlChan)
	}
	
	//Initiate Workers
	go func() {
		defer close(urlChan)
		
		for _, url := range urls {
			select {
				case <-ctx.Done():
					return
				case urlChan <- url:
			}
		}
	}()
	
	//Wait for results
	wg.Wait()
}

func graceful(cancel func()) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	
	go func() {
		<-stop
		cancel()
	}()
}