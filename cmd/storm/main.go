// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"fmt"
	olog "log"
	"net/http"
	"os"

	"github.com/medtune/stormtf/cse"
	"github.com/medtune/stormtf/features"
	"github.com/medtune/stormtf/filters"
	"github.com/medtune/stormtf/log"
	"github.com/medtune/stormtf/storm"
	"github.com/spf13/cobra"
)

var (
	VerboseLevel            int
	GoogleCredentialDefault bool
	GoogleCrendetialsPath   string
	Scopes                  []string
	SearchQuery             string
	QueryOpt                = cse.QueryOption{}
	SearchEngineID          string
	ResizeOpt               string
	ImageOutFormat          string
	OutputFile              string
	DataKeyName             string
	DataLabelName           string
	AddImageDimFeature      bool
	AddImageDescFeature     bool
	NumResults              int
	ContextTimeout          int
)

const tfrecordsExt = "tfrecords"

func init() {
	stormtfCmd.Flags().IntVarP(&VerboseLevel, "verbose", "v", 2, "Set verbosity level")
	stormtfCmd.Flags().StringVarP(&OutputFile, "output", "o", "", "Set outfile")
	stormtfCmd.Flags().StringSliceVarP(&Scopes, "scopes", "s", []string{"https://www.googleapis.com/auth/cse"}, "Set cse scopes")
	stormtfCmd.Flags().StringVarP(&QueryOpt.SearchType, "searchtype", "S", "image", "Set search data type")
	stormtfCmd.Flags().StringVarP(&SearchQuery, "query", "q", "", "Set search query")
	stormtfCmd.Flags().StringVarP(&SearchEngineID, "engine-id", "e", "", "Set engine ID")
	stormtfCmd.Flags().StringVarP(&ResizeOpt, "resize-format", "r", "linear:256x256", "Resize image formula")
	stormtfCmd.Flags().StringVarP(&ImageOutFormat, "imout-type", "i", filters.JPEG, "Image type encoding")
	stormtfCmd.Flags().StringVarP(&DataKeyName, "data-key", "K", "image", "Set data label name")
	stormtfCmd.Flags().StringVarP(&DataLabelName, "label", "L", "", "Set data label name")
	stormtfCmd.Flags().BoolVarP(&AddImageDimFeature, "imgdim", "D", true, "add image dimentions to proto")
	stormtfCmd.Flags().IntVarP(&NumResults, "number", "n", 0, "Set storm crawler max results")
	stormtfCmd.Flags().StringVarP(&GoogleCrendetialsPath, "gcreds", "C", "", "Set google credentials files")
	stormtfCmd.Flags().BoolVarP(&GoogleCredentialDefault, "default-gcreds", "c", true, "Default google creds auth mode")
	stormtfCmd.Flags().IntVarP(&ContextTimeout, "context-timeout", "t", 0, "Context timeout")
}

func must(err error) {
	if err != nil {
		olog.Panic(err)
	}
}

var stormtfCmd = &cobra.Command{
	Use:   "stormtf",
	Short: "StormTF #shortdesc",
	Long:  `stormtf #longdesc`,
	Run: func(cmd *cobra.Command, args []string) {
		log.SetLoggingLevel(VerboseLevel)
		ctx := context.Background()
		if SearchQuery == "" {
			log.Error("Must provide search query")
			return
		}
		if SearchEngineID == "" {
			log.Error("Must provide search engine ID")
			return
		}
		if NumResults == 0 {
			log.Error("Must provide wanted number of results")
			return
		}
		if DataKeyName == "" {
			log.Error("Must provide data key name")
			return
		}
		if QueryOpt.SearchType != "image" {
			log.Error("Unsupported search type %v", QueryOpt.SearchType)
			return
		}

		// # TODO Context timeout
		var client *http.Client
		if !GoogleCredentialDefault && GoogleCrendetialsPath == "" {
			log.Error("Must provide google crendentials auth mode")
			return
		}
		if GoogleCredentialDefault {
			c, err := cse.DefaultGoogleClient(ctx, cse.GoogleCustomSearchScope)
			if err != nil {
				log.Error("Coudlnt get client. Scope: %v Error:%v", cse.GoogleCustomSearchScope, err)
				return
			}
			client = c
		} else {
			c, err := cse.GoogleClientFromJSON(ctx, GoogleCrendetialsPath, cse.GoogleCustomSearchScope)
			if err != nil {
				log.Error("Coudlnt get client. Scope: %v Error:%v", cse.GoogleCustomSearchScope, err)
				return
			}
			client = c
		}
		service, err := cse.NewGCS(client)
		if err != nil {
			log.Error("Coudlnt get GCS Service. Error:%v", err)
			return
		}

		service.SetEngineID(SearchEngineID)
		imgProc := filters.NewImgProcs()
		if DataLabelName == "" {
			DataLabelName = SearchQuery
		}
		imgProc.AddFeature("label", features.LabelFeature(DataLabelName))
		if imf, x, y, err := filters.ResizeImageFilterFromString(ResizeOpt); err != nil {
			log.Error("Error filter string format. Error:%v", err)
			return
		} else {
			imgProc.AddFilter(imf)
			if AddImageDimFeature {
				imgProc.AddFeature("height", features.NewInt64ListFeature(int64(x)))
				imgProc.AddFeature("width", features.NewInt64ListFeature(int64(y)))
			}
		}
		imgProc.SetDefaultKey(DataKeyName)
		storm := storm.New(service, imgProc)
		var op string
		if OutputFile == "" {
			op = SearchQuery + "." + tfrecordsExt
		} else {
			op = OutputFile + "." + tfrecordsExt
		}
		OutputFile = op
		err = storm.Storm(
			ctx,
			SearchQuery,
			QueryOpt,
			int64(NumResults),
			OutputFile,
		)
		if err != nil {
			log.Error("error: %v", err)
		}
		return
	},
}

func Execute() {
	if err := stormtfCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func main() {
	Execute()
}
