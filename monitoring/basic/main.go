package main

import(
	"os"
)

//CLI Args = $name, $port
func main() {
	server := NewServer(os.Args[1], os.Args[2])
	server.Start()
}