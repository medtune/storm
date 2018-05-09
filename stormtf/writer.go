package stormtf

import (
	"io"
	"os"
	"sync"

	"github.com/golang/protobuf/proto"
)

type writer struct {
	w           io.WriteCloser
	mu          sync.Mutex
	wprotoChan  chan *Features
	stopChan    chan struct{}
	errorChan   chan error
	errStopChan chan struct{}
}

func (w *writer) Init(file_path string, strenght int64, handleError func(error)) error {
	os.Remove(file_path)
	file, err := os.Create(file_path)
	if err != nil {
		return err
	}
	//fmt.Println(err, file)
	w.w = file
	w.wprotoChan = make(chan *Features, strenght)
	w.errorChan = make(chan error, strenght)
	w.stopChan = make(chan struct{}, 1)
	go func() {
		total := 0
		totalb := 0
		for {
			//fmt.Println("IN - LOOP : /!\\")
			select {
			case <-w.stopChan:
				logger.Info("    (+++++) Received stop signal after %v successful writes (total bytes: %v)\n", total, totalb)
				return
			case pb := <-w.wprotoChan:
				w.mu.Lock()
				bytes, err := proto.Marshal(&Sample{
					Features: pb,
				})
				if err != nil {
					w.errorChan <- err
				}

				in, err := w.w.Write(bytes)
				if err != nil {
					w.errorChan <- err
				}
				total++
				totalb += in
				w.mu.Unlock()

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
				logger.Info("    (+++++) Received stop error handling message. Total errors: %v\n", totalErrs)
				return
			}
		}
	}()

	return nil
}

func (w *writer) Close() {
	w.mu.Lock()
	w.stopChan <- struct{}{}
	w.errStopChan <- struct{}{}
	err := w.w.Close()
	if err != nil {
		logger.Warn("Can't close W file. Got %v", err)
	}
	w.mu.Unlock()
	close(w.wprotoChan)
	close(w.errorChan)
}
