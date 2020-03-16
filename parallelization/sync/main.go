package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
)

var wg sync.WaitGroup

func main() {
	if len(os.Args) != 3 {
		fmt.Println("usage:", os.Args[0], "data-count", "parallelization-count")
		os.Exit(1)
	}

	in := make(chan int)
	p, _ := strconv.Atoi(os.Args[2])
	wg.Add(p)
	for i := 0; i < p; i++ {
		go process(i, in)
	}

	n, _ := strconv.Atoi(os.Args[1])
	for i := 0; i < n; i++ {
		in <- i
	}
	close(in)

	wg.Wait()

}

func process(i int, in <-chan int) {
	for elem := range in {
		fmt.Printf("goroutine :%d processing: %d\n", i, elem)
	}
	wg.Done()
}
