package pooling

import(
	"os"
	"log"
	"time"
	"testing"
)

func TestMain(m *testing.M) {	
	exit := m.Run()
	
	os.Exit(exit)
}

func TestPooling(t *testing.T) {
	poolingTime := time.Now()
	poolingExample()
	log.Println("Pooling Example:", time.Since(poolingTime))
}