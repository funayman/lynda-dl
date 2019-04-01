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
	"net/url"
	"os"
	"strconv"

	"github.com/funayman/lynda-dl/client"
	"github.com/funayman/lynda-dl/course"
	"github.com/funayman/lynda-dl/util"
	"github.com/spf13/cobra"
)

const (
	CourseRegexp = `https?://(?:www\.)?(?:lynda\.com|educourse\.ga)/(?:(?:[^/]+/){2,3}(?P<courseId>\d+))`
)

var (
	isLearningPath bool
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a Lynda course",
	// Long:  ``,
	PreRun: func(cmd *cobra.Command, args []string) {
		// check if cookie file exists
		if _, err := os.Stat(cookiepath); os.IsNotExist(err) {
			log.Fatalf("cookie file: %s does not exist\n", cookiepath)
		}

		if len(args) == 0 {
			log.Fatal("download command requires a URL")
		}

		// if URLs are passed, ensure they're valid
		for _, argUrl := range args {
			if _, err := url.ParseRequestURI(argUrl); err != nil {
				log.Fatalf("URL %s is malformed", argUrl)
			}
		}

	},
	Run: func(cmd *cobra.Command, args []string) {
		client.Init(cookiepath)

		// check for --learning-path
		if isLearningPath {
			courseIds := util.ExtractCourseIdsFromLearningPath(args[0])
			for _, id := range courseIds {
				c, err := course.Build(id)
				if err != nil {
					log.Fatal(err)
				}

				err = c.Download()
				if err != nil {
					log.Fatal(err)
				}
			}
		}

		params := util.ParseUrl(args[0])
		if params["videoId"] != "" {
			// download video only
		} else if params["courseId"] != "" {
			stringId := params["courseId"]
			courseId, _ := strconv.Atoi(stringId)
			c, err := course.Build(courseId)
			if err != nil {
				log.Fatal(err)
			}

			err = c.Download()
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("All Downloads Complete")
	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().IntVarP(&id, "course-id", "i", 0, "Lynda course id")
	// downloadCmd.MarkFlagRequired("course-id")

	downloadCmd.Flags().StringVarP(&cookiepath, "cookies", "c", "", "path to cookies.txt")
	downloadCmd.MarkFlagRequired("cookies")

	downloadCmd.Flags().BoolVar(&isLearningPath, "learning-path", false, "Url for Learning Path rather than course")
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
