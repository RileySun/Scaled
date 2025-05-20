package limiter

import(
	"time"
	"context"
)

func limiterExample()  {
	api := NewAPICaller("https://picsum.photos/200")
	api.SetLimit(time.Second * 1)
	
	ctx, cancel := context.WithCancel(context.Background())
	
	go api.Start(ctx)
	
	time.Sleep(time.Second * 10)
	cancel()
}