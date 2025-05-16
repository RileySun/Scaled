package errgroup

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

func TestErrgroup(t *testing.T) {
	groupTime := time.Now()
	errGroupExample()
	log.Println("User data retrieved:", time.Since(groupTime))
}