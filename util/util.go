package util

import (
	"strings"

	"github.com/cheggaaa/pb"
)

const BadRunes = "?!/;:öä"

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
