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

	"github.com/funayman/lynda-dl/course"
	"github.com/spf13/cobra"
)

const (
	badRunes = "?!/;:öä"
)

var (
	id         int
	cookiepath string
)

// downloadCmd represents the download command
var downloadCmd = &cobra.Command{
	Use:   "download",
	Short: "Download a Lynda course",
	// Long:  ``,
	Run: func(cmd *cobra.Command, args []string) {
		// check if cookie file exists
		if _, err := os.Stat(cookiepath); os.IsNotExist(err) {
			log.Fatalf("cookie file: %s does not exist\n", cookiepath)
		}

		c := course.Build(id)
		c.Download()

		fmt.Println("COMPLETE")

	},
}

func init() {
	rootCmd.AddCommand(downloadCmd)

	downloadCmd.Flags().IntVarP(&id, "course-id", "i", 0, "Lynda course id")
	downloadCmd.MarkFlagRequired("course-id")

	downloadCmd.Flags().StringVarP(&cookiepath, "cookies", "c", "", "path to cookies.txt")
	downloadCmd.MarkFlagRequired("cookies")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// downloadCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// downloadCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
