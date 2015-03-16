#!/bin/bash -e

VER=0.1.0alpha
ID=$(git rev-parse HEAD | cut -c1-7)

go build -ldflags "-X github.com/funkygao/golib/server.Version $VER -X github.com/funkygao/golib/server.BuildId $ID -w"

#---------
# show ver
#---------
./configd -version
