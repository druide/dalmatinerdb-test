const Dalmatiner = require('dalmatiner')
const argv = require('minimist')(process.argv.slice(2))
const conn = new Dalmatiner('tcp://localhost:5555', argv.bucket, 1) // change to server IP

const events = ['request', 'impression', 'creativeView', 'start', 'firstQuartile', 'midpoint', 'thirdQuartile',
  'complete']
const elen = events.length

const id = '' + argv.id
const n = '' + argv.n

setInterval(send, 1000)

function send () {
  var domain, aid, cmp, event
  for (var i = 0; i < n; i++) {
    domain = 'domain' + i
    for (var j = 0; j < n; j++) {
      aid = '' + j
      for (var k = 0; k < n; k++) {
        cmp = '' + k
        for (var h = 0; h < elen; h++) {
          event = events[h]

          let value = Math.floor(Date.now() / 1000) % 10 ? 1 : 0
          conn.sendData(['a' + id, domain, aid, cmp, event], [value])
        }
      }
    }
  }
  // console.log('Send', n * n * n * elen, 'events')
}

console.log(`Worker ${id} send ${n * n * n * elen} events`)
