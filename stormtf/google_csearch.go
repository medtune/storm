package stormtf

import (
	"context"
	"fmt"
	"net/http"

	cs "google.golang.org/api/customsearch/v1"
)

type QueryOption struct {
	C2coff           string
	Cr               string
	DateRestrict     string
	ExactTerms       string
	ExcludeTerms     string
	FileType         string
	Filter           string
	Gl               string
	Googlehost       string
	HighRange        string
	Hl               string
	Hq               string
	ImgColorType     string
	ImgDominantColor string
	ImgSize          string
	ImgType          string
	LinkSite         string
	LowRange         string
	Lr               string
	Num              int64
	OrTerms          string
	RelatedSite      string
	Rights           string
	Safe             string
	SearchType       string
	SiteSearch       string
	SiteSearchFilter string
	Sort             string
	Start            int64
}

type GoogleSearchEngineService interface {
	SetEngineID(string)
	Search(ctx context.Context, query string, option *QueryOption) (*cs.Search, error)
}

func NewGCS(c *http.Client) (GoogleSearchEngineService, error) {
	s, err := cs.New(c)
	if err != nil {
		return nil, err
	}
	return &GCS{
		s: cs.NewCseService(s),
	}, nil
}

type GCS struct {
	s        *cs.CseService
	engineId string
}

func (g *GCS) SetEngineID(id string) {
	g.engineId = id
}

func (g *GCS) Search(ctx context.Context, query string, option *QueryOption) (*cs.Search, error) {
	if g.engineId == "" {
		return nil, fmt.Errorf("Empty engine ID")
	}
	rs := g.s.List(query)
	rs.Cx(g.engineId)
	if option.C2coff != "" {
		rs.C2coff(option.C2coff)
	}
	if option.Cr != "" {
		rs.Cr(option.Cr)
	}
	if option.SearchType != "" {
		rs.SearchType(option.SearchType)
	}
	if option.Start != 0 {
		rs.Start(option.Start)
	}
	if option.ImgType != "" {
		rs.ImgType(option.ImgType)
	}
	if option.ImgSize != "" {
		rs.ImgSize(option.ImgSize)
	}
	if option.ImgColorType != "" {
		rs.ImgColorType(option.ImgColorType)
	}
	if option.Num != 0 {
		rs.Num(option.Num)
	}
	if option.Rights != "" {
		rs.Rights(option.Rights)
	}
	/*
		rs.DateRestrict(option.DateRestrict)
		rs.ExactTerms(option.ExactTerms)
		rs.ExcludeTerms(option.ExcludeTerms)
		rs.FileType(option.FileType)
		rs.Filter(option.Filter)
		rs.Gl(option.Gl)
		rs.Googlehost(option.Googlehost)
		rs.HighRange(option.HighRange)
		rs.Hl(option.Hl)
		rs.Hq(option.Hq)
		rs.ImgColorType(option.ImgColorType)
		rs.ImgDominantColor(option.ImgDominantColor)
		rs.ImgType(option.ImgType)
		rs.LinkSite(option.LinkSite)
		rs.LowRange(option.LowRange)
		rs.Lr(option.Lr)
		rs.Num(option.Num)
		rs.OrTerms(option.OrTerms)
		rs.RelatedSite(option.RelatedSite)
		rs.Rights(option.Rights)
		rs.Safe(option.Safe)
		rs.SiteSearch(option.SiteSearch)
		rs.SiteSearchFilter(option.SiteSearchFilter)
		rs.Sort(option.Sort)
		rs.Start(option.Start)*/
	return rs.Context(ctx).Do()
}
