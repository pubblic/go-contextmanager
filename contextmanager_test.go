package contextmanager

import (
	"context"
	"fmt"
)

func ExampleSignalContext() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	c := SignalTx()

	producer := func() error {
		co := NewSignalContext(ctx, c)
		data := []int{1, 2, 3, 4, 5}
		for _, n := range data {
			if !co.Yield(n) {
				return co.Err()
			}
		}
		return nil
	}

	process := func(n int) error {
		// do something with the received n
		return nil
	}

	consumer := func() error {
		for s := range c {
			var n int
			if !s.As(&n) {
				return s.Err()
			}
			err := process(n)
			if err != nil {
				return err
			}
		}
		return nil
	}

	go func() {
		producer()
		close(c)
	}()

	err := consumer()
	cancel()
	if err == nil {
		fmt.Println("No Error")
	} else {
		fmt.Println("Error")
	}

	// Output: No Error
}
