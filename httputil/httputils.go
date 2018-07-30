package stormtf

import (
	"context"
	"io"
	"net/http"
)

func httpDo(ctx context.Context, req *http.Request, f func(*http.Response, error) error) error {
	// Run the HTTP request in a goroutine and pass the response to f.
	tr := &http.Transport{}
	client := &http.Client{Transport: tr}
	c := make(chan error, 1)

	go func() {
		c <- f(client.Do(req))
	}()

	select {
	case <-ctx.Done():
		tr.CancelRequest(req)
		<-c // Wait for f to return.
		return ctx.Err()
	case err := <-c:
		return err
	}
}

func downloadBodyRC(ctx context.Context, url string) (io.ReadCloser, error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	var body io.ReadCloser
	err = httpDo(ctx, req, func(r *http.Response, e error) error {
		if e != nil {
			return e
		}
		body = r.Body
		return nil
	})
	if err != nil {
		return nil, err
	}
	return body, nil
}
