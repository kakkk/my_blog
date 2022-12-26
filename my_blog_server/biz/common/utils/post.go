package utils

import (
	"fmt"
	"strings"
	"time"

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

func GetPostInfo(editor string, publishAt time.Time, content string, pv int64) string {
	// 降级
	if editor == "" {
		editor = "kakkk"
	}
	publishAtStr := publishAt.Format("January 02, 2006")
	readTime := len(content) / 275
	if readTime < 1 {
		readTime = 1
	}
	return fmt.Sprintf("%v · %v min · %v · %v visited", publishAtStr, readTime, editor, pv)
}
