package main

import (
	"fmt"
	"time"

	"github.com/EemeliSaari/turso/pkg/constants"
	"github.com/EemeliSaari/turso/pkg/rss"
)

func myCallback(articles []*rss.Article) {
	fmt.Println(len(articles))
}

func main() {
	x := []string{
		"https://feeds.yle.fi/uutiset/v1/majorHeadlines/YLE_UUTISET.rss",
		"https://www.hs.fi/rss/tuoreimmat.xml",
	}
	lst := rss.NewListener(x, constants.DefaultInterval)
	lst.AddCallback(myCallback)

	_, err := lst.Start()
	if err != nil {
		fmt.Println("ERROR")
	}
	time.Sleep(20000 * time.Second)

	_, err = lst.Stop()
	if err != nil {
		fmt.Println(err)
	}
}
