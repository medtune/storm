package main

import (
	"context"
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
	//client, err := defaultGoogleClient(ctx, "https://www.googleapis.com/auth/cse")
	cf := "C:/Users/User13/Downloads/Project Tetra - Tests-4b94221e79fb.json"
	scope1 := "https://www.googleapis.com/auth/cse"

	client, err := gs.GoogleClientFromJSON(ctx, cf, scope1)

	//google.DefaultTokenSource()
	/*
		ctx := context.Background()
		data, err := ioutil.ReadFile("C:/Users/User13/Desktop/client_secret_1040422134366-845gdntao0okonamh4l8l0mq3ep2idia.apps.googleusercontent.com.json")
		if err != nil {
			log.Fatal(err)
		}

		conf, err := google.CredentialsFromJSON(ctx, data, "https://www.googleapis.com/auth/cse")
		if err != nil {
			log.Fatal(err)
		}
		// Initiate an http.Client. The following GET request will be
		// authorized and authenticated on the behalf of
		// your service account.
		client := conf.Client(oauth2.NoContext)

		/*
			tk := "AIzaSyC_lhX5ngPtzmjOZXA-EGxO_kGan5PQE3Q"
			ctx := context.WithValue(context.Background(), oauth2.HTTPClient, &http.Client{
				Transport: &transport.APIKey{Key: tk},
			})

			conf := &oauth2.Config{
				ClientID:     "1040422134366-845gdntao0okonamh4l8l0mq3ep2idia.apps.googleusercontent.com",
				ClientSecret: "Q0ZMdj9nOVJyKa6VLkjqfYZa",
				RedirectURL:  "localhost",
				Scopes:       []string{"https://www.googleapis.com/auth/cse"},
				Endpoint:     google.Endpoint,
			}
			url := conf.AuthCodeURL("state", oauth2.AccessTypeOffline)
			fmt.Printf("Visit the URL for the auth dialog: %v", url)
	*/

	// Use the authorization code that is pushed to the redirect
	// URL. Exchange will do the handshake to retrieve the
	// initial access token. The HTTP Client returned by
	// conf.Client will refresh the token as necessary.
	must(err)
	s, err := gs.NewGCS(client)
	must(err)
	key := "009324899209435307429:zqyjhz9e1ki"
	storm, _ := gs.New(s)
	err = storm.Storm(ctx, "dog profile", gs.QueryOption{SearchType: "image", Start: 0}, key, 40, "records")
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
