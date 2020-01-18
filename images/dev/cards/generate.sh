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

echo "Generating cards"

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

cp full/back.png .

echo "Cropping cards"
mogrify -crop 1100x1600+0+0 -gravity Center *.png
echo "Resizing cards"
mogrify -resize x300 *.png

echo "Removing old card images"
rm ../../cards/*.png
echo "Installing new card images"
mv *.png ../../cards/
