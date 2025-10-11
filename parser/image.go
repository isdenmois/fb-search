package main

import (
	"bytes"
	"image"
	"image/jpeg"
)

func imageToJpg(img image.Image) (bytes.Buffer, error) {
	var buf bytes.Buffer
	err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 95})

	return buf, err
}
