package caching

import(
	"time"
	"testing"
	
	"github.com/RileySun/Babylon"
)

//t.Error(err)
//t.Fail()

var cacher = NewCacher()

func TestSet(t *testing.T) {
	//Random Key/Value Pairs (words)
	babbler := babylon.NewBabylon()
	
	//Set
	var randomKey string
	for i:=0; i<1000; i++ {
		key := babbler.Babble()
		err := cacher.Set(key, babbler.Babble(1), time.Minute * 10)
		if err != nil {
			t.Error(err)
		}
		
		if i == 0 {
			randomKey = key
		}
	}
	
	//Get
	_, err := cacher.Get(randomKey)
	if err != nil {
		t.Error(err)
	}
	
	//Remove
	err = cacher.Remove(randomKey)
	if err != nil {
		t.Error(err)
	}
	
	//Clear
	err = cacher.Clear()
	if err != nil {
		t.Error(err)
	}
}