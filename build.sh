#!/bin/bash -e
ID=$(git rev-parse HEAD | cut -c1-7)

cd $(dirname $0)/servant; make
cd ../cmd/fproxyd

go build -v -ldflags "-X main.BuildID $ID"

#---------
# show ver
#---------
./fproxyd -version
