package gohtml

import (
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/net/html"
)

const UAiPhoneSafari = "Mozilla/5.0 (iPhone; CPU iPhone OS 10_0_1 like Mac OS X) AppleWebKit/602.1.50 (KHTML, like Gecko) Version/10.0 Mobile/14A403 Safari/602.1"

type Article struct {
	Title   string
	Content string
	Images  []string
}

func ParseArticleLink(link string) (*Article, error) {
	req, err := http.NewRequest("GET", link, nil)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	pageURL, _ := url.Parse(link)

	req.Header.Set("User-Agent", UAiPhoneSafari)
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return ParseArticle(string(data), pageURL)
}

func ParseArticle(htmlStr string, pageURL *url.URL) (*Article, error) {
	doc, err := html.Parse(strings.NewReader(htmlStr))
	if err != nil {
		log.Println(err)
		return nil, err
	}

	if doc == nil {
		log.Println("no html document")
		return nil, errors.New("no document")
	}
	title := GetTitle(doc)
	RemoveAttributes(doc, []string{"src", "href"})
	TidyNodes(doc)
	body := GetNodeByTag(doc, "body")
	if body == nil {
		return nil, errors.New("no <body>")
	}

	nodes := make([]*html.Node, 0, 100)
	nodes = append(nodes, body)
	for i := 0; i < len(nodes); i++ {
		n := nodes[i]
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			if c.Type == html.ElementNode {
				nodes = append(nodes, c)
			}
		}
	}

	var maxTextNode *html.Node
	var maxLen int
	for _, n := range nodes {
		if n.Data == "p" {
			l := GetTextLen(n)
			if l > maxLen {
				maxTextNode = n
				maxLen = l
			}
		}
	}

	if maxTextNode == nil || maxTextNode.Parent == nil {
		return nil, errors.New("invalid article")
	}

	a := &Article{}
	a.Title = title
	a.Content = Render(maxTextNode.Parent)
	a.Images = ParseImages(maxTextNode.Parent, pageURL)
	return a, nil
}

func GetTitle(doc *html.Node) string {
	if doc.Type != html.DocumentNode {
		panic("node type is not DocumentNode")
	}

	head := GetNodeByTag(doc, "head")
	if head == nil {
		return ""
	}

	titleNode := GetNodeByTag(head, "title")
	if titleNode == nil {
		log.Println("no title")
		return ""
	}

	if titleNode.FirstChild != nil && titleNode.FirstChild.Type == html.TextNode {
		return titleNode.FirstChild.Data
	}

	return ""
}

func TidyNodes(n *html.Node) {
	for c := (n.FirstChild); c != nil; {
		switch c.Type {
		case html.ElementNode:
			if _unusedElements[c.Data] {
				c = RemoveNode(c)
				break
			}

			if c.Data == "img" && len(c.Attr) == 0 {
				c = RemoveNode(c)
				break
			} else if c.Data == "a" && len(c.Attr) == 0 {
				c = RemoveNode(c)
				break
			}

			if _voidElements[c.Data] {
				c = c.NextSibling
			} else {
				TidyNodes(c)
				if c.FirstChild == nil {
					c = RemoveNode(c)
				} else {
					c = c.NextSibling
				}
			}
		case html.TextNode:
			c.Data = strings.TrimSpace(c.Data)
			if len(c.Data) == 0 {
				c = RemoveNode(c)
			} else {
				c = c.NextSibling
			}
		default:
			c = RemoveNode(c)
		}
	}
}

func GetTextLen(n *html.Node) int {
	if n.Type == html.TextNode {
		return len(n.Data)
	} else if n.Type == html.ElementNode {
		if n.Data == "img" {
			return 1
		} else {
			k := 0
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				k += GetTextLen(c)
			}

			if n.Data == "a" {
				return 1
			}
			return k
		}
	} else {
		return 0
	}
}

func ParseImages(n *html.Node, pageURL *url.URL) []string {
	if n.Type != html.ElementNode {
		return nil
	}

	var links []string
	if n.Data == "img" {
		for _, a := range n.Attr {
			if a.Key == "src" {
				if u, err := pageURL.Parse(a.Val); err != nil {
					log.Println(err, pageURL.String())
				} else {
					links = append(links, u.String())
				}
			}
		}
	} else {
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			links = append(links, ParseImages(c, pageURL)...)
		}
	}

	return links
}
