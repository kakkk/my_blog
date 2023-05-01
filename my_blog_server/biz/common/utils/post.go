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
	return cutContentByCount(content, 97)
}

func GetPostMetaAbstract(content string) string {
	return cutContentByCount(content, 127)
}

func cutContentByCount(content string, count int) string {
	striped := stripmd.Strip(content)
	striped = strings.Replace(striped, "\n", " ", -1)
	runeArr := []rune(striped)
	if len(runeArr) < count {
		return string(runeArr)
	}
	cut := runeArr[0:count]
	return fmt.Sprintf("%v...", string(cut))
}

func GetWordCount(content string) string {
	return cast.ToString(utf8.RuneCountInString(stripmd.Strip(content)))
}

func GetPublishAtStr(publishAt *time.Time) string {
	return publishAt.Format("January 02, 2006")
}

func GetPostInfo(editor string, publishAt time.Time, content string) string {
	publishAtStr := publishAt.Format("January 02, 2006")
	readTime := utf8.RuneCountInString(stripmd.Strip(content)) / 275
	if readTime < 1 {
		readTime = 1
	}
	return fmt.Sprintf("%v · %v min · %v", publishAtStr, readTime, editor)
}
