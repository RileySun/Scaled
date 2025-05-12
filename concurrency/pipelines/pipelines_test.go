package pipelines

import(
	"os"
	"log"
	"time"
	"testing"
)

func TestMain(m *testing.M) {
	//Cleanup Previous
	os.RemoveAll("./images/")
	
	//Create Dirs
	os.Mkdir("./images", 0755)
	
	exit := m.Run()
	
	os.Exit(exit)
}

func TestImage(t *testing.T) {
	pipelineTime := time.Now()
	pipelineExample()
	log.Println("Pipeline Speed:", time.Since(pipelineTime))
}