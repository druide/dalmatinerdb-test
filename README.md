# Dalmatinerdb test

## Setup

See [setup.sh](./ddb-setup/setup.sh).

## Result

AWS m4.4xlarge instance.

```
SELECT sum(sum(domain00001.'00001'.'00001'.start BUCKET 'test2', 1h)) LAST 1d
8.394ms 3.6K metrics

SELECT sum(sum(domain00001.'00001'.'00001'.* BUCKET 'test2', 1h)) LAST 1d
78.037ms 200K metrics

SELECT sum(sum(domain00001.'00001'.*.start BUCKET 'test2', 1h)) LAST 1d
206.475ms 500K metrics

SELECT sum(sum(domain00001.*.*.start BUCKET 'test2', 1d)) LAST 1d
4125.99ms 200M metrics

--------------------------------------------------------------------------------

./haggar -agents=50 -carbon="X.X.X.X:5555" -flush-interval=1s -jitter=1s -metrics=10 -prefix="test2"

100 agents
500K mps
CPU 75%
19GB mem
2e9 metrics
HDD 156Mb
```

## TODO

- When CPU load is above 80%, memory usage is growing fast, dalmatiner frontend
  cannot query due to timeouts. This continue during 5-10 min even after stop
  writing metrics from "haggar".
- Metric queries `a.b.*` are fast (milliseconds), but `a.*.b` are slow
  (up to timeouts). Probably because of big data scan.

```
./haggar -agents=50 -carbon="127.0.0.1:5555" -flush-interval=1s -jitter=1s -metrics=10 -prefix="test2"


SELECT 'dalmatinerdb@127.0.0.1'.'mps' BUCKET 'dalmatinerdb' LAST 30s

SELECT sum('domain'.'domain9'.'aid'.'9'.'cmp'.'9'.'event'.'request'.count BUCKET test1) LAST 1m

SELECT domain001.'001'.'001'.start BUCKET 'test2' LAST 1m


curl -H 'accept: application/json' 'http://127.0.0.1:8080/collections/my_opentsdb_bucket/metrics/AmExB2RvbWFpbjABMAEwCGNvbXBsZXRl/namespaces//tags'

SELECT ALL FROM teststat WHERE ddb:part_1 = 'a1' AND ddb:part_2 = 'domain1' LAST 60
SELECT sum(ALL FROM teststat WHERE ddb:part_1 = 'a1') LAST 60
```
