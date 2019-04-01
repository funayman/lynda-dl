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
	"log"
	"net/url"

	"github.com/funayman/lynda-dl/course"
	"github.com/funayman/lynda-dl/util"
	"github.com/spf13/cobra"
)

// listCmd represents the list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "view the contents of a course",
	PreRun: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			log.Fatal("list command requires a URL")
		}

		// if URLs are passed, ensure they're valid
		for _, argUrl := range args {
			if _, err := url.ParseRequestURI(argUrl); err != nil {
				log.Fatalf("URL %s is malformed", argUrl)
			}
		}
	},
	Run: func(cmd *cobra.Command, args []string) {
		params := util.ParseUrl(args[0])
		if courseId, ok := params["courseId"]; ok {
			c, err := course.Build(courseId)
			if err != nil {
				log.Fatal(err)
			}
			c.Print()
		} else {
			log.Fatal("Could not find course id")
		}
	},
}

func init() {
	rootCmd.AddCommand(listCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// listCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// listCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
