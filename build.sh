#!/bin/bash -e
ID=$(git rev-parse HEAD | cut -c1-7)

cd $(dirname $0)/servant; make
cd ../daemon/faed

go build -ldflags "-X github.com/funkygao/fae/engine.BuildID $ID"

#---------
# show ver
#---------
./faed -version
