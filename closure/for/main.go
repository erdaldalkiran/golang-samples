package main

import "fmt"

func main() {
	fmt.Println("value:")
	j := 0
	for i := 0; i < 3; i++ {
		j := j
		fmt.Printf("%v ", j)
		j++
	}

	fmt.Println()
	fmt.Println("pointer:")

	z := new(int)
	*z = 3
	for i := 0; i < 3; i++ {
		var z = z
		fmt.Printf("%v ", *z)
		z = nil
	}

	fmt.Println()
	fmt.Println("channel:")

	c := make(chan int)
	go func() {
		for i := 0; i < 3; i++ {
			var c = c
			c <- i
			fmt.Printf("channel: %v ", c)
			c = nil
		}

		close(c)
	}()
	for i := range c {
		fmt.Printf("channel value: %v ", i)
	}

	fmt.Println()
	fmt.Println("slice:")
	s := []int{1, 2, 3}
	for i := 0; i < 3; i++ {
		var s = s
		for v := range s {
			fmt.Printf("%v ", v)
		}
		s = nil
	}

	fmt.Println()
	fmt.Println("map:")
	m := map[int]int{
		0: 0,
		1: 1,
		2: 2,
	}
	for i := 0; i < 3; i++ {
		var m = m
		for k := range m {
			fmt.Printf("%v ", k)
		}
		m = nil
	}
}
