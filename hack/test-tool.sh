#!/bin/bash

set -e -x -u -o pipefail

go build 

./arkade get $1 --arch arm64 --os darwin --quiet
file $HOME/.arkade/bin/$1
rm $HOME/.arkade/bin/$1

./arkade get $1 --arch x86_64 --os darwin --quiet
file $HOME/.arkade/bin/$1
rm $HOME/.arkade/bin/$1 

./arkade get $1 --arch x86_64 --os linux --quiet
file $HOME/.arkade/bin/$1
rm $HOME/.arkade/bin/$1 

./arkade get $1 --arch arm64 --os linux --quiet
file $HOME/.arkade/bin/$1
rm $HOME/.arkade/bin/$1 

./arkade get $1 --arch x86_64 --os ming --quiet
file $HOME/.arkade/bin/$1
rm $HOME/.arkade/bin/$1 
