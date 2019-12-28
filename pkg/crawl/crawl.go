package crawl

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

// Crawler ...
type Crawler struct {
}

// NewCrawler ...
func NewCrawler() *Crawler {
	c := Crawler{}
	return &c
}

// FindContent ...
func (c Crawler) FindArticleContent(data []byte) (string, error) {
	if len(data) == 0 {
		return "", errors.New("Empty data received.")
	}

	raw := bytes.NewReader(data)
	doc, err := html.Parse(raw)
	if err != nil {
		return "", err
	}

	var articleBody *html.Node
	var runner func(*html.Node, bool)
	runner = func(n *html.Node, useArticle bool) {
		if useArticle {
			potential, err := search(n, "article")
			if err != nil {
				runner(n, false)
			} else {
				n = potential
			}
		}

		body, err := searchBody(n)

		if err != nil && useArticle {
			if n.Parent == nil {
				return
			} else {
				n.Parent.RemoveChild(n)
				runner(doc, useArticle)
			}
		} else if err != nil && !useArticle {
			return
		} else {
			articleBody = body
		}
	}
	runner(doc, true)

	if articleBody == nil {
		return "", errors.New("Could not find valid articles.")
	}

	// Prune the HTML
	for _, t := range []string{"script", "figure", "img"} {
		remove(articleBody, t)
	}

	return renderNode(articleBody), nil
}

func search(node *html.Node, target string) (*html.Node, error) {
	var res *html.Node
	var f func(*html.Node, string)
	f = func(n *html.Node, t string) {
		if n.Type == html.ElementNode && n.Data == t {
			res = n
			return
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child, t)
		}
	}
	f(node, target)

	if res == nil {
		return nil, fmt.Errorf("Could not find %s", target)
	}

	return res, nil
}

func searchBody(node *html.Node) (*html.Node, error) {
	// Refactor this as configurable
	pattern := "(([aA]rticle)|([pP]ost)){1}[-_]*(([bB]ody)|([sS]ection)|([cC]ontent)|([cC]ontainer)){1}\b"

	var res *html.Node
	var f func(*html.Node)
	f = func(n *html.Node) {
		field := attributeString(n)
		match, err := regexp.MatchString(pattern, field)
		if err != nil {
			panic(err)
		}
		if match {
			res = n
			return
		}

		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}
	f(node)

	if res == nil {
		return nil, errors.New("Could not find article body.")
	}

	return res, nil
}

func remove(node *html.Node, target string) {
	if node.Type == html.ErrorNode {
		return
	} else if node.Type == html.ElementNode && node.Data == target {
		node.Parent.RemoveChild(node)
	} else {
		return
	}
	for c := node.FirstChild; c != nil; c = c.NextSibling {
		remove(c, target)
	}
}

func attributeString(node *html.Node) string {
	text := ""
	for _, a := range node.Attr {
		text += fmt.Sprintf(" %s=%s", a.Key, a.Val)
	}
	return text
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)

	text := buf.String()
	for _, t := range []string{"\n", "\t", "\b"} {
		text = strings.ReplaceAll(text, t, "")
	}
	return strings.Join(strings.Fields(text), " ")
}
