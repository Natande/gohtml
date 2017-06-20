package gohtml

import (
	"golang.org/x/net/html"
	"log"
)

func AppendChildNodes(parent *html.Node, children []*html.Node) {
	for _, c := range children {
		parent.AppendChild(c)
	}
}

func DetachNodes(nodes []*html.Node) {
	for _, n := range nodes {
		n.Parent.RemoveChild(n)
	}
}

func CompactNode(n *html.Node) {
	var appendNodes []*html.Node
	for c := n.FirstChild; c != nil; {
		CompactNode(c)
		if _mergeTextElements[c.Data] {
			appendNodes = append(appendNodes, GetChildNodes(c)...)
			log.Println("delete", c.Data)
			c = RemoveNode(c)
		} else if c.Type == html.ElementNode && c.FirstChild == nil && !_voidElements[c.Data] {
			log.Println("delete", c.Data)
			c = RemoveNode(c)
		} else {
			c = c.NextSibling
		}
	}

	DetachNodes(appendNodes)
	AppendChildNodes(n, appendNodes)
	if n.FirstChild != nil && n.FirstChild.NextSibling == nil {
		if n.FirstChild.Data == n.Data || (n.FirstChild.Data == "br" && (n.Data == "p" || n.Data == "div")) {
			childNodes := GetChildNodes(n.FirstChild)
			log.Println("delete", n.FirstChild.Data)
			n.RemoveChild(n.FirstChild)
			DetachNodes(childNodes)
			AppendChildNodes(n, childNodes)
		} else if n.FirstChild.Data == "img" && n.Data == "a" {
			*n = *n.FirstChild
		}
	}
}
