'use strict'

const indicator = require('./indicator/indicator.js');

function main() {

    // User defined
    let ballCount = parseInt(process.argv[2]) || 0
    let minutes = parseInt(process.argv[3]) || 0

    if (ballCount < 27 || ballCount > 127) {
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
            process.exit(2)
        }

        // Check to notice if this ball is in the original order with the previous one
        consecutive = (i == prevBallNumber+1) ?  consecutive+1 : 0

        prevBallNumber = i

        if (consecutive == ballCount-1) {
            currentCycle++
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
            console.info({ "Min": min.stack, "FiveMin": fiveMin.stack, "Hour": hour.stack, "Main": queue, })
            process.exit(1)
        }
    }

    let diff = end - start

    console.info("%d balls cycle after %d days", ballCount, hour.Cycles / 2)
    console.info("Completed in %d milliseconds (%d seconds)", diff, Number.parseFloat((diff / 1000)).toFixed(3))
}

main();
