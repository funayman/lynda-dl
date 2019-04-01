package util

import (
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/cheggaaa/pb"
)

const (
	BadRunes                    = "?!/;:öä"
	CourseVideoRegexp           = `https?://(?:www\.)?(?:lynda\.com)/(?:(?:[^/]+/){2}(?P<courseId>\d+))(?:/)?(?P<videoId>\d+)?`
	LearningPathCourseIdsRegexp = `data-course-id="(?:(?P<courseId>\d+))"`
)

func CleanText(str string) string {
	for _, r := range BadRunes {
		str = strings.Replace(str, string(r), "", -1)
	}
	return str
}

func NewBar(length int) *pb.ProgressBar {
	bar := pb.New(length)
	bar.ShowTimeLeft = true
	bar.ShowSpeed = true
	bar.SetWidth(80)
	bar.SetMaxWidth(80)
	bar.SetUnits(pb.U_BYTES)
	return bar
}

func ExtractCourseIdsFromLearningPath(url string) []int {
	log.Print("Learning Path! Attempting to parse all course ids")
	uniqueIds := make(map[string]bool)
	courseIds := make([]int, 0)

	r := regexp.MustCompile(LearningPathCourseIdsRegexp)

	log.Printf("Requesting URL: %s\n", url)
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Success! Scraping resp.Body")
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

	log.Println("Found Learning Path IDs. Continuing to download")
	return courseIds
}

func ParseUrl(url string) (paramsMap map[string]int) {
	r := regexp.MustCompile(CourseVideoRegexp)
	match := r.FindStringSubmatch(url)

	paramsMap = make(map[string]int)
	for i, name := range r.SubexpNames() {
		if i > 0 && i <= len(match) {
			idMatch, err := strconv.Atoi(match[i])
			if err != nil {
				continue
			}
			paramsMap[name] = idMatch
		}
	}

	return
}
