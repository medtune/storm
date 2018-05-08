package stormtf

import (
	"fmt"
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
	w.wprotoChan = make(chan *Features, strenght*2)
	w.errorChan = make(chan error, strenght)
	w.stopChan = make(chan struct{}, 1)
	go func() {
		for {
			//fmt.Println("IN - LOOP : /!\\")
			select {
			case pb := <-w.wprotoChan:

				//fmt.Println("GOT ONE !!!")
				bytes, err := proto.Marshal(&Sample{
					Features: pb,
				})
				//fmt.Println("GOT ONE ---")
				if err != nil {
					//fmt.Println("ERROR W", err)
					w.errorChan <- err
				}
				//fmt.Println("GOT ONE ---")
				w.mu.Lock()
				in, err := w.w.Write(bytes)
				w.mu.Unlock()

				if err != nil {
					//fmt.Println("ERROR W 2", err)
					w.errorChan <- err
				}
				//fmt.Println("WROTE ", in, "bytes to file")
				fmt.Println("Write", in, "bytes:", bytes[0:10])

			case <-w.stopChan:
				//fmt.Println("received stop alert 1")
				return

			}

		}
	}()

	w.errStopChan = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case err := <-w.errorChan:
				handleError(err)
			case <-w.errStopChan:
				//fmt.Println("received stop alert 2")
				return
			}
		}
	}()

	return nil
}

func (w *writer) Close() {
	//fmt.Println("hello")
	w.mu.Lock()
	w.stopChan <- struct{}{}
	w.errStopChan <- struct{}{}
	//fmt.Println("lol")
	w.w.Close()
	//fmt.Println("yeah")
	w.mu.Unlock()
	close(w.wprotoChan)
	close(w.errorChan)
	//fmt.Println("------XXXXXXXX----------")
}
