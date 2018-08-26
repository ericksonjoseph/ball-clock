'use strict'

const indicator = require('./indicator/indicator.js');
const log = require('./logger.js');


function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

// User defined
let ballCountArg = 30
let timeArg = 60000

let ballCount = ballCountArg
let slowdown = timeArg

if (ballCount < 0 || ballCount > 127) {
    log.error("Ball count must be in the range of 27 to 127. %d given\n", ballCount)
    process.exit(1)
}

// System variables
var consecutive = 0;
var currentCycle = 0;
let firstBallNumber = 1
let prevBallNumber = firstBallNumber

let queue = []
let hour = new indicator("Hour", 11, queue, queue)
let fiveMin = new indicator("Five", 11, queue, hour.Track)
let min = new indicator("Min", 4, queue, fiveMin.Track)


// Push inital balls to bottom queue
for (let i = firstBallNumber; i < (firstBallNumber + ballCount); i++) {
    log.debug("pushing ball %d to the queue\n", i)
    queue.push(indicator.Ball(i))
}

// Every minute, send a ball to the minute indicator
while(true) {

    // Run the indicators
    min.Run()
    fiveMin.Run()
    hour.Run()

    // Grab the next ball from the queue
    let i = queue.shift()
    if (i === undefined) {
        console.log("Queue empty")
        process.exit(2)
        continue
    }

    // Check to notice if this ball is in the original order with the previous one
    if (i.Number == prevBallNumber+1) {
        log.debug("++\n")
        consecutive++
    } else {
        consecutive = 0
    }

    prevBallNumber = i.Number

    if (consecutive == ballCount-1) {
        currentCycle++
        log.debug("--------- Cycle %d started ---------\n", currentCycle)
        if (currentCycle == 2) {
            log.debug("%d balls cycle after %d days\n", ballCount, hour.Cycles / 2)
            log.debug("Completed in %d milliseconds (%.3f seconds)\n")
            process.exit(0)
        }
        consecutive = 0
        prevBallNumber = firstBallNumber
    }

    // Send this ball to the first indicator
    min.Track.push(i)

    //time.Sleep(slowdown * time.Millisecond)
}
