package main

import(
	"context"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	
	server := NewVideoServer()
	server.Start(ctx)
	
	//server.Stop(cancel)
}