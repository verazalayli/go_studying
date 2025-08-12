package concurrency

import (
	"context"
	"fmt"
	"time"
)

/*
context is an interface from the context package that is used for:
canceling work (for example, in case of timeout, error, or connection closure),
deadline transfers (deadline is the deadline before which the operation must be completed),
storing and transferring values (for example, request ID, trace ID, etc. between layers).
*/
func doWork(ctx context.Context) {
	for {
		select {
		case <-ctx.Done(): // <- канал закрывается, если произошла отмена
			fmt.Println("Cancelled:", ctx.Err())
			return
		default:
			fmt.Println("Working...")
			time.Sleep(500 * time.Millisecond)
		}
	}
}

func contextFunc() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	go doWork(ctx)

	time.Sleep(3 * time.Second)
	fmt.Println("Main done")
}

/*
context.Background()
is the base empty context. It is used as a "root" for others.

context.TODO()
is a placeholder when the context is needed, but you haven't decided which one yet.

context.With Cancel(parent)
A context that can be manually canceled.

context.WithTimeout(parent, duration)
A context that is automatically canceled after a duration.

context.WithDeadline(parent, time.Time)
A context with a fixed deadline.

context.WithValue(parent, key, value)
The context that stores the value.
It is not used for "everything in a row", but only for end-to-end parameters (for example, request ID).
*/
