package stormtf

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

var AllowOverWritingFeature bool

const (
	PNG    = "png"
	JPEG   = "jpeg"
	UNKOWN = "unkown"

	DefaultJPEGQuality = 98
)

func isSupported(s string) bool {
	return s == PNG || s == JPEG
}

type imageProcessor struct {
	filters         []imageFilter
	defaultDataKey  string
	defaultEncoding string
	defaultFeatures map[string]*Feature
}

func NewImgProcs() *imageProcessor {
	return &imageProcessor{defaultFeatures: make(map[string]*Feature)}
}

func NewImgProcsWithFilters(fs ...imageFilter) *imageProcessor {
	return &imageProcessor{filters: fs}
}

func (ip *imageProcessor) SetDefaultKey(s string) {
	ip.defaultDataKey = s
}

func (ip *imageProcessor) AddFeature(ft string, f *Feature) error {
	_, ok := ip.defaultFeatures[ft]
	if ok && !AllowOverWritingFeature {
		return fmt.Errorf("Cannot overwrite default feature (%v), overwrting rule is %v", "image", ft, AllowOverWritingFeature)
	}
	ip.defaultFeatures[ft] = f
	return nil
}

func (ip *imageProcessor) SetEncoding(encoding string) error {
	if !isSupported(encoding) {
		return fmt.Errorf("Unsupported encoding format %v", encoding)
	}
	ip.defaultEncoding = encoding
	return nil
}

func (ip *imageProcessor) AddFilter(fs interface{}) {
	ip.filters = append(ip.filters, fs.(imageFilter))
}

func (ip *imageProcessor) Process(r io.ReadCloser, kind string, extraFeatures map[string]*Feature) (*Features, error) {
	var img image.Image
	defer r.Close()
	if kind == PNG {
		i, err := png.Decode(r)
		if err != nil {
			return nil, err
		}
		img = i
	} else if kind == JPEG {
		i, err := jpeg.Decode(r)
		if err != nil {
			return nil, err
		}
		img = i
	} else if kind == UNKOWN {
		i, k, err := image.Decode(r)
		if err != nil {
			return nil, err
		}
		img = i
		if k != PNG && k != JPEG {
			kind = JPEG
		} else {
			kind = k
		}
	} else {
		return nil, fmt.Errorf("Unkown image encoding type")
	}

	for _, ofilter := range ip.filters {
		img = ofilter(img)
	}

	var outtype string
	if ip.defaultEncoding != "" {
		outtype = ip.defaultEncoding
	} else {
		outtype = kind
	}

	var buf bytes.Buffer
	if outtype == JPEG {
		err := jpeg.Encode(&buf, img, &jpeg.Options{Quality: DefaultJPEGQuality})
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

	size := 0
	if extraFeatures != nil {
		size += len(extraFeatures)
	}
	if ip.defaultFeatures != nil {
		size += len(ip.defaultFeatures)
	}

	m := make(map[string]*Feature, size)
	fts := &Features{Feature: m}
	for i, v := range ip.defaultFeatures {
		fts.Feature[i] = v
	}

	_, ok := fts.Feature[ip.defaultDataKey]
	if ok && !AllowOverWritingFeature {
		return nil, fmt.Errorf("Cannot overwrite default feature (%v), overwrting rule is %v", "image", AllowOverWritingFeature)
	}

	fts.Feature[ip.defaultDataKey] = imgfeature
	for i, v := range extraFeatures {
		_, ok := fts.Feature[i]
		if ok && !AllowOverWritingFeature {
			return nil, fmt.Errorf("Cannot overwrite extra feature (%v) , overwrting rule is %v", i, AllowOverWritingFeature)
		}
		fts.Feature[i] = v
	}

	return fts, nil
}
