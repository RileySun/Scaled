package pipelines

import(
	"log"
	"sync"
)

func pipelineExample() {
	//Waitgroup is for last stage in pipeline
	var wg sync.WaitGroup
	wg.Add(1)
	
	//Create Pipeline
	downloadStage := download()
	greyscaleStage := greyscale(downloadStage)
	resizeStage := resize(greyscaleStage)
	go save(resizeStage, &wg) //Trigger Pipeline
	
	//Wait
	log.Println("Pipeline Initiated...")
	wg.Wait()
	log.Println("Pipeline Finished")
}