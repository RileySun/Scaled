package pooling

import(
	"io"
	"time"
	"errors"
	"context"
	"net/http"
	"io/ioutil"
)


type Client struct {
	httpClient *http.Client
	Headers map[string]string
}

func NewClient() *Client {
	return &Client{
		httpClient:&http.Client{Timeout:time.Second * 5},
	}
}

//Methods
func (c *Client) Get(mainCtx context.Context, url string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(mainCtx, time.Duration(time.Second * 5))
	defer cancel()
	
	//Request
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return []byte{}, err
	}
	
	//Set Headers
	for _, key := range c.Headers {
		req.Header.Set(key, c.Headers[key])
	}
	
	//Download
	res, reqErr := c.httpClient.Do(req)
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

func (c *Client) Post(mainCtx context.Context, url string, body io.Reader) ([]byte, error) {
	ctx, cancel := context.WithTimeout(mainCtx, time.Duration(time.Second * 5))
	defer cancel()
	
	//Request
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, body)
	if err != nil {
		return []byte{}, err
	}
	
	//Set Headers
	for _, key := range c.Headers {
		req.Header.Set(key, c.Headers[key])
	}
	
	//Download
	res, reqErr := c.httpClient.Do(req)
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

//Headers
func (c *Client) AddHeader(key string, value string) error {
	_, ok := c.Headers[key]
	if ok {
		return errors.New("Header already exists, please remove before adding")
	}
	
	c.Headers[key] = value
	return nil
}

func (c *Client) RemoveHeader(key string) error {
	_, ok := c.Headers[key]
	if !ok {
		return errors.New("Header does not exist, please add before removing")
	}
	
	delete(c.Headers, key)
	return nil
}