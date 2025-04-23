package main

import(
	"os"
)

//CLI Args = $name, $port
func main() {
	//Start server on port
	server := NewServer(os.Args[1], os.Args[2])
	server.Start()
}