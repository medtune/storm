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
	b, err := ioutil.ReadFile("filexd")
	must(err)
	fmt.Println(len(b))
	count := 0
	for {
		var pb stormtf.Sample
		count++
		if err := proto.Unmarshal(b, &pb); err == io.EOF {
			fmt.Println("end of file")
			return
		} else if err != nil {
			fmt.Println("----", err)
			return
		} else if err == nil {
			image := pb.Features.Feature["image"].GetBytesList().Value[0]
			_b, err := proto.Marshal(&pb)
			must(err)
			size := len(_b)
			fmt.Println(size)
			b = b[size:]
			fmt.Println(len(b), "-----")
			if len(b) == 0 {
				fmt.Println("endddd")
				break
			}
			buf := bytes.NewReader(image)
			img, err := jpeg.Decode(buf)
			must(err)
			img = transform.Resize(img, 128, 128, transform.Linear)
			if err := imgio.Save(strconv.Itoa(count)+"filename.png", img, imgio.PNGEncoder()); err != nil {
				panic(err)
			}
			//fmt.Println(len(pb.Features.Feature))
			fmt.Println(count)
		}
	}
}
