package pooling

import(	
	"log"
	"sync"
	"context"
)

func poolingExample()  {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	var wg sync.WaitGroup
	
	pool := NewPool(3)
	for i:=0; i<10; i++ {
		wg.Add(1)
		
		go func() {
			defer wg.Done()
			
			client := pool.Acquire()
			defer pool.Release(client)
			
			_, err := client.Get(ctx, "https://picsum.photos/200")
			if err != nil {
				log.Println(i, err)
			}
		}()
	}
	
	wg.Wait()
}