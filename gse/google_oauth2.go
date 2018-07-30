package stormtf

import (
	"context"
	"io/ioutil"
	"net/http"

	"golang.org/x/oauth2/google"
)

const (
	GoogleCustomSearchScope = "https://www.googleapis.com/auth/cse"
)

func DefaultGoogleClient(ctx context.Context, scopes ...string) (*http.Client, error) {
	return google.DefaultClient(ctx, scopes...)
}

func GoogleClientFromJSON(ctx context.Context, jfilepath string, scopes ...string) (*http.Client, error) {
	data, err := ioutil.ReadFile(jfilepath)
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(data, scopes...)
	if err != nil {
		return nil, err
	}
	return conf.Client(ctx), nil
}
