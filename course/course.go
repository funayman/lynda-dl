package course

import (
	"encoding/json"
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
	homedir "github.com/mitchellh/go-homedir"
)

const LyndaCourseUrlFormat = "https://www.lynda.com/ajax/player?courseId=%d&type=course"
const LyndaVideoUrlFormat = "https://www.lynda.com/ajax/player?courseId=%d&videoId=%d&type=video"

type LyndaCourse struct {
	ID               int    `json:"ID"`
	Title            string `json:"Title"`
	Description      string `json:"Description"`
	ShortDescription string `json:"ShortDescription"`
	Chapters         []struct {
		CourseID         int    `json:"CourseId"`
		FullChapterIndex string `json:"FullChapterIndex"`
		Description      string `json:"Description"`
		SortIndex        string `json:"SortIndex"`
		ID               int    `json:"ID"`
		Title            string `json:"Title"`
		Videos           []struct {
			Description string `json:"Description"`
			ID          int    `json:"ID"`
			CourseID    int    `json:"CourseID"`
			Title       string `json:"Title"`
			CourseTitle string `json:"CourseTitle"`
			ChapterID   int    `json:"ChapterID"`
		} `json:"Videos"`
	} `json:"Chapters"`
	Authors []struct {
		Fullname    string `json:"Fullname"`
		FirstName   string `json:"FirstName"`
		LastName    string `json:"LastName"`
		Bio         string `json:"Bio"`
		Biographies struct {
			Main string `json:"1"`
		} `json:"Biographies"`
		Image string `json:"Image"`
	} `json:"Authors"`
}

type LyndaVideo struct {
	Streams struct {
		Main struct {
			Format360  string `json:"360"`
			Format540  string `json:"540"`
			Format720  string `json:"720"`
			Format1080 string `json:"1080"`
		} `json:"0"`
		Backup struct {
			Format360  string `json:"360"`
			Format540  string `json:"540"`
			Format720  string `json:"720"`
			Format1080 string `json:"1080"`
		} `json:"1"`
	} `json:"PrioritizedStreams"`
	VideoIndex   int    `json:"VideoIndex"`
	ChapterIndex int    `json:"ChapterIndex"`
	Description  string `json:"Description"`
	ID           int    `json:"ID"`
	CourseID     int    `json:"CourseID"`
	Title        string `json:"Title"`
	CourseTitle  string `json:"CourseTitle"`
	ChapterID    int    `json:"ChapterID"`
}

func Build(id int) (c LyndaCourse, err error) {
	client := client.GetInstance()

	// get course data
	url := fmt.Sprintf(LyndaCourseUrlFormat, id)
	resp, err := client.Get(url)
	if err != nil {
		return
	}
	return unmarshal(resp.Body)
}

func (c *LyndaCourse) writeReadme() (err error) {
	// write content.md
	content, err := os.Create("CONTENT.md")
	if err != nil {
		return
	}
	defer content.Close()

	var buf strings.Builder
	fmt.Fprintf(&buf, "# %s\n", c.Title)

	if c.Description != "" {
		buf.WriteString(c.Description)
	} else {
		buf.WriteString(c.ShortDescription)
	}
	buf.WriteString("\n\n")

	for _, chapter := range c.Chapters {
		fmt.Fprintf(&buf, "## %s\n%s\n\n", chapter.Title, chapter.Description)
		for i, video := range chapter.Videos {
			fmt.Fprintf(&buf, "### %02d - %s\n%s\n\n", i+1, video.Title, video.Description)
		}
	}

	content.WriteString(buf.String())
	return nil
}

func (c *LyndaCourse) Download() (err error) {
	// move to home directory
	home, err := homedir.Dir()
	if err != nil {
		return err
	}
	os.Chdir(home)

	c.buildFolders()
	c.writeReadme()
	c.downloadCourse()
	return nil
}

func (c *LyndaCourse) downloadCourse() (err error) {
	// download videos
	fmt.Printf("*** Downloading Videos for %s ***\n", c.Title)
	for _, chapter := range c.Chapters {
		err = os.Chdir(chapter.Title)
		if err != nil {
			log.Fatal(err)
		}

		for _, video := range chapter.Videos {
			fmt.Printf("Video: %s\n", video.Title)
			url := fmt.Sprintf(LyndaVideoUrlFormat, video.CourseID, video.ID)

			// get JSON feed for video
			fmt.Println("Grabbing JSON feed...")
			data, err := exec.Command("curl", "-L", url, "-b", cookiepath).Output()
			if err != nil {
				log.Fatal(err)
			}

			// unmarshal cURL output
			v, err := unmarshalVideo(data)
			if err != nil {
				log.Fatal(err)
			}

			if v.Title == "" { // something went wrong
				log.Fatalf("error parsing data; no title\n-> cURL output: %s\n-> url: %s", string(data), url)
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

}

func (c *LyndaCourse) buildFolders() error {
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

}

func unmarshal(data io.Reader) (c LyndaCourse, err error) {
	err = json.NewDecoder(data).Decode(&c)
	return
}

func unmarshalVideo(data []byte) (v LyndaVideo, err error) {
	err = json.Unmarshal(data, &v)
	return
}