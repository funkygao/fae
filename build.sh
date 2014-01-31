#!/bin/bash -e

cd $(dirname $0)/cmd/fproxyd
ID=$(git rev-parse HEAD | cut -c1-7)

go build -v -ldflags "-X main.BuildID $ID"

#---------
# show ver
#---------
./fproxyd -version
