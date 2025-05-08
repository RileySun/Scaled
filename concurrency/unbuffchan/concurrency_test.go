package unbuffchan

import(
	"os"
	"log"
	"time"
	"testing"
)

func TestMain(m *testing.M) {
	//Cleanup Previous
	os.RemoveAll("./images/")
	os.RemoveAll("./products/")
	
	//Create Dirs
	os.Mkdir("./products", 0755)
	os.Mkdir("./images", 0755)
	
	exit := m.Run()
	
	os.Exit(exit)
}

func TestSort(t *testing.T) {
	sortTime := time.Now()
	sortExample()
	log.Println("Sorted in:", time.Since(sortTime))
}

func TestDownload(t *testing.T) {
	downTime := time.Now()
	downloadExample("./products")
	log.Println("Downloaded in:", time.Since(downTime))
}

func TestImage(t *testing.T) {
	imageTime := time.Now()
	imageExample("./images")
	log.Println("Images edited in:", time.Since(imageTime))
}