'use strict'

const LOG = console

// Bridge to handle all logging using singleton logger instance and giving us
// the opportunity to decorate each log e.g. with an aws request ID which is
// very helpful when analyzing cloudwatch logs
module.exports = class Logger {
  static trace () {
    LOG.trace.apply(LOG, arguments)
  }

  static debug () {
    //LOG.debug.apply(LOG, arguments)
  }

  static info () {
    LOG.info.apply(LOG, arguments)
  }

  static warn () {
    LOG.warn.apply(LOG, arguments)
  }

  static error () {
    LOG.error.apply(LOG, arguments)
  }

  static fatal () {
    LOG.fatal.apply(LOG, arguments)
  }
}
