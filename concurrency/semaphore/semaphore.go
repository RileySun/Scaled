package semaphore

import(	
	"sync"
	"context"
)

func semaphoreExample()  {
	ctx, cancel := context.WithCancel(context.Background())
	wg := sync.WaitGroup{}
	wg.Add(5)
	
	slicer := NewSlicer()
	go slicer.Start(ctx, &wg)
	
	slicer.input <- "yes/maybe/no"
	slicer.input <- "tomato/potato/beets"
	slicer.input <- "phil/joel/lukas"
	slicer.input <- "ellie/riley/tricia"
	slicer.input <- "house/apartment/shack"
	
	wg.Wait()
	cancel()
}