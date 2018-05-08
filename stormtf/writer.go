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
	w.wprotoChan = make(chan *Features, strenght+1)
	w.errorChan = make(chan error, strenght+1)
	w.stopChan = make(chan struct{}, 1)
	go func() {
		total := 0
		for {
			//fmt.Println("IN - LOOP : /!\\")
			select {
			case <-w.stopChan:
				//fmt.Println("received stop alert 1")
				return
			case pb := <-w.wprotoChan:
				//fmt.Println("GOT ONE !!!")
				fmt.Println("GOTCHA", total)
				w.mu.Lock()
				bytes, err := proto.Marshal(&Sample{
					Features: pb,
				})
				fmt.Println("UNMARSHELED", total)
				//fmt.Println("GOT ONE ---")
				if err != nil {
					//fmt.Println("ERROR W", err)
					w.errorChan <- err
				}
				//fmt.Println("GOT ONE ---")

				in, err := w.w.Write(bytes)
				fmt.Println("JUST WROTE LOL", total)
				if err != nil {
					//fmt.Println("ERROR W 2", err)
					w.errorChan <- err
				}
				//fmt.Println("WROTE ", in, "bytes to file")
				fmt.Println("Write", in, "bytes:", bytes[0:10])
				w.errorChan <- fmt.Errorf("take this newbie %v", total)
				total += 1

				w.mu.Unlock()

			}

		}
		fmt.Println("Wrote total ", total, "bytes")

	}()

	w.errStopChan = make(chan struct{}, 1)
	go func() {
		for {
			select {
			case err := <-w.errorChan:
				handleError(err)
			case <-w.errStopChan:
				fmt.Println("received stop alert 2")
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
	err := w.w.Close()
	if err != nil {
		fmt.Println(err)
	}
	//fmt.Println("yeah")
	w.mu.Unlock()
	close(w.wprotoChan)
	close(w.errorChan)
	fmt.Println("------XXXXXXXX----------")
}
