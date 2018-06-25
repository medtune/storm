package stormtf

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"strings"
)

type Processor interface {
	AddFilter(interface{})
	AddFeature(string, *Feature) error
	SetEncoding(string) error
	Process(io.ReadCloser, string, map[string]*Feature) (*Features, error)
}

type StormTF struct {
	// TODO
	itemCount     int16
	maxThread     int8
	bucketSize    int8
	maxGoroutines int8

	errorHandler errorHandler

	googleSearch GoogleSearchEngineService
	downloader   func(context.Context, string) (io.ReadCloser, error)
	processor    Processor
	writer       *tfrWriter
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
		logger.Error("/errorChannel: received: %v", e)
	}
}

func Testing(gs GoogleSearchEngineService) *StormTF {
	fs := make(map[string]*Feature)
	fs["label"] = &Feature{
		Kind: &Feature_BytesList{BytesList: &BytesList{
			Value: [][]byte{[]byte("dog")},
		}},
	}
	ip := &imageProcessor{
		defaultFeatures: fs,
		filters:         []imageFilter{ResizeImFilter256x256},
	}
	return &StormTF{
		writer:       &tfrWriter{},
		googleSearch: gs,
		processor:    ip,
		errorHandler: defaultErrorHandler(),
		downloader:   downloadBodyRC,
	}
}

func New(gs GoogleSearchEngineService, proc Processor) *StormTF {
	return &StormTF{
		writer:       &tfrWriter{},
		googleSearch: gs,
		processor:    proc,
		errorHandler: defaultErrorHandler(),
		downloader:   downloadBodyRC,
	}
}

func (stf *StormTF) Storm(ctx context.Context, query string, queryOption QueryOption,
	numSamples int64, destination string) error {
	logger.Debug("Storming started\n")

	if numSamples%10 != 0 {
		return fmt.Errorf("numSamples must be in form 10 * k, got :%v", numSamples)
	}

	err := stf.writer.Init(destination, numSamples, stf.errorHandler)
	if err != nil {
		logger.Debug("writer couldnt make it. RIP")
		return err
	}

	logger.Debug("Writer is ready\n")
	opt := queryOption
	var start int64 = 0
	logger.Debug("Starting operations ... query:%v | count:%v\n", query, numSamples)
	rt1 := time.Now()

	for start*10 < numSamples {
		logger.Log("Operation number %v has started", start+1)
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
					stf.writer.errorChan <- err
					return
				}

				kind := getImgType(i.Mime)
				ft, err := stf.processor.Process(b, kind, nil)
				if err != nil {
					stf.writer.errorChan <- err
					return
				}
				stf.writer.wprotoChan <- ft
			}()
		}

		wg.Wait()
		logger.Info("Step %v timing: %v Goroutines: %v)\n", start+1, time.Since(t1), 10)
		start++
	}

	logger.Log("Total operations %v timing: %v", start, time.Since(rt1))
	stf.writer.mu.Lock()
	stf.writer.mu.Unlock()

	if err := stf.writer.Close(); err == nil {
		logger.Log("Shipped file '%v'", destination)
	} else {
		logger.Log("Error closing file %v", destination)
	}

	return err
}

func (stf *StormTF) GetProcessor() Processor {
	return stf.processor
}
