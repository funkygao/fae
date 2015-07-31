#!/bin/bash -e

if [[ $1 = "-loc" ]]; then
    cd $(dirname $0)/servant; make clean; cd -
    find . -name '*.go' | xargs wc -l | sort -n
    exit
fi

VER=0.3.2stable
GIT_ID=$(git rev-parse HEAD | cut -c1-7)

if [[ $1 = "-php" ]]; then
    cp -f servant/gen-php/fun/rpc/* /Users/gaopeng/myphp/fae
    exit
fi

FAE_HOME=/sgn/app/fae
if [[ $1 = "-install" ]]; then
    mkdir -p $FAE_HOME/bin $FAE_HOME/var $FAE_HOME/etc
    cp -f bin/faed.linux $FAE_HOME/bin/faed
    cp -f etc/faed.cf.sample $FAE_HOME/etc/faed.cf
    cp -f etc/faed /etc/init.d/faed
    echo 'update config: metrics_logfile' 
    echo 'Done'
    exit
fi

BUILD_FLAGS=''
if [[ $1 = "-race" ]]; then
    BUILD_FLAGS="$BUILD_FLAGS -race"
fi
if [[ $1 = "-gc" ]]; then
    BUILD_FLAGS="$BUILD_FLAGS -gcflags '-m=1'"
fi

if [[ $1 = "-cpu" ]]; then
    go tool pprof ./daemon/faed/faed prof/cpu.pprof
    exit
fi
if [[ $1 = "-mem" ]]; then
    #go tool pprof ./daemon/faed/faed prof/mem.pprof
    echo "================="
    echo "allocated objects"
    echo "================="
    go tool pprof -alloc_objects -text http://localhost:9102/debug/pprof/heap
    echo "================="
    echo "inuse objects"
    echo "================="
    go tool pprof -inuse_objects -text http://localhost:9102/debug/pprof/heap
    #go tool pprof -alloc_space -text http://localhost:9102/debug/pprof/heap
    exit
fi

cd $(dirname $0)/servant; make
cd ../daemon/faed

if [[ $1 = "-linux" ]]; then
    #cp -f ../../servant/gen-php/fun/rpc/* /Users/gaopeng/fun/dragon-server-code/v2/fae
    #cd $GOROOT/src 
    #sudo CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ./make.bash
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $BUILD_FLAGS -ldflags "-X github.com/funkygao/golib/server.Version $VER -X github.com/funkygao/golib/server.BuildId $GIT_ID"
    exit
else
    go build $BUILD_FLAGS -tags release -ldflags "-X github.com/funkygao/golib/server.Version $VER -X github.com/funkygao/golib/server.BuildId $GIT_ID -w"
fi

#---------
# show ver
#---------
./faed -version
