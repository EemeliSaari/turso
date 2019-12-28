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
	Item      *gofeed.Item
	Loaded    bool   `json:"loaded"`
	Erroneous string `json:"errorneous"`
}

// NewArticle ...
func NewArticle(item *gofeed.Item) *Article {
	return &Article{
		Item:      item,
		Loaded:    len(item.Content) > 0,
		Erroneous: "",
	}
}

// FetchHTML ...
func (a Article) fetchHTML() {
	if a.Loaded {
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

	crawler := crawl.NewCrawler()
	content, err := crawler.FindArticleContent(html)
	if err != nil {
		a.Erroneous = err.Error()
	} else {
		a.Item.Content = content
	}
	a.Loaded = true
}

// Hash ...
func (a Article) Hash() [16]byte {
	jsonBytes, _ := json.Marshal(a.Item)
	return md5.Sum(jsonBytes)
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
