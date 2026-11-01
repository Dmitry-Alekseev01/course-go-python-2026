package urls

import (
	"context"
	"log/slog"
	"net/http"
	"sync"
)

type ResponseCodeAggregator struct {
	client HttpClient
}

func NewURLAggregator(client HttpClient) *ResponseCodeAggregator {
	return &ResponseCodeAggregator{client: client}
}

func (u *ResponseCodeAggregator) Aggregate(ctx context.Context, urls <-chan string) (AggregationResult, error) {
	result := newResult()
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			select {
			case url, ok := <-urls:
				if !ok {
					return
				}
				wg.Add(1)
				go func(url string) {
					defer wg.Done()

					type response struct {
						resp *http.Response
						err  error
					}
					ch := make(chan response, 1)

					go func() {
						resp, err := u.client.Get(url)
						ch <- response{resp: resp, err: err}
					}()

					select {
					case r := <-ch:
						if r.err != nil {
							slog.Error(r.err.Error())
							return
						}
						statusCode := r.resp.StatusCode
						r.resp.Body.Close()
						result.addResponse(statusCode)
					case <-ctx.Done():
						return
					}
				}(url)
			case <-ctx.Done():
				return
			}
		}
	}()

	go func() {
		wg.Wait()
		result.setDone()
	}()

	return result, nil
}
