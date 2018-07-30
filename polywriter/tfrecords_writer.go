package stormtf

import (
	"io"
	"os"
	"sync"
)

type errorHandler func(error)

type tfrWriter struct {
	w           io.WriteCloser
	mu          sync.Mutex
	wprotoChan  chan *Features
	stopChan    chan struct{}
	errorChan   chan error
	errStopChan chan struct{}
	doneChan    chan struct{}
}

func (w *tfrWriter) Init(file_path string, strenght int64, handleError func(error)) error {
	os.Remove(file_path)
	file, err := os.Create(file_path)
	if err != nil {
		return err
	}

	w.w = file
	w.wprotoChan = make(chan *Features, strenght)
	w.errorChan = make(chan error, strenght)
	w.stopChan = make(chan struct{}, 1)
	w.doneChan = make(chan struct{}, 2)

	go func() {
		total := 0
		totalb := 0
		for {
			select {
			case pb := <-w.wprotoChan:
				w.mu.Lock()
				in, err := writeTFRecordExample(w.w, &Example{
					Features: pb,
				})
				if err != nil {
					w.errorChan <- err
					w.mu.Unlock()
					continue
				}
				total++
				totalb += in
				w.mu.Unlock()

			case <-w.stopChan:
				logger.Debug("Received stop signal after %v successful writes (total bytes: %v)\n", total, totalb)
				w.doneChan <- struct{}{}
				return
			}

		}

	}()

	w.errStopChan = make(chan struct{}, 1)
	go func() {
		totalErrs := 0
		for {
			select {
			case err := <-w.errorChan:
				handleError(err)
				totalErrs++

			case <-w.errStopChan:
				logger.Debug("Received stop error handling message. Total errors: %v\n", totalErrs)
				w.doneChan <- struct{}{}
				return
			}
		}
	}()

	return nil
}

func (w *tfrWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()

	w.stopChan <- struct{}{}
	w.errStopChan <- struct{}{}
	if err := w.w.Close(); err != nil {
		logger.Warn("Can't close W file. Got %v\n", err)
		return err
	}

	close(w.wprotoChan)
	close(w.errorChan)
	<-w.doneChan
	<-w.doneChan
	close(w.doneChan)
	return nil
}
