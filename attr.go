package gohtml

import (
	"golang.org/x/net/html"
	"strings"
)

func RemoveAttributes(n *html.Node, keepAttrs []string) {
	attrs := make([]html.Attribute, 0, len(n.Attr))
	dataAttrs := make(map[string]html.Attribute)
	originalAttrs := make(map[string]html.Attribute)
	for _, a := range n.Attr {
		if indexOfString(keepAttrs, a.Key) >= 0 {
			attrs = append(attrs, a)
			originalAttrs[a.Key] = a
		} else if i := strings.Index(a.Key, "data-"); i == 0 {
			a.Key = a.Key[5:]
			if len(a.Key) > 0 && indexOfString(keepAttrs, a.Key) >= 0 {
				dataAttrs[a.Key] = a
			}
		}
	}

	for k, v := range dataAttrs {
		if _, found := originalAttrs[k]; !found {
			attrs = append(attrs, v)
		}
	}

	n.Attr = attrs
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		RemoveAttributes(c, keepAttrs)
	}
}

func RemoveSomeAttributes(n *html.Node, someAttrs []string) {
	attrMap := make(map[string]bool, len(someAttrs))
	for _, attr := range someAttrs {
		attrMap[attr] = true
	}

	for i := len(n.Attr) - 1; i >= 0; i-- {
		a := n.Attr[i]
		flag, _ := attrMap[a.Key]
		if flag {
			if i < len(n.Attr)-1 {
				n.Attr = append(n.Attr[0:i], n.Attr[i+1:]...)
			} else {
				n.Attr = n.Attr[0:i]
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		RemoveSomeAttributes(n, someAttrs)
	}
}

func indexOfString(strs []string, s string) int {
	for i, str := range strs {
		if str == s {
			return i
		}
	}

	return -1
}
