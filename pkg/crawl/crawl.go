package crawl

import (
	"bytes"
	"fmt"
	"io"

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
func (c Crawler) FindContent(data []byte) (string, error) {
	raw := bytes.NewReader(data)
	doc, _ := html.Parse(raw)
	fmt.Println(doc.Type)
	fmt.Println(c)
	articles := c.articles(doc)

	for i, n := range articles {
		text = renderNode(n)

		fmt.Println(n.Attr)
		//fmt.Println(n.Token().Text())
		fmt.Println()
	}
	ret := ""
	return ret, nil
}

func (c Crawler) articles(node *html.Node) []*html.Node {
	results := []*html.Node{}

	var search func(*html.Node)
	search = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "article" {
			results = append(results, n)
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			search(child)
		}
	}

	search(node)

	return results
}

func (c Crawler) attribute(node *html.Node, attr string) (string, error) {
	for _, a := range node.Attr {
		if a.Key == attr {
			return a.Val, nil
		}
	}
	return "", fmt.Errorf("Could not find %s", attr)
}

func findArticleBody(nodes *[]html.Node) *html.Node {
	if len(nodes) == 1 {
		return nodes[0]
	}

	text := ""
	longest := text
	best := 0
	for i, n := range nodes {
		text = renderNode(n)
		if node_length := len(text); node_length >= best {
			longest = text
		}
		fmt.Println(n.Attr)
		//fmt.Println(n.Token().Text())
		fmt.Println()
	}
}

func texts(node *html.Node) {
	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			fmt.Println(n.Data)
		}
		for child := n.FirstChild; child != nil; child = child.NextSibling {
			f(child)
		}
	}

	f(node)
}

func renderNode(n *html.Node) string {
	var buf bytes.Buffer
	w := io.Writer(&buf)
	html.Render(w, n)
	return buf.String()
}
