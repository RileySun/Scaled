package pipelines

import(
	"os"
	"log"
	"sync"
	"bytes"
	"image"
	"image/jpeg"
	"path/filepath"
	
	"github.com/google/uuid"
)

//Pipeline Function
func save(in <-chan ResizedImage, wg *sync.WaitGroup) {
	defer wg.Done()
	
	for resized := range in {
		//Error Checking
		if resized.err != nil {
			log.Println("Error Resizing Image", resized.id)
			continue
		}
			
		err := saveImage(resized.img)
		if err != nil {
			log.Println(err)
		}
	}
}

//Action
func saveImage(img image.Image) error {
	//Create File
	buf := new(bytes.Buffer)
	encodeErr := jpeg.Encode(buf, img, nil)
	if encodeErr != nil {
		return encodeErr
	}
	fileBytes := buf.Bytes()
	
	//Save File
	fileName := uuid.New()
	filePath := filepath.Join("./images", fileName.String() + ".jpeg")
	os.WriteFile(filePath, fileBytes, 0755);
	
	//Ok
	return nil
}