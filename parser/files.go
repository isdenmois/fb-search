package main

import (
	"archive/zip"
	"fmt"
	"image"

	"github.com/gen2brain/jpegxl"
)

func getCover(file string, path string) (image.Image, error) {
	return getImage("files/covers/"+file, fileNameWithoutExtension(path))
}

func getImage(archive string, path string) (image.Image, error) {
	fmt.Print("getImage archive=", archive, "; path=", path)

	// Open the ZIP file for reading
	r, err := zip.OpenReader(archive)
	if err != nil {
		return nil, err
	}
	defer r.Close() // Ensure the reader is closed when done

	file, err := r.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	image, err := jpegxl.Decode(file)

	if err != nil {
		return nil, err
	}

	return image, nil
}
