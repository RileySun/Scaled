package limiter

import(
	"log"
	"time"
	"context"
	"net/http"
	"io/ioutil"
)

type Limiter interface {
	Action(context.Context)
	SetLimit(time.Duration)
	Start()
	Stop()
}

type APICaller struct {
	url string //address
	limit *time.Ticker //rate
	Output chan []byte //download output
}

func NewAPICaller(url string) *APICaller {
	return &APICaller{
		url:url,
	}
}

func (a *APICaller) SetLimit(d time.Duration) {
	a.limit = time.NewTicker(d)
}

func (a *APICaller) Start(ctx context.Context) {
	a.Output = make(chan []byte)
	
	for range a.limit.C {
		select {
			case <-ctx.Done():
				a.Stop()
				return
			default:
				_, err := a.Action(ctx) //black holing output
				if err != nil {
					log.Println("API call failure")
					continue
				}
				log.Println("API call success")
				//a.Output <- byt
		}
	}
}

func (a *APICaller) Stop() {
	a.limit.Stop()
	close(a.Output)
}

func (a *APICaller) Action(mainCtx context.Context) ([]byte, error) {
	//Timeout
	ctx, cancel := context.WithTimeout(mainCtx, time.Duration(time.Second * 10))
	defer cancel()
	
	//Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, a.url, nil)
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