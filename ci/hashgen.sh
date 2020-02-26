#!/bin/sh

for f in bin/arkade*; do shasum -a 256 $f > $f.sha256; done
