package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var sliceWork []func() error
	for i := 1; i < 50; i++ {
		work := work(i)
		sliceWork = append(sliceWork, work)
	}
	doWork(sliceWork, 10, 200)
}

func doWork(fs []func() error, n int, nError int) {
	chError := make(chan error, nError)
	var wg sync.WaitGroup
	quit := make(chan bool)
	for _, f := range fs {
		for i := 0; i < n; i++ {
			wg.Add(1)
			go func() {
				defer func() { wg.Done() }()
				for {
					select {
					case <-quit:
						return
					default:
						chError <- f()
					}
				}
			}()
		}
	}
	var count int
	for err := range chError {
		fmt.Println(err)
		count++
		if count == nError {
			close(quit)
			break
		}
	}

	fmt.Println("after max error")
	go func() {
		wg.Wait()
		close(chError)
	}()
	for err := range chError {
		fmt.Println(err)
	}
	fmt.Println("all done")
}

func work(x int) func() error {
	return func() error {
		for i := 1; i > 0; i++ {
			if i%x*1000 == 0 {
				time.Sleep(1000 * time.Millisecond)
				return fmt.Errorf("fail work: %d", x)
			}
		}
		return nil
	}
}
