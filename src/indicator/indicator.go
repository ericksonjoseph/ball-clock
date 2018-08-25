package indicator

import (
	"fmt"
	"time"
)

var (
	slowdown time.Duration = 0
)

type Ball struct {
	Number int
}

type stack []Ball

func (s stack) push(i Ball) stack {
	return append(s, i)
}

func (s stack) pop() (stack, Ball) {
	l := len(s)
	return s[:l-1], s[l-1]
}

type indicator struct {
	Name     string
	Track    chan Ball // Incoming balls from the queue
	queue    chan Ball // Channel to release balls to bottom queue
	carry    chan Ball // Channel to carry last ball to next indicator
	stack    stack
	capacity int
}

func New(name string, cap int, queue, carry chan Ball) indicator {

	fmt.Println("creating new indicator")

	return indicator{
		name,
		make(chan Ball),
		queue,
		carry,
		stack{},
		cap,
	}
}

func (i indicator) Run() {
	go func() {

		for {

			fmt.Println(".")

			ball := <-i.Track
			i.stack = i.stack.push(ball)
			fmt.Printf("Indicator %s <---  %+v. Current count: %d\n", i.Name, ball, len(i.stack))

			if len(i.stack) >= i.capacity {

				var ballToCarry Ball
				fmt.Printf("%s Full !!\n", i.Name)
				time.Sleep(slowdown * time.Millisecond)

				i.stack, ballToCarry = i.stack.pop()

				for len(i.stack) > 0 {
					var ballToRelease Ball
					i.stack, ballToRelease = i.stack.pop()
					fmt.Printf("%s Releasing ball %d\n", i.Name, ballToRelease.Number)
					i.queue <- ballToRelease
					time.Sleep(slowdown * time.Millisecond)
				}

				fmt.Printf("%s Carrying ball %d to next indicator\n", i.Name, ballToCarry.Number)
				i.carry <- ballToCarry
			}
		}
	}()
}
