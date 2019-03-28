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
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strconv"

	"github.com/funayman/lynda-dl/client"
	"github.com/funayman/lynda-dl/course"
	"github.com/spf13/cobra"
)

const (
	CourseRegexp       = `https?://(?:www\.)?(?:lynda\.com|educourse\.ga)/(?:(?:[^/]+/){2,3}(?P<courseId>\d+))`
	RespCourseIdRegexp = `data-course-id="(?:(?P<courseId>\d+))"`
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

	},
	Run: func(cmd *cobra.Command, args []string) {
		uniqueIds := make(map[string]bool)
		courseIds := make([]int, 0)

		if isLearningPath {
			if len(args) == 0 {
				log.Fatal("URL required for learning path")
			}

			r := regexp.MustCompile(RespCourseIdRegexp)

			resp, err := http.Get(args[0])
			if err != nil {
				log.Fatal(err)
			}

			body, _ := ioutil.ReadAll(resp.Body)

			courseIdMatches := r.FindAllStringSubmatch(string(body), -1)
			for _, cid := range courseIdMatches {
				stringId := cid[1]

				if _, ok := uniqueIds[stringId]; ok {
					continue
				}

				uniqueIds[stringId] = true
				id, _ := strconv.Atoi(stringId)
				courseIds = append(courseIds, id)
			}

			fmt.Println(courseIds)
		} else {
			courseIds = append(courseIds, id)
		}

		client.Init(cookiepath)
		for _, id := range courseIds {

			c, err := course.Build(id)
			if err != nil {
				// fmt.Errorf("[ERROR] Cannot build course '%s'\n\terr := %s\n", c.Title, err)
				log.Fatal(err)
			}

			err = c.Download()
			if err != nil {
				fmt.Println(err)
				// fmt.Errorf("[ERROR] Cannot download course '%s'\n\terr := %s\n", c.Title, err.Error())
				continue
			}
		}
		fmt.Println("COMPLETE")
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

func getParams(regEx, url string) (paramsMap map[string]string) {
	var compRegEx = regexp.MustCompile(regEx)
	match := compRegEx.FindStringSubmatch(url)

	paramsMap = make(map[string]string)
	for i, name := range compRegEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}
	return
}
