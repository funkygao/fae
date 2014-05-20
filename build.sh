#!/bin/bash -e

if [[ $1 = "-loc" ]]; then
    cd $(dirname $0)/servant; make clean; cd -
    find . -name '*.go' | xargs wc -l | sort -n
    exit
fi

ID=$(git rev-parse HEAD | cut -c1-7)
cd $(dirname $0)/servant; make
cd ../daemon/faed

if [[ $1 = "-linux" ]]; then
    cp -f ../../servant/gen-php/fun/rpc/* /Users/gaopeng/fun/royalstory-server-code/system/fae/
    cp -f ../../servant/gen-php/fun/rpc/* /Users/gaopeng/fun/dragon-server-code/v2/fae
    # cd $GOROOT/src; CGO_ENABLED=0 GOOS=linux GOARCH=386 ./make.bash
    CGO_ENABLED=0 GOOS=linux GOARCH=386 go build -ldflags "-X github.com/funkygao/fae/engine.BuildID $ID"
    cp -f faed /Users/gaopeng/fun/royalstory-server-code/daemon/
    exit
else
    #go build -race -v -ldflags "-X github.com/funkygao/fae/engine.BuildID $ID"
    go build -ldflags "-X github.com/funkygao/fae/engine.BuildID $ID -w"
fi

#---------
# show ver
#---------
./faed -version
