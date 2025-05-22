package semaphore

import(
	"log"
	"sync"
	"time"
	"context"
	"strings"
	
	"golang.org/x/sync/semaphore"
	"github.com/RileySun/Scaled/utils"
)

type Slicer struct {
	sema *semaphore.Weighted
	input chan string
}

func NewSlicer() *Slicer {
	return &Slicer{
		sema:semaphore.NewWeighted(int64(10)),
		input:make(chan string),
	}
}


func (s *Slicer) Start(ctx context.Context, wg *sync.WaitGroup) {
	for str := range s.input {
		if err := s.sema.Acquire(ctx, 1); err != nil {
			//semaphore acquire has failed
			break
		}
		go s.Action(str, wg)
	}
}

func (s *Slicer) Action(str string, wg *sync.WaitGroup) {
	defer s.sema.Release(1)
	time.Sleep(time.Second * time.Duration(utils.GetRandomInt(0, 4)))
	out := strings.Split(str, "/")[1]
	log.Println(out)
	wg.Done()
}