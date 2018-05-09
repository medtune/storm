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

package cmd

import (
	"fmt"
	"log"
	"os"

	"github.com/iallabs/stormtf/stormtf"
	"github.com/spf13/cobra"
)

var (
	Storm *stormtf.StormTF

	VerboseLevel          int    = 1
	Provider              string = "google"
	GoogleCrendetialsPath        = ""
	Scopes                       = []string{}
	DataType              string = "image"
	SearchQuery                  = ""
	QueryOpt                     = stormtf.QueryOption{}
	SearchEngineID               = ""
	ResizeOpt             string = ""
	ProtoFormat           string = "features"
	OutputFile            string = ""
	NumResults                   = 0
)

func must(err error) {
	if err != nil {
		log.Panic(err)
	}
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stormtf",
	Short: "Google.com crawler to generate ready to use tfrecords",
	Long:  `Google.com crawler to generate ready to use tfrecords`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Command start")
		fmt.Println(VerboseLevel, OutputFile, Scopes, GoogleCrendetialsPath)

		/*
			ctx := context.Background()
			client, err := stormtf.GoogleClientFromJSON(ctx, GoogleCrendetialsPath, Scopes...)
			must(err)
			s, err := stormtf.NewGCS(client)
			must(err)
			stormAgent, err := stormtf.New(s)
			must(err)
			err = stormAgent.Storm(
				ctx,
				SearchQuery,
				QueryOpt,
				SearchEngineID,
				NumResults,
				OutputFile,
			)
		*/
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ez.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.

	rootCmd.Flags().IntVarP(&VerboseLevel, "verbose", "v", 1, "Set verbosity level")
	rootCmd.Flags().StringVarP(&OutputFile, "output", "o", "records", "Set outfile")
	rootCmd.Flags().StringVarP(&Provider, "provider", "p", "google", "Set search provider")
	rootCmd.Flags().StringSliceVarP(&Scopes, "scopes", "s", []string{"https://www.googleapis.com/auth/cse"}, "Set cse scopes")
	rootCmd.Flags().StringVarP(&DataType, "datatype", "d", "image", "Set search datatype")
	rootCmd.Flags().StringVarP(&SearchQuery, "query", "q", "", "Set search query")
	rootCmd.Flags().StringVarP(&SearchEngineID, "engine-id", "e", "-", "Set engine ID")
	rootCmd.Flags().StringVarP(&ResizeOpt, "resize-format", "r", "linear:256x256", "Resize image formula")
	rootCmd.Flags().IntVarP(&NumResults, "number", "n", 10, "Set storm crawler max results")
	rootCmd.Flags().StringVarP(&GoogleCrendetialsPath, "credentials", "c", "", "Set google credentials files")
}
