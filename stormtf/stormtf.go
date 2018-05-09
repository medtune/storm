package stormtf

import (
	"context"
	"fmt"
	"image"
	"io"
	"sync"
	"time"

	"strings"

	"github.com/anthonynsimon/bild/transform"
)

type errorHandler func(error)

type StormTF struct {
	itemCount     int16
	maxThread     int8
	bucketSize    int8
	maxGoroutines int8

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
			return transform.Resize(img, 512, 512, transform.Linear)
			//return resize.Resize(256, 256, img, resize.Lanczos3)
		})},
	}
	return &StormTF{
		writer:         &writer{},
		googleSearch:   gs,
		imageProcessor: ip,
		errorHandler:   func(e error) { logger.Error("ERROR HANDLER: %v", e) },
		downloader:     DownloadBodyRC,
	}, nil
}

func (stf *StormTF) Storm(ctx context.Context, query string, queryOption QueryOption, id string,
	numSamples int64, destination string) error {
	logger.Debug("Storming started\n")

	if numSamples%10 != 0 {
		return fmt.Errorf("numSamples must be in form 10 * k, got :%v", numSamples)
	}

	err := stf.writer.Init(destination, numSamples, stf.errorHandler)
	if err != nil {
		logger.Debug("writer couldnt make it. RIP\n")
		return err
	}

	logger.Debug("Writer is ready\n")
	stf.googleSearch.SetEngineID(id)
	opt := queryOption
	var start int64 = 0
	logger.Debug("Engine is ready\n")
	rt1 := time.Now()

	for start*10 < numSamples {
		logger.Log("Query number %v is pending...\n", start+1)
		t1 := time.Now()
		opt.Start = start * 10
		search, err := stf.googleSearch.Search(ctx, query, &opt)
		if err != nil {
			return err
		}

		var wg sync.WaitGroup
		for index, _ := range search.Items {
			i := *search.Items[index]
			wg.Add(1)

			go func() {
				defer wg.Done()
				b, err := stf.downloader(ctx, i.Link)
				if err != nil {
					//logger.Log("Error downloading image link:%v\n", i.Link)
					stf.writer.errorChan <- err
					return
				}

				kind := getImgType(i.Mime)
				ft, err := stf.imageProcessor.Process(b, kind)
				if err != nil {
					//logger.Log("Error processing image link:%v type:%v\n", i.Link, kind)
					stf.writer.errorChan <- err
					return
				}

				stf.writer.wprotoChan <- ft
			}()
		}

		wg.Wait()
		logger.Info("       > Step %v timestamp it in %v (GOROUTINES: %v)\n", start+1, time.Since(t1), 10)
		start++
	}
	logger.Log("Total %v queries took in total: %v\n\n", start, time.Since(rt1))
	stf.writer.mu.Lock()
	stf.writer.mu.Unlock()
	stf.writer.Close()
	return nil
}
