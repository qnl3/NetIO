#!/bin/bash

targetHost=$1
portNumber=$2

if [ -z $targetHost ] ; then
  echo "usage:"
  echo 
  echo "$( basename $0) <server ip> <port>"
  echo 
  echo "$( basename $0) generates a NetIO configuration script targeting the specified host."
  echo 
fi

if [ $# -lt 3 ] ; then
  portNumber=54321
fi

jsonnet hello.jsonnet  -o  hello.json --ext-str "targetHost=$targetHost" --ext-str "portNumber=$portNumber"


