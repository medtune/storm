package filters

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"

	"github.com/medtune/storm/features"
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

type ImageProcessor struct {
	Filters         []ImageFilter
	DefaultDataKey  string
	DefaultEncoding string
	DefaultFeatures map[string]*features.Feature
}

func NewImgProcs() *ImageProcessor {
	return &ImageProcessor{DefaultFeatures: make(map[string]*features.Feature)}
}

func NewImgProcsWithFilters(fs ...ImageFilter) *ImageProcessor {
	return &ImageProcessor{Filters: fs}
}

func (ip *ImageProcessor) SetDefaultKey(s string) {
	ip.DefaultDataKey = s
}

func (ip *ImageProcessor) AddFeature(ft string, f *features.Feature) error {
	_, ok := ip.DefaultFeatures[ft]
	if ok && !AllowOverWritingFeature {
		return fmt.Errorf("Cannot overwrite Default feature (%v), overwrting rule is %v", "image", ft, AllowOverWritingFeature)
	}
	ip.DefaultFeatures[ft] = f
	return nil
}

func (ip *ImageProcessor) SetEncoding(encoding string) error {
	if !isSupported(encoding) {
		return fmt.Errorf("Unsupported encoding format %v", encoding)
	}
	ip.DefaultEncoding = encoding
	return nil
}

func (ip *ImageProcessor) AddFilter(fs interface{}) {
	ip.Filters = append(ip.Filters, fs.(ImageFilter))
}

func (ip *ImageProcessor) Process(r io.ReadCloser, kind string, extraFeatures map[string]*features.Feature) (*features.Features, error) {
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

	for _, ofilter := range ip.Filters {
		img = ofilter(img)
	}

	var outtype string
	if ip.DefaultEncoding != "" {
		outtype = ip.DefaultEncoding
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
	imgfeature := &features.Feature{
		Kind: &features.Feature_BytesList{BytesList: &features.BytesList{
			Value: [][]byte{bif},
		}},
	}

	size := 0
	if extraFeatures != nil {
		size += len(extraFeatures)
	}
	if ip.DefaultFeatures != nil {
		size += len(ip.DefaultFeatures)
	}

	m := make(map[string]*features.Feature, size)
	fts := &features.Features{Feature: m}
	for i, v := range ip.DefaultFeatures {
		fts.Feature[i] = v
	}

	_, ok := fts.Feature[ip.DefaultDataKey]
	if ok && !AllowOverWritingFeature {
		return nil, fmt.Errorf("Cannot overwrite Default feature (%v), overwrting rule is %v", "image", AllowOverWritingFeature)
	}

	fts.Feature[ip.DefaultDataKey] = imgfeature
	for i, v := range extraFeatures {
		_, ok := fts.Feature[i]
		if ok && !AllowOverWritingFeature {
			return nil, fmt.Errorf("Cannot overwrite extra feature (%v) , overwrting rule is %v", i, AllowOverWritingFeature)
		}
		fts.Feature[i] = v
	}

	return fts, nil
}
