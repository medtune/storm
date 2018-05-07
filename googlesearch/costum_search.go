package gcs

import (
	cs "google.golang.org/api/customsearch/v1"
)

type QueryOpt struct {
}

type GoogleSearchEngineService interface {
	Query(query string, opt QueryOpt) (*cs.Search, error)
}

type GCS struct {
	*cs.CseService
}
