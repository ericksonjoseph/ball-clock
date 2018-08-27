'use strict'

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

        let ball = this.Track.pop()
        if (ball === undefined) {
            return
        }

        this.stack.push(ball)

        if (this.stack.length > this.capacity) {

            this.Cycles++

            let ballToCarry = null;

            ballToCarry = this.stack.pop()

            while (this.stack.length > 0)  {
                var ballToRelease = null;
                ballToRelease = this.stack.pop()
                this.queue.push(ballToRelease)
            }

            this.carry.push(ballToCarry)
        }
    }
}
