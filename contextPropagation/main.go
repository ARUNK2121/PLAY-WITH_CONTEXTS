package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type output struct {
	count int
	err   error
}

func dbTask1(ctx context.Context, wg *sync.WaitGroup) (int, error) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		fmt.Println("access DB task 1 error:", ctx.Err())
	case <-time.After(5 * time.Second):
		return 20, nil
	}

	return 3, nil
}

func dbTask2(ctx context.Context, wg *sync.WaitGroup) (int, error) {
	defer wg.Done()
	select {
	case <-ctx.Done():
		fmt.Println("access db task2 error:", ctx.Err())
	case <-time.After(5 * time.Second):
		return 30, nil
	}
	return 3, nil
}

func WebAPI(ctx context.Context) (int, error) {
	wg := sync.WaitGroup{}
	outputCh1 := make(chan output, 1)
	outputCh2 := make(chan output, 2)
	wg.Add(1)
	go func() {
		count1, err := dbTask1(ctx, &wg)
		o := output{
			count: count1,
			err:   err,
		}
		outputCh1 <- o
	}()

	wg.Add(1)
	go func() {
		count2, err := dbTask2(ctx, &wg)
		o := output{
			count: count2,
			err:   err,
		}
		outputCh2 <- o
	}()

	wg.Wait()
	output1 := <-outputCh1
	if output1.err != nil {
		return 0, output1.err
	}

	output2 := <-outputCh2
	if output2.err != nil {
		return 0, output2.err
	}

	output := output1.count + output2.count
	return output, nil

}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 7*time.Second)
	defer cancel()
	count, err := WebAPI(ctx)
	if err != nil {
		fmt.Println("long running task exited with error:", err)
		return
	}
	fmt.Println("count is:", count)
}
