package main

import (
	"context"
	"log"
	"net/http"
	"runtime"
	"time"
)

type result struct {
	url     string
	err     error
	latency time.Duration
}

func get(ctx context.Context, url string, ch chan<- result) {
	start := time.Now()
	var r result
	ticker := time.NewTicker(1 * time.Second).C
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if resp, err := http.DefaultClient.Do(req); err != nil {
		r = result{url, err, 0}
	} else {
		t := time.Since(start).Round(time.Millisecond)
		r = result{url, err, t}
		resp.Body.Close()
	}
	for {
		select {
		case ch <- r:
			return
		case <-ticker:
			log.Println("tick", r)
		}
	}

}
func first(ctx context.Context, urls []string) (*result, error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	results := make(chan result, len(urls))

	for _, url := range urls {
		go get(ctx, url, results)
	}
	select {
	case r := <-results:
		return &r, nil
	case <-ctx.Done():
		return nil, ctx.Err()

	}
}
func main() {
	//results := make(chan result)
	list := []string{
		"https://google.com",
		"https://amazon.in",
		"https://wsj.com",
		"https://bing.com",
	}
	r, _ := first(context.Background(), list)
	if r.err != nil {
		log.Printf("%-20s %s", r.url, r.err)
	} else {
		log.Printf("%-20s %s", r.url, r.latency)
	}
	time.Sleep(5 * time.Second)
	log.Println("quit anyways ...", runtime.NumGoroutine(), "still running")
}
