package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage:", os.Args[0], "data-count", "parallelization-count")
		os.Exit(1)
	}

	in := make(chan int)
	var outs []chan string

	p, _ := strconv.Atoi(os.Args[2])
	for i := 0; i < p; i++ {
		out := make(chan string)
		outs = append(outs, out)
		go process(i, in, out)
	}

	n, _ := strconv.Atoi(os.Args[1])
	go func(n int) {
		for i := 0; i < n; i++ {
			in <- i
		}
		close(in)
	}(n)

	for elem := range merge(outs) {
		fmt.Println("out", elem)
	}

}

func process(i int, in <-chan int, out chan string) {
	for elem := range in {
		out <- fmt.Sprintf("goroutine :%d processing: %d", i, elem)
	}

	close(out)
}

func merge(cs []chan string) chan string {
	out := make(chan string)

	var wg sync.WaitGroup
	wg.Add(len(cs))

	for _, c := range cs {
		go func(c chan string) {
			for v := range c {
				out <- v
			}
			wg.Done()
		}(c)
	}

	go func() {
		wg.Wait()
		close(out)
	}()

	return out
}
