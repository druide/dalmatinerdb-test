const cp = require('child_process')
const argv = require('minimist')(process.argv.slice(2))

const w = argv.w || 10
let id = 1

for (var i = 0; i < w; i++) {
  cp.fork('./worker-graphite', [
    '--id=' + id++,
    '--n=' + (argv.n || 10),
    '--bucket=' + (argv.bucket || 'test1')
  ])
}
