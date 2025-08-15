#!/usr/bin/env bash

while true; do
    go build -o _build/gentable ./examples && pkill -f '_build/gentable'
    inotifywait -e attrib $(find . -name '*.go') || exit
done
