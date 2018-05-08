package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"

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
	count := 0
	for {
		var pb stormtf.Sample
		count++
		if err := proto.Unmarshal(b, &pb); err == io.EOF {
			fmt.Println("end of file")
			return
		} else if err == nil {
			fmt.Println(len(pb.Features.Feature))
			fmt.Println(count)
		} else if err != nil {
			fmt.Println("XD")
		}
	}
}
