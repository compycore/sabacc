#!/usr/bin/env bash

cp frames/*.png .
mogrify -resize 100 *.png
apngasm *.png -o d-0.png
