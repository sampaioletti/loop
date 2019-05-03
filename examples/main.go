package main

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/sampaioletti/loop/pkg/loop"
)

func main() {
	count := make(chan int)
	i := 0
	//create the first loop, wher its created doesn't matter
	l1 := loop.NewLoop()
	l1.AddCall(func(context.Context) {
		fmt.Println("Hello From l1")
		time.Sleep(time.Second * 1)
	})
	go func() {
		l1.AddCall(func(context.Context) {
			fmt.Println("Hello Again From l1")
			time.Sleep(time.Second * 5)
		})
		l1.Run(nil)
	}()
	// Watch for values on the count channel
	go func() {
		for {
			select {
			case n := <-count:
				fmt.Println("count", n)
			}
		}
	}()
	l2 := loop.NewLoop()
	errChan := make(chan error)
	intChan := make(chan int)

	l2.AddCall(func(context.Context) {
		sum, err := add(i, 10)
		if err != nil {
			errChan <- err
			return
		}
		intChan <- sum
	})

	go func() {
		for {
			select {
			case n := <-intChan:
				fmt.Println("add", n)
			case err := <-errChan:
				fmt.Println("add Error", err)
			}
		}
	}()
	l2.AddCall(func(context.Context) {
		fmt.Println("Hello From l2")
		time.Sleep(time.Second * 1)
		count <- i
		i++
	})
	go func() {
		//this should fail
		time.Sleep(time.Second * 2)
		err := l2.Run(nil)
		if err != nil {
			fmt.Println(err)
		}
	}()
	fmt.Println(l2.Run(nil))
}

func add(a, b int) (int, error) {
	if a < 1 {
		return -1, errors.New("Doesn't support adding zero")
	}
	return a + b, nil
}
