package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/EemeliSaari/turso/pkg/constants"
	"github.com/EemeliSaari/turso/pkg/rss"
)

func myCallback(articles []*rss.Article) {
	fmt.Println(len(articles))
	for _, a := range articles {
		buffer := new(bytes.Buffer)
		encoder := json.NewEncoder(buffer)
		encoder.SetIndent("", "\t")
		encoder.SetEscapeHTML(false)

		err := encoder.Encode(a)
		if err != nil {
			panic(err)
		}
		file, err := os.OpenFile(fmt.Sprintf("data/%x.json", a.Hash()), os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}
		_, err = file.Write(buffer.Bytes())
		if err != nil {
			panic(err)
		}
	}
}

func feedsFromFile(path string) []string {
	file, err := os.Open(path)
	defer file.Close()

	if err != nil {
		panic(err)
	}
	reader := bufio.NewReader(file)

	lines := []string{}
	var line string
	for {
		line, err = reader.ReadString('\n')
		if err != nil {
			break
		} else {
			line = strings.Trim(line, " \n\t\r")
			lines = append(lines, line)
		}
	}
	return lines
}

func main() {
	feeds := feedsFromFile("feeds.txt")
	lst := rss.NewListener(feeds, constants.DefaultInterval)
	lst.AddCallback(myCallback)

	err := lst.Start()
	if err != nil {
		fmt.Println("ERROR")
	}
	time.Sleep(20000 * time.Second)

	err = lst.Stop()
	if err != nil {
		fmt.Println(err)
	}
}
