package storm

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"strings"

	"github.com/medtune/stormtf/cse"
	"github.com/medtune/stormtf/features"
	"github.com/medtune/stormtf/filters"
	"github.com/medtune/stormtf/httputil"
	"github.com/medtune/stormtf/log"
	"github.com/medtune/stormtf/multiwriter"
)

type Processor interface {
	AddFilter(interface{})
	AddFeature(string, *features.Feature) error
	SetEncoding(string) error
	Process(io.ReadCloser, string, map[string]*features.Feature) (*features.Features, error)
}

type StormTF struct {
	itemCount     int16
	maxThread     int8
	bucketSize    int8
	maxGoroutines int8

	errorHandler multiwriter.ErrorHandler

	googleSearch cse.GoogleSearchEngineService
	downloader   func(context.Context, string) (io.ReadCloser, error)
	processor    Processor
	writer       *multiwriter.TfrWriter
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

func defaultErrorHandler() func(e error) {
	return func(e error) {
		log.Error("/errorChannel: received: %v", e)
	}
}

func Testing(gs cse.GoogleSearchEngineService) *StormTF {
	fs := make(map[string]*features.Feature)
	fs["label"] = &features.Feature{
		Kind: &features.Feature_BytesList{BytesList: &features.BytesList{
			Value: [][]byte{[]byte("dog")},
		}},
	}

	ip := &filters.ImageProcessor{
		DefaultFeatures: fs,
		Filters:         []filters.ImageFilter{filters.ResizeImFilter256x256},
	}

	return &StormTF{
		writer:       &multiwriter.TfrWriter{},
		googleSearch: gs,
		processor:    ip,
		errorHandler: defaultErrorHandler(),
		downloader:   httputil.DownloadBodyRC,
	}
}

func New(gs cse.GoogleSearchEngineService, proc Processor) *StormTF {
	return &StormTF{
		writer:       &multiwriter.TfrWriter{},
		googleSearch: gs,
		processor:    proc,
		errorHandler: defaultErrorHandler(),
		downloader:   httputil.DownloadBodyRC,
	}
}

func (stf *StormTF) Storm(ctx context.Context, query string, queryOption cse.QueryOption,
	numSamples int64, destination string) error {
	log.Debug("Storming started\n")

	if numSamples%10 != 0 {
		return fmt.Errorf("numSamples must be in form 10 * k, got :%v", numSamples)
	}

	err := stf.writer.Init(destination, numSamples, stf.errorHandler)
	if err != nil {
		log.Debug("writer couldnt make it. RIP")
		return err
	}

	log.Debug("Writer is ready\n")
	opt := queryOption
	var start int64 = 0
	log.Debug("Starting operations ... query:%v | count:%v\n", query, numSamples)
	rt1 := time.Now()

	for start*10 < numSamples {
		log.Log("Operation number %v has started", start+1)
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
					stf.writer.ErrorChan <- err
					return
				}

				kind := getImgType(i.Mime)
				ft, err := stf.processor.Process(b, kind, nil)
				if err != nil {
					stf.writer.ErrorChan <- err
					return
				}
				stf.writer.WprotoChan <- ft
			}()
		}

		wg.Wait()
		log.Info("Step %v timing: %v Goroutines: %v)\n", start+1, time.Since(t1), 10)
		start++
	}

	log.Log("Total operations %v timing: %v", start, time.Since(rt1))
	stf.writer.Lock()
	stf.writer.Unlock()

	if err := stf.writer.Close(); err == nil {
		log.Log("Shipped file '%v'", destination)
	} else {
		log.Log("Error closing file %v", destination)
	}

	return err
}

func (stf *StormTF) GetProcessor() Processor {
	return stf.processor
}
