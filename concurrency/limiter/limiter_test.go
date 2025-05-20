package limiter

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

func TestLimiter(t *testing.T) {
	limiterTime := time.Now()
	limiterExample()
	log.Println("Rate Limiter Example:", time.Since(limiterTime))
}