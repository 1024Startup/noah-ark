#!/bin/bash

path=$(cd `dirname $0`; pwd)

case $1 in
    run)
        echo $1
        cd /go/src/noah
        echo "building..."
        go build
        echo "build done"
        ./noah
        ;;
    docker|*)
        echo $1
        docker run -it --rm -v ${path}:/go -p 1024:8080 golang:1.10 /bin/bash -c "/go/run.sh run"
        ;;
esac