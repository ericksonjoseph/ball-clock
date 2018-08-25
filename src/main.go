package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"./indicator"
)

func main() {

	// User defined
	ballCountArg := flag.Int("c", 27, "ball count")
	timeArg := flag.Int64("t", 60000, "Milliseconds to consider 1 min")

	flag.Parse()

	ballCount := *ballCountArg
	slowdown := time.Duration(*timeArg)

	if ballCount < 27 || ballCount > 127 {
		fmt.Printf("Ball count must be in the range of 27 to 127. %d given\n", ballCount)
		os.Exit(1)
	}

	// System variables
	var consecutive int
	var cycleCount int
	firstBallNumber := 1
	prevBallNumber := firstBallNumber

	queue := make(chan indicator.Ball, ballCount)
	hour := indicator.New("Hour", 12, queue, queue)
	fiveMin := indicator.New("Five", 12, queue, hour.Track)
	min := indicator.New("Min", 5, queue, fiveMin.Track)

	// Run the indicators
	min.Run()
	fiveMin.Run()
	hour.Run()

	// Push inital balls to bottom queue
	for i := firstBallNumber; i < (firstBallNumber + ballCount); i++ {
		fmt.Printf("pushing ball %d to the queue\n", i)
		queue <- indicator.Ball{i}
	}

	// Every minute, send a ball to the minute indicator

	for {
		time.Sleep(slowdown * time.Millisecond)

		// Grab the next ball from the queue
		i := <-queue

		// Check to notice if this ball is in the original order with the previous one
		if i.Number == prevBallNumber+1 {
			fmt.Printf("++\n")
			consecutive++
		} else {
			consecutive = 0
		}

		prevBallNumber = i.Number

		if consecutive == ballCount-1 {
			fmt.Printf("--------- Cycle %d started ---------\n", cycleCount)
			cycleCount++
			consecutive = 0
			prevBallNumber = firstBallNumber
		}

		// Send this ball to the first indicator
		min.Track <- i
	}
}
