#!/usr/bin/env bash

function generate() {
	echo $4
	convert shapes/$2.png frame/frame.png -composite $4
	convert $4 frame/header.png -composite $4
	convert $4 frame/ring.png -composite $4
	convert $4 staves/$1.png -composite $4
	convert $4 numbers/$3.png -composite $4
	convert $4 reality/gloss.png -composite $4
	convert $4 reality/grunge/$(ls reality/grunge | shuf -n 1) -composite $4
}

for COLOR in red green
do
	for STAVE in circle square triangle
	do
		for NUMBER in 1 2 3 4 5 6 7 8 9 10
		do
			FILENAME=$STAVE-$COLOR-$NUMBER.png
			generate $STAVE $COLOR $NUMBER $FILENAME
		done
	done
done

generate zero black 0 zero.png
