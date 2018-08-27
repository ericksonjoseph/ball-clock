package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"./indicator"
)

var (
	//logger = log.New(&bytes.Buffer{}, "logger: ", log.Ldate)
	logger = log.New(os.Stdout, "logger: ", log.Ldate)
)

func main() {

	// User defined
	ballCount, err := strconv.Atoi(os.Args[1])

	if err != nil || ballCount < 27 || ballCount > 127 {
		fmt.Printf("Ball count must be in the range of 27 to 127. %d given\n", ballCount)
		os.Exit(1)
	}

	// System variables
	var consecutive int
	var currentCycle int
	firstBallNumber := 1
	prevBallNumber := firstBallNumber

	queue := indicator.NewQueue(ballCount)
	hour := indicator.New("Hour", 11, queue, queue)
	fiveMin := indicator.New("Five", 11, queue, hour)
	min := indicator.New("Min", 4, queue, fiveMin)

	// Push inital balls to bottom queue
	for i := firstBallNumber; i < (firstBallNumber + ballCount); i++ {
		logger.Printf("pushing ball %d to the queue\n", i)
		queue.Push(indicator.Ball{i})
	}

	// Start the timer
	start := time.Now().UTC()

	// Every minute, send a ball to the minute indicator
	for {
		// Grab the next ball from the queue
		i := queue.Next()

		// Check to notice if this ball is in the original order with the previous one
		if i.Number == prevBallNumber+1 {
			logger.Printf("++\n")
			consecutive++
		} else {
			consecutive = 0
		}

		prevBallNumber = i.Number

		if consecutive == ballCount-1 {
			currentCycle++
			logger.Printf("--------- Cycle %d started ---------\n", currentCycle)
			if currentCycle == 2 {
				runTime := time.Since(start)
				fmt.Printf("%d balls cycle after %d days\n", ballCount, (hour.Cycles / 2))
				fmt.Printf("Completed in %d milliseconds (%.3f seconds)\n", runTime/time.Millisecond, (float64(runTime) / float64(time.Second)))
				os.Exit(0)
			}
			consecutive = 0
			prevBallNumber = firstBallNumber
		}

		// Send this ball to the first indicator
		min.Push(i)

		// Run the indicators
		min.Run()
		fiveMin.Run()
		hour.Run()
	}
}
