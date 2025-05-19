package mutex

import(
	"sync"
	"time"
	"errors"
)

func mutexExample() error {
	data := NewDatabank()
	
	//Add to databank
	err := data.Add("Potato", 3)
	if err != nil {
		return err
	}
	
	//Check it was added
	if !data.Has("Potato") {
		return errors.New("Potato key does not exist")
	}
	
	//Run two operations at once (would cause issues if not mutex locked)
	wg := sync.WaitGroup{}
	wg.Add(2)
	
	go func(wg *sync.WaitGroup, err error) {
		time.Sleep(time.Second * 2)
		err = data.Set("Potato", 0)
		wg.Done()
	}(&wg, err)
	
	go func(wg *sync.WaitGroup, err error) {
		time.Sleep(time.Second * 2)
		err = data.Remove("Potato") 
		wg.Done()
	}(&wg, err)
	
	//Wait for goroutines & return nil
	wg.Wait()
	return err
}