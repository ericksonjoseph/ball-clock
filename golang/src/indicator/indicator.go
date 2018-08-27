package indicator

import (
	"log"
	"os"
)

var (
	//logger                 = log.New(&bytes.Buffer{}, "logger: ", log.Ldate)
	logger = log.New(os.Stdout, "logger: ", log.Ldate)
)

type Ball struct {
	Number int
}

type stack []Ball

func (s *stack) push(i Ball) {
	*s = append(*s, i)
}

func (s stack) pop() (stack, Ball) {
	l := len(s)
	return s[:l-1], s[l-1]
}

type queue chan Ball

func (q *queue) Push(ball Ball) {
	*q <- ball
}

func (q *queue) Next() Ball {
	return <-*q
}

func NewQueue(capacity int) *queue {
	q := make(queue, capacity)
	return &q
}

type pusher interface {
	Push(ball Ball)
}

type indicator struct {
	Name     string
	Cycles   int    // Keeps track of how many times this indicator breached its capacity
	queue    *queue // Channel to release balls to bottom queue
	carry    pusher // Channel to carry last ball to next indicator
	stack    stack  // Stack of ball-bearings indicating the time
	capacity int
}

func New(name string, cap int, queue *queue, carry pusher) *indicator {

	logger.Println("creating new indicator")

	return &indicator{
		name,
		0,
		queue,
		carry,
		stack{},
		cap,
	}
}

func (i *indicator) Push(ball Ball) {
	logger.Printf("Indicator %s <---  %+v. Current count: %d\n", i.Name, ball, len(i.stack))
	i.stack.push(ball)
}

func (i *indicator) Run() {

	logger.Println(".")

	if len(i.stack) > i.capacity {

		i.Cycles++

		var ballToCarry Ball
		logger.Printf("%s Full !!\n", i.Name)

		i.stack, ballToCarry = i.stack.pop()

		for len(i.stack) > 0 {
			var ballToRelease Ball
			i.stack, ballToRelease = i.stack.pop()
			logger.Printf("%s Releasing ball %d\n", i.Name, ballToRelease.Number)
			i.queue.Push(ballToRelease)
		}

		logger.Printf("%s Carrying ball %d to next indicator\n", i.Name, ballToCarry.Number)
		i.carry.Push(ballToCarry)
	}
}
