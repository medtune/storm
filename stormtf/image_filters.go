package stormtf

import (
	"fmt"
	"image"
	"strconv"
	"strings"

	"github.com/anthonynsimon/bild/transform"
)

type imageFilter func(image.Image) image.Image

func ResizeImageFilter(widgth, hight int, fl transform.ResampleFilter) imageFilter {
	return imageFilter(func(img image.Image) image.Image {
		return transform.Resize(img, widgth, hight, fl)
	})
}

func ResizeImageFilterFromString(s string) (imageFilter, int, int, error) {
	ls := strings.Split(s, ":")
	kind := ls[0]
	dims := strings.Split(ls[1], "x")
	x, err := strconv.Atoi(dims[0])
	if err != nil {
		return nil, 0, 0, err
	}
	y, err := strconv.Atoi(dims[1])
	if err != nil {
		return nil, 0, 0, err
	}
	var fl transform.ResampleFilter
	if kind == "linear" {
		fl = transform.Linear
	} else if kind == "lancoz" {
		fl = transform.Lanczos
	} else {
		return nil, 0, 0, fmt.Errorf("Unkown resize filter")
	}

	return ResizeImageFilter(x, y, fl), x, y, nil

}

var ResizeImFilter512x512 = ResizeImageFilter(512, 512, transform.Linear)
var ResizeImFilter256x256 = ResizeImageFilter(256, 256, transform.Linear)
var ResizeImFilter128x128 = ResizeImageFilter(128, 128, transform.Linear)
var ResizeImFilter64x64 = ResizeImageFilter(64, 64, transform.Linear)
