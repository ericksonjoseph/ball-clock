'use strict'

const indicator = require('./indicator/indicator.js');
const log = require('./logger.js');

function main() {

    // User defined
    let ballCount = parseInt(process.argv[2]) || 0
    let minutes = parseInt(process.argv[3]) || 0

    if (ballCount < 27 || ballCount > 100027) {
        log.error("Ball count must be in the range of 27 to 127. %d given", ballCount)
        process.exit(1)
    }

    // System variables
    var mintuesElapsed = 0;
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
        log.debug("pushing ball %d to the queue", i)
        queue.push(i)
    }

    // Start the timer
    let start = Date.now();
    let end = 0

    // Every minute, send a ball to the minute indicator
    while(true) {

        // Grab the next ball from the queue
        let i = queue.shift()
        if (i === undefined) {
            console.log("Queue empty. Deadlock!")
            process.exit(2)
            continue
        }

        // Check to notice if this ball is in the original order with the previous one
        if (i == prevBallNumber+1) {
            log.debug("++")
            consecutive++
        } else {
            consecutive = 0
        }

        prevBallNumber = i

        if (consecutive == ballCount-1) {
            currentCycle++
            log.debug("Cycle %d started", currentCycle)
            if (currentCycle == 2 && !minutes) {
                end = Date.now();
                break;
            }
            consecutive = 0
            prevBallNumber = firstBallNumber
        }

        // Send this ball to the first indicator
        min.Track.push(i)

        // Run the indicators
        min.Run()
        fiveMin.Run()
        hour.Run()

        mintuesElapsed++
        if (minutes && minutes === mintuesElapsed) {
            log.info({
                "Min": min.stack, 
                "FiveMin": fiveMin.stack, 
                "Hour": hour.stack, 
                "Main": queue, 
            })
            process.exit(1)
        }
    }

    let diff = end - start

    log.info("%d balls cycle after %d days", ballCount, hour.Cycles / 2)
    log.info("Completed in %d milliseconds (%d seconds)", diff, Number.parseFloat((diff / 1000)).toFixed(3))
}

main();
