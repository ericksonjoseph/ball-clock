'use strict'

const log = require('../logger.js');
const slowdown = 0

class Ball {
    constructor(n) {
        this.Number = n
    }
}

function sleep(ms) {
    return new Promise(resolve => setTimeout(resolve, ms));
}

module.exports = class indicator {
    constructor(name, cap, queue, carry) {
        this.Name = name
        this.Track = []
        this.Cycles = 0
        this.queue = queue
        this.carry = carry
        this.stack = []
        this.capacity = cap
    }

    static Ball(n) {
        return new Ball(n)
    }

    Run() {

        log.debug(".")

        let ball = this.Track.pop()
        if (ball === undefined) {
            return
        }

        this.stack.push(ball)
        log.debug("Indicator " + this.Name + " <---  Ball(" + ball.Number + ") Current count: " + this.stack.length + "\n")

        if (this.stack.length > this.capacity) {

            this.Cycles++

            let ballToCarry = null;
            log.debug(this.Name + " Full !!\n")
            sleep(slowdown)

            ballToCarry = this.stack.pop()

            while (this.stack.length > 0)  {
                var ballToRelease = null;
                ballToRelease = this.stack.pop()
                log.debug(this.Name + " Releasing ball " + ballToRelease.Number)
                this.queue.push(ballToRelease)
                sleep(slowdown)
            }

            log.debug(this.Name + " Carrying ball %d to next indicator" + ballToCarry.Number)
            this.carry.push(ballToCarry)
        }
    }
}
