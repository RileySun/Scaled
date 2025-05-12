package pipelines

import(
	"log"
	"image"
	
	"github.com/anthonynsimon/bild/effect"
)

//Pipeline Result
type GreyscaleImage struct {
	img image.Image
	id string
	err error
}

//Pipeline Function
func greyscale(in <-chan Image) <-chan GreyscaleImage {
	out := make(chan GreyscaleImage)
	
	go func() {
		defer close(out)
		
		for download := range in {
			//Error Checking
			if download.err != nil {
				log.Println("Error Downloading Image", download.id)
				continue
			}
		
			resized := greyscaleImage(download.img)
			out <- GreyscaleImage{img:resized, id:download.id, err:nil}
		}
	}()
	
	return out
}

//Action
func greyscaleImage(img image.Image) image.Image {
	return effect.Grayscale(img)
}