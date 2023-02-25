#!/bin/bash

if [ -z $1 ];
then
  echo usage: test-tool.sh TOOLNAME
  exit 1
fi

set -x 

./arkade get $1 --arch arm64 --os darwin --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

./arkade get $1 --arch x86_64 --os darwin --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

./arkade get $1 --arch x86_64 --os linux --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

./arkade get $1 --arch arm64 --os linux --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

./arkade get $1 --arch x86_64 --os ming --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

