package fanin

import(
	"os"
	"os/signal"
	
	"sync"
	"context"
)

func faninExample() {
	//Ctx & Graceful Shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	graceful(cancel)
	
	//Wg
	var wg sync.WaitGroup
	wg.Add(1)
	
	mockData := mock()
	go collector(ctx, &wg, mockData)
	
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