const argv = require('minimist')(process.argv.slice(2))
const graphite = require('graphite')
const client = graphite.createClient('tcp://localhost:2003/')

const events = ['request', 'impression', 'creativeView', 'start', 'firstQuartile', 'midpoint', 'thirdQuartile',
  'complete']
const elen = events.length

const id = '' + argv.id
const n = '' + argv.n

setInterval(send, 1000)

let metrics = {}

function send () {
  let domain, aid, cmp, event
  for (var i = 0; i < n; i++) {
    domain = 'domain' + i
    for (var j = 0; j < n; j++) {
      aid = '' + j
      for (var k = 0; k < n; k++) {
        cmp = '' + k
        for (var h = 0; h < elen; h++) {
          event = events[h]

          let value = (Date.now() / 1000) % 10 ? 1 : 0
          sendMetric(`a${id}.${domain}.${aid}.${cmp}.${event}`, value)
        }
      }
    }
  }
}

function sendMetric (metric, value) {
  metrics[metric] = value
  if (Object.keys(metrics).length >= 100) {
    client.write(metrics, (err) => {
      if (err) console.error(err.message)
    })
    metrics = {}
  }
}

console.log(`Worker ${id} send ${n * n * n * elen} events`)
