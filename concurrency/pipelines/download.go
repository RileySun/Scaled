package pipelines

import(
	"time"
	"bytes"
	"context"
	"strconv"
	"net/http"
	"io/ioutil"
	
	"image"
	"image/jpeg"
)

//Pipeline Result
type Image struct {
	img image.Image
	id string
	err error
}

//Pipeline Function
func download() <-chan Image {
	out := make(chan Image)
	
	go func() {
		defer close(out)
		
		for i:=0; i<50; i++ {
			img, err := downloadImage("https://picsum.photos/200")
			out <- Image{img:img, id:strconv.Itoa(i), err:err}
		}
	}()
	
	return out
}

//Action
func downloadImage(url string) (image.Image, error) {
	//Timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second * 10))
	defer cancel()
	
	//Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return image.NewRGBA(image.Rect(0, 0, 1, 1)), err
	}
	
	//Download
	client := &http.Client{}
	res, reqErr := client.Do(req)
	if reqErr != nil {
		return image.NewRGBA(image.Rect(0, 0, 1, 1)), reqErr
	}
	defer res.Body.Close()
	
	//Read
	byt, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return image.NewRGBA(image.Rect(0, 0, 1, 1)), readErr
	}
	
	//Convert to image.Image
	reader := bytes.NewReader(byt)
	img, convErr := jpeg.Decode(reader)
	if err != nil {
		return image.NewRGBA(image.Rect(0, 0, 1, 1)), convErr
	}
	
	return img, nil
}