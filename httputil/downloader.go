package httputil

import (
	"context"
	"io/ioutil"
	"net/http"
)

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)
	go func() { c <- f(client.Do(req)) }()
	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}

// MakeRequest do http request and exec F function on it Response
func MakeRequest(ctx context.Context, method, url string, f func(*http.Response, error) error) error {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return err
	}
	return httpDo(ctx, req, f)
}

// MakeRequestCC Concurent
func MakeRequestCC(ctx context.Context, method, url string, f func(*http.Response, error) error, done chan error) {
	go func() { done <- MakeRequest(ctx, method, url, f) }()
}

// DoRequest executes the request and return it response body content
func DoRequest(ctx context.Context, method, url string) ([]byte, error) {
	var bytes []byte
	var save = func(r *http.Response, err error) error {
		if r.Body == nil {
			return nil
		}
		defer r.Body.Close()
		bytes, err = ioutil.ReadAll(r.Body)
		return err
	}
	err := MakeRequest(ctx, method, url, save)
	return bytes, err
}

// GetBody no context
func GetBody(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// GetBodyConcurrent no context
func GetBodyCC(url string, b chan []byte) {
	go func() {
		bytes, _ := GetBody(url)
		b <- bytes
	}()
}
