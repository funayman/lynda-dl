package util

import "testing"

const (
	KeyCourseId = "courseId"
	KeyVideoId  = "videoId"
)

func TestParseUrl(t *testing.T) {
	tests := []struct {
		url string
		grp map[string]string
	}{
		{"https://www.lynda.com/IT-Infrastructure-tutorials/What-you-should-know-before-watching-course/369186/418860-4.html", map[string]string{KeyVideoId: "418860", KeyCourseId: "369186"}},
		{"https://www.lynda.com/IT-Infrastructure-tutorials/Securing-IoT-Privacy/609023-2.html", map[string]string{KeyCourseId: "609023"}},
	}

	for _, tt := range tests {
		rslt := ParseUrl(tt.url)
		if rslt[KeyVideoId] != tt.grp[KeyVideoId] {
			t.Errorf("video ids do not match! expected: %s; actual %s\n", tt.grp[KeyVideoId], rslt[KeyVideoId])
		}

		if rslt[KeyCourseId] != tt.grp[KeyCourseId] {
			t.Errorf("course ids do not match! expected: %s; actual %s\n", tt.grp[KeyCourseId], rslt[KeyCourseId])
		}
	}
}
