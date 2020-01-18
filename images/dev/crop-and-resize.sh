#!/usr/bin/env bash

mogrify -crop 1100x1600+0+0 -gravity Center *.png
mogrify -resize 20% *.png
