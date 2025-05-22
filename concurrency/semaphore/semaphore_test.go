package semaphore

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

func TestSemaphore(t *testing.T) {
	semaphoreTime := time.Now()
	semaphoreExample()
	log.Println("Semaphore Example:", time.Since(semaphoreTime))
}