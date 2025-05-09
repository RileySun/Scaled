package waitgroup

import(
	"os"
	"log"
	"sync"
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
	var wg sync.WaitGroup
	wg.Add(50)
		
	for i:=0; i<50; i++ {
		url := "https://dummyjson.com/products/?limit=1&skip="+strconv.Itoa(i)
		go download(url, dir, &wg)
	}
	
	wg.Wait()
	
	log.Println(count, "files downloaded")
}

func download(url string, dir string, wg *sync.WaitGroup) {
	//Waitgroup
	defer wg.Done()
	
	//Timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second * 10))
	defer cancel()
	
	//Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		log.Println("Failed to request: " + url)
		wg.Done()
	}
	
	//Download
	client := &http.Client{}
	res, reqErr := client.Do(req)
	if reqErr != nil {
		log.Println("Failed to download: " + url)
		wg.Done()
	}
	defer res.Body.Close()
	
	//Read
	byt, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		log.Println("Failed to read request: " + url)
		wg.Done()
	}
	
	//Save
	fileName := uuid.New()
	filePath := filepath.Join(dir, fileName.String() + ".json")
	os.WriteFile(filePath, byt, 0755);
	
	count++
}