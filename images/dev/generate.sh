#!/usr/bin/env bash

convert shapes/black.png frame/frame.png -composite test.png
convert test.png frame/header.png -composite test.png
convert test.png frame/ring.png -composite test.png
convert test.png staves/zero.png -composite test.png
convert test.png numbers/0.png -composite test.png
convert test.png reality/gloss.png -composite test.png
convert test.png reality/grunge-1.png -composite test.png
