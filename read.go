package main

import (
	"bytes"
	"fmt"
	"image/jpeg"
	"io"
	"io/ioutil"
	"log"
	"strconv"

	"github.com/anthonynsimon/bild/imgio"
	"github.com/anthonynsimon/bild/transform"
	"github.com/golang/protobuf/proto"
	"github.com/iallabs/stormtf/stormtf"
)

func must(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	b, err := ioutil.ReadFile("records")
	must(err)
	fmt.Println(len(b), err)
	count := 0
	for {
		var pb stormtf.Sample
		if err := proto.Unmarshal(b, &pb); err == io.EOF {
			fmt.Println("end of file")
			return
		} else if err != nil {
			fmt.Println("--------", err, len(b), b[0:10], count)
			break
		} else if err == nil {
			image := pb.Features.
				Feature["image"].
				GetBytesList().Value[0]
			label := pb.Features.
				Feature["label"].
				GetBytesList().Value[0]
			fmt.Println(string(label))
			_b, err2 := proto.Marshal(&pb)
			must(err2)
			size := len(_b)
			fmt.Println("size......", size)
			buf := bytes.NewReader(image)
			img, err2 := jpeg.Decode(buf)
			must(err2)
			img = transform.Resize(img, 512, 512, transform.Linear)
			if err2 := imgio.Save(strconv.Itoa(count)+"filename.png", img, imgio.PNGEncoder()); err != nil {
				panic(err2)
			}
			b = b[:len(b)-size]
			fmt.Println(len(b), "-----")
			if len(b) == 0 {
				fmt.Println("endddd")
				return
			}
			//fmt.Println(len(pb.Features.Feature))
			fmt.Println(count)
		}
		count++

	}
}
