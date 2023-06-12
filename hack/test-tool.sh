#!/bin/bash

if [ -z $1 ];
then
  echo "test-tool.sh"
  echo ""
  echo "Downloads the tool for all supported operating systems and"
  echo "CPU architectures. Print the file type for manual checking"
  echo "This often finds issues with PRs that are not caught by"
  echo "unit test alone"
  echo
  echo usage: test-tool.sh TOOLNAME
  exit 1
fi

set -x 

./arkade get $1 --arch arm64 --os darwin --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

./arkade get $1 --arch x86_64 --os darwin --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

./arkade get $1 --arch x86_64 --os linux --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

./arkade get $1 --arch aarch64 --os linux --quiet && file $HOME/.arkade/bin/$1 && rm $HOME/.arkade/bin/$1 && echo 

./arkade get $1 --arch x86_64 --os mingw --quiet && file $HOME/.arkade/bin/$1.exe && rm $HOME/.arkade/bin/$1.exe && echo

