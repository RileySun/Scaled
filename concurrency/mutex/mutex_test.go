package mutex

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

func TestFanout(t *testing.T) {
	mutexTime := time.Now()
	err := mutexExample()
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	log.Println("Mutex example:", time.Since(mutexTime))
}