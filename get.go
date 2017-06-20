package gohtml

import (
	"golang.org/x/net/html"
	"strings"
)

func matchSelector(n *html.Node, selector string) bool {
	var t int
	var val string
	switch selector[0] {
	case '.':
		t = 1
		val = selector[1:]
	case '#':
		t = 2
		val = selector[1:]
	default:
		t = 0
	}

	if n.Type != html.ElementNode {
		return false
	}

	if t == 0 {
		if n.Data == selector {
			return true
		}
	} else {
		for _, a := range n.Attr {
			switch t {
			case 1:
				if a.Key == "class" && strings.Index(a.Val, val) >= 0 {
					return true
				}
			case 2:
				if a.Key == "id" && a.Val == val {
					return true
				}
			}
		}
	}

	return false
}

func GetNodeBySelector(n *html.Node, selector string) *html.Node {
	if matchSelector(n, selector) {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if nn := GetNodeBySelector(c, selector); nn != nil {
			return nn
		}
	}

	return nil
}

func GetNodesBySelector(n *html.Node, selector string) (nodes []*html.Node) {
	if matchSelector(n, selector) {
		nodes = append(nodes, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, GetNodesBySelector(c, selector)...)
	}

	return
}

func GetNodeByTag(n *html.Node, tag string) *html.Node {
	if n.Type == html.ElementNode && n.Data == tag {
		return n
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if nod := GetNodeByTag(c, tag); nod != nil {
			return nod
		}
	}

	return nil
}

func GetNodesByTag(n *html.Node, tag string) (nodes []*html.Node) {
	if n.Type == html.ElementNode && n.Data == tag {
		nodes = append(nodes, n)
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, GetNodesByTag(c, tag)...)
	}

	return
}

func GetChildNodes(n *html.Node) []*html.Node {
	var nodes []*html.Node
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		nodes = append(nodes, c)
	}
	return nodes
}
