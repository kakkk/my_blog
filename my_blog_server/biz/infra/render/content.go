package render

import "github.com/russross/blackfriday/v2"

func RenderContent(content string) string {
	extensions := blackfriday.CommonExtensions | blackfriday.HardLineBreak
	return string(blackfriday.Run(
		[]byte(content),
		blackfriday.WithExtensions(extensions),
	))
}
