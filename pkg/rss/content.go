package rss

import (
	"crypto/md5"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"runtime"
	"sync"

	"github.com/EemeliSaari/turso/pkg/crawl"
	"github.com/mmcdole/gofeed"
)

// Article ...
type Article struct {
	Item *gofeed.Item

	html   []byte
	loaded bool
}

// NewArticle ...
func newArticle(item *gofeed.Item) *Article {
	return &Article{
		Item: item, html: []byte{},
		loaded: len(item.Content) > 0,
	}
}

// FetchHTML ...
func (a Article) fetchHTML() {
	if a.loaded {
		return
	}

	url := a.Item.Link
	resp, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer resp.Body.Close()

	html, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	a.html = html
	a.loaded = true
}

// Hash ...
func (a Article) Hash() [16]byte {
	jsonBytes, _ := json.Marshal(a.Item)
	return md5.Sum(jsonBytes)
}

func (a Article) crawl() {
	if len(a.html) == 0 {
		a.fetchHTML()
	}
	html := a.html

	crawler := crawl.NewCrawler()
	content, _ := crawler.FindContent(html)
	a.Item.Content = content
}

func loadContent(articles []*Article) {
	var wg sync.WaitGroup

	wg.Add(len(articles))

	for i := 0; i < len(articles); i++ {
		go func(article *Article) {
			defer wg.Done()
			runtime.Gosched()
			article.fetchHTML()
		}(articles[i])
	}

	wg.Wait()
}
