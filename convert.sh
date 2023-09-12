#/bin/bash

mv "$1.ts" all.ts

ffmpeg -i all.ts -acodec copy -vcodec copy all.mp4

mv all.mp4 "$1.mp4"

rm all.ts