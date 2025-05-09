package waitgroup

import(
	"os"
	"log"
	"sync"
	"time"
	"bytes"
	"image"
	"image/jpeg"
	"context"
	"net/http"
	"io/ioutil"
	"path/filepath"
	
	"github.com/google/uuid"
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
)

func imageExample(dir string) {	
	var wg sync.WaitGroup
	wg.Add(50)
	
	for i:=0; i<50; i++ {
		go greyscale("https://picsum.photos/200", dir, &wg)
	}
	
	wg.Wait()
	
	log.Println(count, "images edited")
}

func greyscale(url string, dir string, wg *sync.WaitGroup) {
	defer wg.Done()
	
	//Get Image
	imageBytes, err := downloadImage(url)
	if err != nil {
		log.Println("Could not download image", url)
		wg.Done()
	}
	
	//Image from bytes
	img := imageFromBytes(imageBytes)
	
	//GreyScale
	img = effect.Grayscale(img)
	
	//Erode
	img = effect.Erode(img, 1)
	
	//Resize
	img = transform.Resize(img, 200, 200, transform.Linear) //Just for show (they small enough already)
	
	//Create File
	buf := new(bytes.Buffer)
	encodeErr := jpeg.Encode(buf, img, nil)
	if encodeErr != nil {
		log.Println("Could not create image file:", url)
		wg.Done()
	}
	fileBytes := buf.Bytes()
	
	//Save File
	fileName := uuid.New()
	filePath := filepath.Join(dir, fileName.String() + ".jpeg")
	os.WriteFile(filePath, fileBytes, 0755);
}

func downloadImage(url string) ([]byte, error) {
	//Timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second * 10))
	defer cancel()
	
	//Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, err
	}
	
	//Download
	client := &http.Client{}
	res, reqErr := client.Do(req)
	if reqErr != nil {
		return []byte{}, reqErr
	}
	defer res.Body.Close()
	
	//Read
	byt, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return []byte{}, readErr
	}
	
	return byt, nil
}

func imageFromBytes(byt []byte) image.Image {
	r := bytes.NewReader(byt)
	i, err := jpeg.Decode(r)
	if err != nil {
		log.Fatal("Byt2Img - " + err.Error())
	}
	return i
}