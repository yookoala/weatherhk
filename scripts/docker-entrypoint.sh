#!/bin/bash

case $1 in

  "test")
    go "$@"
    ;;

  *)
    weatherhk-server "$@"
    ;;

esac
