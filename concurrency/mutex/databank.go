package mutex

import(
	"sync"
	"errors"
)

//Struct
type Databank struct {
	Data map[string]int
	mu sync.RWMutex
}

//Creae 
func NewDatabank() *Databank {
	return &Databank{
		Data:make(map[string]int),
	}
}

//Methods
func (d *Databank) Has(key string) bool {
	d.mu.RLock()
	_, ok := d.Data[key]
	d.mu.RUnlock()
	return ok
}

func (d *Databank) Get(key string) (int, error) {
	if d.Has(key) {
		d.mu.RLock()
		defer d.mu.RUnlock()
		return d.Data[key], nil
	} else {
		return 0, errors.New("Invalid Key")
	}
}

func (d *Databank) Add(key string, value int) error {	
	if !d.Has(key) {
		d.mu.Lock()
		d.Data[key] = value
		d.mu.Unlock()
		return nil
	} else {
		return errors.New("Key Already Exists")
	}
}

func (d *Databank) Set(key string, value int) error {
	if d.Has(key) {
		d.mu.Lock()
		d.Data[key] = value
		d.mu.Unlock()
		return nil
	} else {
		return errors.New("Invalid Key")
	}
}

func (d *Databank) Remove(key string) error {
	if !d.Has(key) {
		d.mu.Lock()
		delete(d.Data, key)
		d.mu.Unlock()
		return nil
	} else {
		return errors.New("Invalid Key")
	}
}