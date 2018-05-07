package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/iallabs/stormtf/httputil"
)

var url string = "https://pbs.twimg.com/profile_images/616309728688238592/pBeeJQDQ.png"
var url2 string = "https://www.aha.io/assets/github.7433692cabbfa132f34adb034e7909fa.png"

func test1() {
	b, err := httputil.GetBody(url)
	//fmt.Println(b, err)

	file, err := os.Create("./test.jpg")
	if err != nil {
		log.Fatal(err)
	}
	w, err := file.Write(b)
	file.Close()
	fmt.Println(w, err)
}

func test2() {
	var f = func(r *http.Response, err error) error {
		b := r.Body
		if b == nil {
			return fmt.Errorf("lol")
		}
		defer r.Body.Close()
		bytes, err := ioutil.ReadAll(b)
		fmt.Println("xd", bytes[0])
		return err
	}
	err := httputil.MakeRequest(context.Background(), "GET", url, f)
	fmt.Println(err)
}

func test3() {
	t := time.Now()
	_, err := httputil.DoRequest(context.Background(), "GET", url)
	//fmt.Println(b, err)
	fmt.Println(time.Since(t), err)
}

func main() {
	//test2()
	go test3()
	go test3()
	go test3()
	go test3()

	go test3()
	go test3()
	go test3()
	go test3()

	time.Sleep(1 * 1000000000)
}
