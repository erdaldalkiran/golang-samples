package main

import (
	"cache-invalidator/invalidator"
	"fmt"
)

func main() {
	inv := invalidator.New()
	fmt.Println(inv)

	for i := 0; i < 100; i++ {
		go run(inv)
	}

	done := make(chan bool)
	<-done
}

func run(inv *invalidator.Invalidator) {
	for {
		err := inv.Invalidate("ciko-tanesi")
		if err != nil {
			fmt.Println(err)
		}
	}
}
