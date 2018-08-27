package indicator

type Stack []int

func (s *Stack) push(i int) {
	*s = append(*s, i)
}

func (s Stack) pop() (Stack, int) {
	l := len(s)
	return s[:l-1], s[l-1]
}

func (i *indicator) Next() int {
	next := i.Stack[0]
	i.Stack = i.Stack[1:]
	return next
}

type pusher interface {
	Push(ball int)
}

type indicator struct {
	Name     string
	Cycles   int    // Keeps track of how many times this indicator breached its capacity
	queue    pusher // Channel to release balls to bottom queue
	carry    pusher // Channel to carry last ball to next indicator
	Stack    Stack  // Stack of ball-bearings indicating the time
	capacity int
}

func New(name string, cap int, queue, carry pusher) *indicator {
	return &indicator{
		name,
		0,
		queue,
		carry,
		Stack{},
		cap,
	}
}

func (i *indicator) Push(ball int) {
	i.Stack.push(ball)
}

func (i *indicator) Run() {

	if len(i.Stack) > i.capacity {

		i.Cycles++

		var ballToCarry int

		i.Stack, ballToCarry = i.Stack.pop()

		for len(i.Stack) > 0 {
			var ballToRelease int
			i.Stack, ballToRelease = i.Stack.pop()
			i.queue.Push(ballToRelease)
		}

		i.carry.Push(ballToCarry)
	}
}
