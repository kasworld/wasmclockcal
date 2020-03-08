#!/usr/bin/env bash

echo "open bgclock http://localhost:8080/?bgimg=image"
echo "open youtube clock http://localhost:8080/?mvid=youtubeid"

python3 -m http.server 8080 
