package fanin

import(
	"time"
	"context"
	"net/http"
)

type Result struct {
	url, status string
	err error
}

func status(parentCtx context.Context, url string) *Result {
	ctx, cancel := context.WithTimeout(parentCtx, time.Duration(time.Second * 10))
	defer cancel()
	
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return &Result{url:url, status:"ERROR", err:err}
	}
	
	client := &http.Client{}
	res, reqErr := client.Do(req)
	if reqErr != nil {
		return &Result{url:url, status:"ERROR", err:reqErr}
	}
	defer res.Body.Close()
	
	return &Result{url:url, status:res.Status, err:nil}
}