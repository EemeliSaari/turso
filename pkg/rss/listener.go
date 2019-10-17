package rss

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/EemeliSaari/turso/internal/buffer"
	"github.com/mmcdole/gofeed"
)

// CallbackHandler ...
type CallbackHandler func([]*Article)

// Listener ...
type Listener struct {
	parser    *gofeed.Parser
	ticker    *time.Ticker
	buffer    *buffer.Buffer
	callbacks []CallbackHandler
	running   bool
	stop      chan bool

	Sources  []string
	Interval int
}

// NewListener ...
func NewListener(src []string, interval time.Duration) *Listener {
	fp := gofeed.NewParser()
	fl := Listener{
		parser:  fp,
		Sources: src,
		running: false,
		buffer:  buffer.NewBuffer(),
	}
	fl.stop = make(chan bool)
	fl.ticker = time.NewTicker(interval * time.Second)

	return &fl
}

// Start ...
func (fl Listener) Start() (chan bool, error) {
	if fl.running {
		return fl.stop, errors.New("can't start running")
	}

	go func() {
		for {
			select {
			case <-fl.stop:
				fmt.Println("DONE")
				return
			case <-fl.ticker.C:
				fl.tick()
			}
		}
	}()

	fl.running = true

	return fl.stop, nil
}

// Stop ...
func (fl Listener) Stop() (chan bool, error) {
	if fl.running != false {
		return fl.stop, errors.New("can't stop stopped")
	}

	fl.stop <- false
	fl.running = false

	return fl.stop, nil
}

// AddCallback ...
func (fl *Listener) AddCallback(callback CallbackHandler) {
	fl.callbacks = append(fl.callbacks, callback)
}

func (fl Listener) tick() {

	articles := []*Article{}

	for i := 0; i < len(fl.Sources); i++ {
		url := fl.Sources[i]
		feed, _ := fl.parser.ParseURL(url)

		for idx := 0; idx < len(feed.Items); idx++ {
			item := feed.Items[idx]
			article := newArticle(item)

			isNew := fl.buffer.Add(article)
			if isNew {
				articles = append(articles, article)
			}
		}
	}

	loadContent(articles)

	if len(articles) > 0 {
		fl.callback(articles)
	}
}

func (fl Listener) callback(articles []*Article) {
	var wg sync.WaitGroup
	wg.Add(len(fl.callbacks))

	for i := 0; i < len(fl.callbacks); i++ {
		go func(idx int) {
			defer wg.Done()
			fl.callbacks[idx](articles)
		}(i)
	}
	wg.Wait()
}
