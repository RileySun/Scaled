package pipelines

import(
	"log"
	"image"
	
	"github.com/anthonynsimon/bild/effect"
	"github.com/anthonynsimon/bild/transform"
)

//Pipeline Result
type ResizedImage struct {
	img image.Image
	id string
	err error
}

//Pipeline Function
func resize(in <-chan GreyscaleImage) <-chan ResizedImage {
	out := make(chan ResizedImage)
	
	go func() {
		defer close(out)
		
		for greyscale := range in {
			//Error Checking
			if greyscale.err != nil {
				log.Println("Error Greyscaling Image", greyscale.id)
				continue
			}
			
			resized := resizeImage(greyscale.img)
			out <- ResizedImage{img:resized, id:greyscale.id, err:nil}
		}
	}()
	
	return out
}

//Action
func resizeImage(img image.Image) image.Image {
	//Erode
	img = effect.Erode(img, 1)
	
	//Resize
	img = transform.Resize(img, 200, 200, transform.Linear) //Just for show (they small enough already)
	
	return img
}