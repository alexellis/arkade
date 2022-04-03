#!/bin/bash

set -e -x -u -o pipefail

go build 

./arkade get dagger --arch arm64 --os darwin --quiet
file $HOME/.arkade/bin/dagger
rm $HOME/.arkade/bin/dagger

./arkade get dagger --arch x86_64 --os darwin --quiet
file $HOME/.arkade/bin/dagger
rm $HOME/.arkade/bin/dagger 

./arkade get dagger --arch x86_64 --os linux --quiet
file $HOME/.arkade/bin/dagger
rm $HOME/.arkade/bin/dagger 

./arkade get dagger --arch arm64 --os linux --quiet
file $HOME/.arkade/bin/dagger
rm $HOME/.arkade/bin/dagger 

./arkade get dagger --arch x86_64 --os ming --quiet
file $HOME/.arkade/bin/dagger
rm $HOME/.arkade/bin/dagger 

