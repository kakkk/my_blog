package utils

import (
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/spf13/cast"
	"github.com/writeas/go-strip-markdown"
)

func GetPostPageTitle(title string) string {
	return fmt.Sprintf("%v - kakkk's Blog", title)
}

func GetPostPageDescription(content string) string {
	striped := stripmd.Strip(content)
	striped = strings.Replace(striped, "\n", " ", -1)
	if len(striped) < 97 {
		return striped
	}
	cut := striped[0:96]
	return fmt.Sprintf("%v...", cut)
}

func GetWordCount(content string) string {
	return cast.ToString(utf8.RuneCountInString(stripmd.Strip(content)))
}

func GetPublishAtStr(publishAt *time.Time) string {
	return publishAt.Format("January 02, 2006")
}
