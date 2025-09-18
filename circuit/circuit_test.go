package circuit

import(
	"testing"
	
	"time"
	"errors"
	"context"
	"net/http"
	"io/ioutil"
)

//Intentionally fail download from non-existant URL
func TestFailingURL(t *testing.T) {
	circuit := NewCircuit(5)
	url := "https://nowaythisexists.com/"
	
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(time.Second * 10))
	defer cancel()
	
	action := func()(any, error){return download(ctx, url)}
	
	for i:=0; i<7; i++ {
		//We know there will be no result and err will always return an error
		_, err := circuit.Execute(ctx, action)
		
		if err == nil {
			t.Error(errors.New("Request should have returned and error"))
			t.Fail()
		}
		
		if circuit.State() == "Broken" {
			break
		}
	}
	
	
	//Final Assert
	if circuit.State() != "Broken" {
		t.Error(errors.New("Circuit breaker did not trip"))
		t.Fail()
	}
}


func download(ctx context.Context, url string) ([]byte, error) {
	//Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
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