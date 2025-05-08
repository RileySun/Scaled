package unbuffchan

import(
	"os"
	"log"
	"time"
	"context"
	"strconv"
	"net/http"
	"io/ioutil"
	"path/filepath"
	
	"github.com/google/uuid"
)

var count int

func downloadExample(dir string) {
	done := make(chan bool)
		
	for i:=0; i<50; i++ {
		url := "https://dummyjson.com/products/?limit=1&skip="+strconv.Itoa(i)
		go download(url, dir, done)
	}
	
	for i:=0; i<50; i++ {
		<-done
	}
	
	log.Println(count, "files downloaded")
}

func download(url string, dir string, done chan bool) {
	//Timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second * 10))
	defer cancel()
	
	//Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println("Failed to request: " + url)
		done <- true
	}
	
	//Download
	client := &http.Client{}
	res, reqErr := client.Do(req)
	if reqErr != nil {
		log.Println("Failed to download: " + url)
		done <- true
	}
	defer res.Body.Close()
	
	//Read
	byt, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("Failed to read request: " + url)
		done <- true
	}
	
	//Save
	fileName := uuid.New()
	filePath := filepath.Join(dir, fileName.String() + ".json")
	os.WriteFile(filePath, byt, 0755);
	
	count++
	done <- true
}