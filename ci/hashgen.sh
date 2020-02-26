#!/bin/sh

for f in bin/bazaar*; do shasum -a 256 $f > $f.sha256; done
