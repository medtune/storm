package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	gs "github.com/iallabs/stormtf/stormtf"
	"golang.org/x/oauth2/google"
)

func defaultGoogleClient(ctx context.Context, scopes ...string) (*http.Client, error) {
	//return google.AppEngineTokenSource
	return google.DefaultClient(ctx, scopes...)
}

func must(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func main() {
	ctx := context.Background()
	client, err := defaultGoogleClient(ctx, "https://www.googleapis.com/auth/cse")
	must(err)
	s, err := gs.NewGCS(client)
	must(err)
	key := "009324899209435307429:zqyjhz9e1ki"
	storm, _ := gs.New(s)
	_ = storm.Storm(ctx, "dog profile", gs.QueryOption{SearchType: "image", Start: 11}, key, 10, "records")
	fmt.Println("finished...")
	/*
		r, err := s.Search(ctx, "boein 778", &gs.QueryOption{
			SearchType: "image",
		})
		must(err)
		fmt.Println(r.Queries["nextPage"])
		fmt.Println(r.Items)
		for _, i := range r.Items {
			fmt.Println(i.Mime)
		}
	*/
}
