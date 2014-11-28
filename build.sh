#!/bin/bash -e

if [[ $1 = "-loc" ]]; then
    cd $(dirname $0)/servant; make clean; cd -
    find . -name '*.go' | xargs wc -l | sort -n
    exit
fi

VER=0.2.1stable
ID=$(git rev-parse HEAD | cut -c1-7)
LDFLAGS="-X github.com/funkygao/golib/server.VERSION $VER -X github.com/funkygao/golib/server.BuildID $ID -w"
FAE_HOME=/sgn/app/fae

if [[ $1 = "-dw" ]]; then
    cp -f servant/gen-php/fun/rpc/* /Users/gaopeng/fun/dragon-server-code/v2/fae
    exit
fi

if [[ $1 = "-install" ]]; then
    mkdir -p $FAE_HOME/bin $FAE_HOME/var $FAE_HOME/etc
    cp -f bin/faed.linux $FAE_HOME/bin/faed
    cp -f etc/faed.cf.sample $FAE_HOME/etc/faed.cf
    cp -f etc/faed /etc/init.d/faed
    echo 'update config: metrics_logfile' 
    echo 'Done'
    exit
fi

cd $(dirname $0)/servant; make
cd ../daemon/faed

if [[ $1 = "-linux" ]]; then
    #cd $GOROOT/src 
    #sudo CGO_ENABLED=0 GOOS=linux GOARCH=amd64 ./make.bash
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags $LDFLAGS
    exit
else
    #go build -race -v -ldflags $LDFLAGS
    go build -ldflags $LDFLAGS
fi

#---------
# show ver
#---------
./faed -version
