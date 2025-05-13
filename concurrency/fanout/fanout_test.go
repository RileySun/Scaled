package fanout

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
	fanoutTime := time.Now()
	urls := []string{"http://yes.com", "http://google.com", "http://potatotomatopotato.com", "http://notarealurl.null"}
	fanoutExample(urls...)
	log.Println("Urls checked in:", time.Since(fanoutTime))
}