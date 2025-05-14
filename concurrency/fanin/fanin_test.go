package fanin

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

func TestFanin(t *testing.T) {
	faninTime := time.Now()
	faninExample()
	log.Println("Statuses checked in:", time.Since(faninTime))
}