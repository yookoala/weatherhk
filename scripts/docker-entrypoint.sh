#!/bin/bash

case $1 in

  "test" )
    go test -v -cover ./...
    ;;

  "*" )
    weatherhk-server "$@"
    ;;

esac
