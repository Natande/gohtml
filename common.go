package gohtml

import (
	"bytes"
	"golang.org/x/net/html"
	"regexp"
	"strings"
)

var _importantAttrs = map[string]bool{
	"alt":  true,
	"src":  true,
	"href": true,
}

var _unusedElements = map[string]bool{
	"script": true,
	"style":  true,
	"meta":   true,
	"link":   true,
}

var _mergeTextElements = map[string]bool{
	"b":      true,
	"i":      true,
	"strong": true,
	"font":   true,
	"span":   true,
}

var _voidElements = map[string]bool{
	"area":    true,
	"base":    true,
	"br":      true,
	"col":     true,
	"command": true,
	"embed":   true,
	"hr":      true,
	"img":     true,
	"input":   true,
	"keygen":  true,
	"link":    true,
	"meta":    true,
	"param":   true,
	"source":  true,
	"track":   true,
	"wbr":     true,
}

func RemoveUTF8BOM(htmlStr string) string {
	return strings.Replace(htmlStr, "\xEF\xBB\xBF", "", -1)
}

func RemoveHTMLTags(htmlStr string, tags []string) string {
	for _, tag := range tags {
		htmlStr = strings.Replace(htmlStr, "<"+tag+">", "", -1)
		htmlStr = strings.Replace(htmlStr, "</"+tag+">", "", -1)
		re := regexp.MustCompile("<" + tag + " [^<]*>")
		strs := re.FindAllString(htmlStr, -1)
		for _, s := range strs {
			htmlStr = strings.Replace(htmlStr, s, "", -1)
		}
	}

	return htmlStr
}

func Render(n *html.Node) string {
	buf := new(bytes.Buffer)
	html.Render(buf, n)
	return buf.String()
}
