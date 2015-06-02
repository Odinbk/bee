package main

import "fmt"
import "sync"
import "time"

func main() {
	c1 := make(chan int)
	c2 := make(chan int, 5)
	var wg sync.WaitGroup

	go func() {
		for i := 0; i < 10; i++ {
			c1 <- i
		}
	}()

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case c := <-c1:
					c2 <- c
				case <-time.After(time.Second):
					return
				}
			}
		}()
	}

	for i := 0; i < 4; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				select {
				case c := <-c2:
					fmt.Println(c)
				case <-time.After(time.Second):
					return
				}
			}
		}()
	}

	fmt.Println("Done")
	wg.Wait()
}
