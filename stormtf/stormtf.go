package stormtf

import (
	"context"
	"fmt"
	"image"
	"io"
	"sync"

	"strings"

	"github.com/anthonynsimon/bild/transform"
)

type errorHandler func(error)

type StormTF struct {
	/*itemCount     int16
	maxThread     int8
	bucketSize    int8
	maxGoroutines int8*/

	errorHandler errorHandler

	googleSearch   GoogleSearchEngineService
	downloader     func(context.Context, string) (io.ReadCloser, error)
	imageProcessor ImageProcessor
	writer         *writer
}

func getImgType(s string) string {
	if strings.Contains(s, "jpeg") {
		return "jpeg"
	}
	if strings.Contains(s, "png") {
		return "png"
	}
	return "unkown"
}

func New(gs GoogleSearchEngineService) (*StormTF, error) {
	fs := make(map[string]*Feature)
	fs["label"] = &Feature{
		Kind: &Feature_BytesList{BytesList: &BytesList{
			Value: [][]byte{[]byte("dog")},
		}},
	}
	ip := &imageProcess{
		defaultFeatures: fs,
		filters: []filter{filter(func(img image.Image) image.Image {
			return transform.Resize(img, 256, 256, transform.Linear)
			//return resize.Resize(256, 256, img, resize.Lanczos3)
		})},
	}
	return &StormTF{
		writer:         &writer{},
		googleSearch:   gs,
		imageProcessor: ip,
		errorHandler:   func(e error) { fmt.Println("--------", e) },
		downloader:     DownloadBodyRC,
	}, nil
}

func (stf *StormTF) Storm(ctx context.Context, query string, queryOption QueryOption, id string,
	numSamples int64, destination string) error {
	if numSamples%10 != 0 {
		return fmt.Errorf("numSamples must be in form 10 * k, got :%v", numSamples)
	}
	err := stf.writer.Init(destination, numSamples, stf.errorHandler)
	if err != nil {
		return err
	}
	stf.googleSearch.SetEngineID(id)
	opt := queryOption
	var start int64 = 0
	for start*10 < numSamples {
		//fmt.Println("called loop")
		opt.Start = 20
		search, err := stf.googleSearch.Search(ctx, query, &opt)
		if err != nil {
			return err
		}
		//fmt.Println("searched and found", len(search.Items))
		var wg sync.WaitGroup
		fmt.Println("GOT;", len(search.Items))
		for index, _ := range search.Items {
			i := *search.Items[index]
			fmt.Println(i.Link)

			wg.Add(1)
			go func() {
				defer wg.Done()
				//fmt.Println(i.Link)
				fmt.Println(i.Mime)
				b, err := stf.downloader(ctx, i.Link)
				if err != nil {
					return
				}
				//fmt.Println("will process image", item.Mime)
				kind := getImgType(i.Mime)
				if kind == "unkown" {
					return
				}
				ft, err := stf.imageProcessor.Process(b, kind)
				//fmt.Println("got ft", err)
				if err != nil {
					return
				}
				fmt.Println("PASSED")
				//fmt.Println("processed and the winner izzzzzzz", err)
				stf.writer.wprotoChan <- ft
				fmt.Println("PASSED2222")

				//fmt.Println("lool..")

			}()
		}
		//fmt.Println("waiting...")
		wg.Wait()
		//fmt.Println("freeeeed.....")
		start++
	}
	//fmt.Println("ended......")
	stf.writer.Close()
	//fmt.Println("xd-------")
	return nil
	//done := make(chan struct{}, 1)
}

/*
stormtf --labels=cat;kity/dog --proto-format=features --resize=64*64 -o catdogs.tfrecord

stormtf --query=cat:cat+kitten/dog:dog --image=true
stormtf [-q --query] [-i --image] [-r --image-resize] [-o --output] [-v --verbose] [-p --proto-format]

*/
