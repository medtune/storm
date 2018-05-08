package stormtf

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

const (
	PNG  = "png"
	JPEG = "jpeg"
)

func isSupported(s string) bool {
	return s == PNG || s == JPEG
}

type filter func(img image.Image) image.Image

type ImageProcessor interface {
	SetFilters(...filter)
	SetFeatures(map[string]*Feature)
	SetEncoding(string) error

	Process(io.ReadCloser, string) (*Features, error)
}

type imageProcess struct {
	filters         []filter
	defaultEncoding string
	defaultFeatures map[string]*Feature
}

func newImageProcessor() *imageProcess {
	return &imageProcess{}
}

func newImageProcessorFilters(fs ...filter) *imageProcess {
	return &imageProcess{filters: fs}
}

func (ip *imageProcess) SetFeatures(feas map[string]*Feature) {
	ip.defaultFeatures = feas
}

func (ip *imageProcess) SetEncoding(encoding string) error {
	if !isSupported(encoding) {
		return fmt.Errorf("Unsupported encoding format %v", encoding)
	}
	ip.defaultEncoding = encoding
	return nil
}

func (ip *imageProcess) SetFilters(fs ...filter) {
	ip.filters = fs
}

func (ip *imageProcess) Process(r io.ReadCloser, kind string) (*Features, error) {
	var img image.Image
	if kind == PNG {
		i, err := png.Decode(r)
		//fmt.Println("decoded png")
		if err != nil {
			return nil, err
		}
		img = i
	} else if kind == JPEG {
		i, err := jpeg.Decode(r)
		//fmt.Println("decoded jped")
		if err != nil {
			return nil, err
		}
		img = i
	} else {
		//fmt.Println("unsupport no decode")
		return nil, fmt.Errorf("Unkown image encoding type")
	}
	r.Close()
	for _, ofilter := range ip.filters {
		img = ofilter(img)
	}

	var outtype string
	if ip.defaultEncoding != "" {
		outtype = ip.defaultEncoding
	} else {
		outtype = kind
	}
	//fmt.Println("decided type", outtype)
	var buf bytes.Buffer
	if outtype == JPEG {
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: 100})
		if err != nil {
			return nil, err
		}
	}

	if outtype == PNG {
		err := png.Encode(&buf, img)
		if err != nil {
			return nil, err
		}
	}
	bif := buf.Bytes()
	imgfeature := &Feature{
		Kind: &Feature_BytesList{BytesList: &BytesList{
			Value: [][]byte{bif},
		}},
	}
	fts := &Features{
		Feature: ip.defaultFeatures,
	}
	fts.Feature["image"] = imgfeature
	//fmt.Println("----<<++++", len(bif))
	return fts, nil
}
