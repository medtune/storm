package multiwriter

import (
	"io"
	"os"
	"sync"

	"github.com/medtune/storm/features"
	"github.com/medtune/storm/log"
	"github.com/medtune/storm/tfrecords"
)

type ErrorHandler func(error)

type TfrWriter struct {
	w           io.WriteCloser
	mu          sync.Mutex
	WprotoChan  chan *features.Features
	StopChan    chan struct{}
	ErrorChan   chan error
	ErrStopChan chan struct{}
	DoneChan    chan struct{}
}

func (w *TfrWriter) Lock()   { w.mu.Lock() }
func (w *TfrWriter) Unlock() { w.mu.Unlock() }

func (w *TfrWriter) Init(file_path string, strenght int64, handleError func(error)) error {
	os.Remove(file_path)
	file, err := os.Create(file_path)
	if err != nil {
		return err
	}

	w.w = file
	w.WprotoChan = make(chan *features.Features, strenght)
	w.ErrorChan = make(chan error, strenght)
	w.StopChan = make(chan struct{}, 1)
	w.DoneChan = make(chan struct{}, 2)

	go func() {
		total := 0
		totalb := 0
		for {
			select {
			case pb := <-w.WprotoChan:
				w.mu.Lock()
				in, err := tfrecords.WriteTFRecordExample(w.w, &features.Example{
					Features: pb,
				})
				if err != nil {
					w.ErrorChan <- err
					w.mu.Unlock()
					continue
				}
				total++
				totalb += in
				w.mu.Unlock()

			case <-w.StopChan:
				log.Debug("Received stop signal after %v successful writes (total bytes: %v)\n", total, totalb)
				w.DoneChan <- struct{}{}
				return
			}

		}

	}()

	w.ErrStopChan = make(chan struct{}, 1)
	go func() {
		totalErrs := 0
		for {
			select {
			case err := <-w.ErrorChan:
				handleError(err)
				totalErrs++

			case <-w.ErrStopChan:
				log.Debug("Received stop error handling message. Total errors: %v\n", totalErrs)
				w.DoneChan <- struct{}{}
				return
			}
		}
	}()

	return nil
}

func (w *TfrWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.StopChan <- struct{}{}
	w.ErrStopChan <- struct{}{}
	if err := w.w.Close(); err != nil {
		log.Warn("Can't close W file. Got %v\n", err)
		return err
	}

	close(w.WprotoChan)
	close(w.ErrorChan)
	<-w.DoneChan
	<-w.DoneChan
	close(w.DoneChan)
	return nil
}
