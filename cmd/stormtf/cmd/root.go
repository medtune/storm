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
	"os"

	"github.com/iallabs/stormtf"
	"github.com/spf13/cobra"
)

var (
	Storm *stormtf.StormTF

	Verbose     int    = 1
	Provider    string = "google"
	DataType    string = "image"
	QueryOpt    string = ""
	ResizeOpt   string = ""
	ProtoFormat string = "features"
	OutputFile  string = ""
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "stormtf",
	Short: "Google.com crawler to generate ready to use tfrecords",
	Long:  `Google.com crawler to generate ready to use tfrecords`,
	Run: func(cmd *cobra.Command, args []string) {

		if Verbose > 0 {
			fmt.Println("hey", args)
		}
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
	Storm = stormtf.New()
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ez.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().IntVarP(&Verbose, "verbose", "v", 1, "Set verbosity level")
}
