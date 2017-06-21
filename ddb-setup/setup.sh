# prerequosites
sudo apt-get update
sudo apt-get install zlib1g-dev make erlang build-essential g++
sudo apt-get install postgresql postgresql-contrib

# ZFS
sudo apt-get install zfs
cd / && touch vdev1 && truncate -s 5G vdev1 # change 900G to appropriate allocation
zpool create data /vdev1
zfs create data/dalmatinerdb -o compression=lz4 -o atime=off -o logbias=throughput

export TARGET_DIRECTORY=/usr/local

# dalmatinerdb
cd ~
git clone https://github.com/dalmatinerdb/dalmatinerdb.git
cd dalmatinerdb
git checkout d4f3663a5043c33635e8b1c4e939c9ccc628d426 # Use rel 0.3.0
make all rel
cp -r _build/prod/rel/ddb $TARGET_DIRECTORY
cd $TARGET_DIRECTORY/ddb
mkdir -p /data/dalmatinerdb/etc
mkdir -p /data/dalmatinerdb/log/sasl
mkdir -p /data/dalmatinerdb/db/ring
cp etc/dalmatinerdb.conf.example /data/dalmatinerdb/etc/dalmatinerdb.conf
nano /data/dalmatinerdb/etc/dalmatinerdb.conf # configure and change settings as needed
nano $TARGET_DIRECTORY/ddb/bin/ddb # comment RUNNER_USER
nano $TARGET_DIRECTORY/ddb/bin/ddb-admin # comment RUNNER_USER
./bin/ddb start

# dalmatinerfe
cd ~
git clone https://github.com/dalmatinerdb/dalmatiner-frontend.git
cd dalmatiner-frontend
git checkout 8fcb9299fd72c421daaf684349183817365d5039 # Use rel 0.3.0
make deps all rel
cp -r _build/prod/rel/dalmatinerfe $TARGET_DIRECTORY
cd $TARGET_DIRECTORY/dalmatinerfe
mkdir -p /data/dalmatinerfe/etc
cp etc/dalmatinerfe.conf.example /data/dalmatinerfe/etc/dalmatinerfe.conf
nano /data/dalmatinerfe/etc/dalmatinerfe.conf # check the settings and adjust if needed
su - postgres -c "psql -f /usr/local/dalmatinerfe/lib/dqe_idx_pg-0.3.6/priv/schema.sql"
service postgresql restart
nano $TARGET_DIRECTORY/dalmatinerfe/bin/dalmatinerfe # comment RUNNER_USER
./bin/dalmatinerfe start
curl -H 'accept: application/json' 'http://127.0.0.1:8080/buckets' # []

# dpx
cd ~
git clone https://github.com/dalmatinerdb/ddb_proxy.git
cd ddb_proxy
git checkout b46fed0628e3b3489dd6a2c5301e6ce9e7c84ca5 # Use rel 0.3.1
# Fix Makefile ./rebar3
make all rel
cp -r _build/prod/rel/dpx $TARGET_DIRECTORY
cd $TARGET_DIRECTORY/dpx
mkdir -p /data/dalmatinerpx/etc
cp etc/dalmatinerpx.conf.example /data/dalmatinerpx/etc/dalmatinerpx.conf
nano /data/dalmatinerpx/etc/dalmatinerpx.conf # check the settings and adjust if needed
nano $TARGET_DIRECTORY/dpx/bin/dpx # comment RUNNER_USER
./bin/dpx start

# some checks
zfs get all data/dalmatinerdb | grep compressratio
