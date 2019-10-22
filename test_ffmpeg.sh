#!/bin/bash

ffmpeg -f lavfi -re -i smptebars=duration=6000:size=320x200:rate=30 -f lavfi -i sine=frequency=1000:duration=6000:sample_rate=48000 -pix_fmt yuv420p -c:v libx264 -b:v 180k -g 60 -keyint_min 60 -profile:v baseline -preset veryfast -c:a aac -b:a 96k -f mpegts tcp://127.0.0.1:9555