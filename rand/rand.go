package rand

import (
	"context"
	"math/rand"
	"time"

	"github.com/narahari92/vaali"
)

func Generator() vaali.RunnableFunc {
	return func(ctx context.Context, stop <-chan struct{}) {
		randNums := make([]int, 0)
		rand.Seed(time.Now().UnixNano())
		for {
			select {
			case <-stop:
				randNums = nil
				return
			case <-ctx.Done():
				randNums = nil
				return
			default:
				randNums = append(randNums, rand.Intn(1000))
			}
		}
	}
}
