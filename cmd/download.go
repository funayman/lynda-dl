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
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
	"github.com/funayman/lynda-dl/client"
	"github.com/funayman/lynda-dl/course"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
)

const (
	badRunes = "?!/;:öä"
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		client := client.New()

		id := "753913"

		// move to home directory
		home, err := homedir.Dir()
		if err != nil {
			log.Fatal(err)
		}
		os.Chdir(home)

		// get course data
		url := fmt.Sprintf(course.LyndaCourseUrlFormat, id)
		resp, err := client.Get(url)
		if err != nil {
			log.Fatal(err)
		}
		c, err := course.Unmarshal(resp.Body)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println(c.BuildReadMe())
		os.Exit(0)

		// build folders
		folder := c.Title
		for _, r := range badRunes {
			folder = strings.Replace(folder, string(r), "", -1)
		}
		err = os.Mkdir(folder, os.ModePerm)
		if err != nil {
			log.Fatal(err)
		}
		os.Chdir(folder)
		for _, chapter := range c.Chapters {
			folder := chapter.Title
			for _, r := range badRunes {
				folder = strings.Replace(folder, string(r), "", -1)
			}
			err = os.Mkdir(folder, os.ModePerm)
			if err != nil {
				log.Fatal(err)
			}
			chapter.Title = folder // replace incase of any manipluation
		}

		// download videos
		fmt.Printf("Downloading Videos for %s\n", c.Title)
		for _, chapter := range c.Chapters {
			err = os.Chdir(chapter.Title)
			if err != nil {
				log.Fatal(err)
			}

			for _, video := range chapter.Videos {
				fmt.Printf("Video: %s\n", video.Title)
				url := fmt.Sprintf(course.LyndaVideoUrlFormat, video.CourseID, video.ID)

				// get JSON feed for video
				fmt.Println("Grabbing JSON feed...")
				data, err := exec.Command("curl", "-L", url, "-b", "/Volumes/Macintosh HD/Users/funayman/Downloads/cookies.txt").Output()
				if err != nil {
					log.Fatal(err)
				}

				// unmarshal cURL output
				v, err := course.UnmarshalVideo(data)
				if err != nil {
					log.Fatal(err)
				}

				if v.Title == "" {
					// something went wrong
					fmt.Println(url)
					log.Fatal(errors.New("error parsing data; no title; cURL output: " + string(data)))
				}

				var videoUrl string
				if v.Streams.Main.Format1080 != "" {
					videoUrl = v.Streams.Main.Format1080
				} else if v.Streams.Main.Format720 != "" {
					videoUrl = v.Streams.Main.Format720
				} else if v.Streams.Main.Format540 != "" {
					videoUrl = v.Streams.Main.Format540
				} else if v.Streams.Main.Format360 != "" {
					videoUrl = v.Streams.Main.Format360
				}

				if videoUrl == "" {
					log.Fatal(errors.New("no available videos"))
				}

				fmt.Printf("Downloading %s...\n", v.Title)

				// ready to go do the actual downloading
				resp, err := client.Get(videoUrl)
				if err != nil {
					log.Fatal(err)
				}

				contentLength := resp.Header.Get("Content-Length")
				length, _ := strconv.Atoi(contentLength)

				bar := pb.New(length).SetUnits(pb.U_BYTES)
				bar.Start()

				fileName := fmt.Sprintf("%02d - %s.mp4", v.VideoIndex, v.Title)
				f, err := os.Create(fileName)
				if err != nil {
					log.Fatal(err)
				}

				reader := bar.NewProxyReader(resp.Body)
				io.Copy(f, reader)

				bar.Finish()
			}

			err = os.Chdir("..")
			if err != nil {
				log.Fatal(err)
			}
		}

		fmt.Println("COMPLETE")

		// out, err := exec.Command("curl", "-L", "https://www.lynda.com/ajax/player?courseId=753913&videoId=775920&type=video", "-b", "../../../../../../Downloads/cookies.txt").Output()
		// if err != nil {
		// 	log.Fatal(err)
		// }
		// fmt.Println(string(out))

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
