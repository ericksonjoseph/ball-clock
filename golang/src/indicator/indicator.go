package indicator

import (
	"log"
	"os"
	"time"
)

var (
	slowdown time.Duration = 0
	//logger                 = log.New(&bytes.Buffer{}, "logger: ", log.Ldate)
	logger = log.New(os.Stdout, "logger: ", log.Ldate)
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
	Cycles   int       // Keeps track of how many times this indicator breached its capacity
	queue    chan Ball // Channel to release balls to bottom queue
	carry    chan Ball // Channel to carry last ball to next indicator
	stack    stack     // Stack of ball-bearings indicating the time
	capacity int
}

func New(name string, cap int, queue, carry chan Ball) indicator {

	logger.Println("creating new indicator")

	return indicator{
		name,
		make(chan Ball),
		0,
		queue,
		carry,
		stack{},
		cap,
	}
}

func (i indicator) Run() {
	go func() {

		for {

			logger.Println(".")

			ball := <-i.Track
			i.stack = i.stack.push(ball)
			logger.Printf("Indicator %s <---  %+v. Current count: %d\n", i.Name, ball, len(i.stack))

			if len(i.stack) > i.capacity {

				i.Cycles++

				var ballToCarry Ball
				logger.Printf("%s Full !!\n", i.Name)
				time.Sleep(slowdown * time.Millisecond)

				i.stack, ballToCarry = i.stack.pop()

				for len(i.stack) > 0 {
					var ballToRelease Ball
					i.stack, ballToRelease = i.stack.pop()
					logger.Printf("%s Releasing ball %d\n", i.Name, ballToRelease.Number)
					i.queue <- ballToRelease
					time.Sleep(slowdown * time.Millisecond)
				}

				logger.Printf("%s Carrying ball %d to next indicator\n", i.Name, ballToCarry.Number)
				i.carry <- ballToCarry
			}
		}
	}()
}
