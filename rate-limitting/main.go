package main

import (
	"flag"
	"fmt"
	"time"
)

var limit = flag.Int("l", 1, "max number of requests")
var duration = flag.Duration("d", time.Duration(1000000000), "duration")

func main() {
	flag.Parse()
	fmt.Printf("All arguments: %v %v\n", *limit, *duration)
	r := NewRateLimiter(*limit, *duration)

	go func() {
		timer := time.NewTimer(time.Duration(4 * time.Second))
		<-timer.C
		r.Close()
	}()

	for range r.Tokens {
		fmt.Println("sending requets")
		time.Sleep(200 * time.Millisecond)
	}

	time.Sleep(2 * time.Second)
}

//RateLimiter ...
type RateLimiter struct {
	limit    int
	duration time.Duration
	Tokens   chan struct{}
	done     chan bool
}

//NewRateLimiter ...
func NewRateLimiter(l int, d time.Duration) *RateLimiter {
	ch := make(chan struct{}, l)
	for i := 0; i < l; i++ {
		ch <- struct{}{}
	}

	rl := &RateLimiter{l, d, ch, make(chan bool)}
	rl.autoFill()
	return rl
}

//Close ...
func (r *RateLimiter) Close() {
	close(r.done)
}

func (r *RateLimiter) autoFill() {
	defer fmt.Println("closing autoFill")

	fillBucket := func() {
		for i := 0; i < r.limit; i++ {
			select {
			case r.Tokens <- struct{}{}:
				fmt.Println("writing tokens")
			default:
				fmt.Println("tokens bucket is full")
				return
			}
		}
	}

	go func() {
		defer fmt.Println("closing go autoFill")
		defer close(r.Tokens)
		ticker := time.NewTicker(r.duration)
		for {
			select {
			case <-r.done:
				fmt.Println("ending autoFill")
				return
			case <-ticker.C:
				fillBucket()
			}
		}
	}()
}
