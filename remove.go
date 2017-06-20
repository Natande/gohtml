package gohtml

import (
	"golang.org/x/net/html"
)

func RemoveNodesBySelector(n *html.Node, selector string) {
	nodes := GetNodesBySelector(n, selector)
	for _, nod := range nodes {
		if nod.Parent != nil {
			nod.Parent.RemoveChild(nod)
		}
	}
}

func RemoveNode(child *html.Node) (nextSibling *html.Node) {
	ns := child.NextSibling
	child.Parent.RemoveChild(child)
	return ns
}
