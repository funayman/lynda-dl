package course

import (
	"encoding/json"
	"fmt"
	"io"
	"strings"
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

func (c *LyndaCourse) BuildReadMe() string {
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

	return buf.String()
}

func Unmarshal(data io.Reader) (c LyndaCourse, err error) {
	err = json.NewDecoder(data).Decode(&c)
	return
}

func UnmarshalVideo(data []byte) (v LyndaVideo, err error) {
	err = json.Unmarshal(data, &v)
	return
}
