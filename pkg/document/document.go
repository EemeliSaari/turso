package document

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/EemeliSaari/turso/internal/utils"
)

// Part ...
type Part struct {
	ByteStart  int `json:"byte_start"`
	ByteEnd    int `json:"byte_end"`
	TokenStart int `json:"token_start"`
	TokenEnd   int `json:"token_end"`
}

// Paragraph ...
type Paragraph struct {
	Part
}

// Title ...
type Title struct {
	Part
}

type htmlSearch struct {
	*Token
	tag   string
	isEnd bool
}

// Document ...
type Document struct {
	Tokens     []*Token     `json:"tokens"`
	Paragraphs []*Paragraph `json:"paragraphs"`
	Titles     []*Title     `json:"title"`
	Raw        string       `json:"raw"`
	Checksum   string       `json:"checksum"`
}

// New ...
func New(b []byte) (*Document, error) {
	t := NewTokenizer()

	tokens := []*Token{}
	htmlTokens := []*htmlSearch{}
	for token := range t.TokenReader(&b) {
		tokens = append(tokens, token)
		if token.Type == HTMLToken && token.Content[0] == '<' {
			start := 1
			isEnd := false
			end := len(token.Content)
			if token.Content[1] == '/' {
				start++
				end--
				isEnd = true
			}
			htmlTokens = append(htmlTokens, &htmlSearch{
				Token: token,
				tag:   strings.Split(string(token.Content[start:end]), " ")[0],
				isEnd: isEnd})
		}
	}

	var end *htmlSearch
	paragraphs := []*Paragraph{}
	titles := []*Title{}
	for idx, token := range htmlTokens {
		if !token.isEnd {
			switch token.tag {
			case "p", "section", "div":
				end = findEnd(token, idx, &htmlTokens)
				paragraphs = append(paragraphs, &Paragraph{Part: *newPart(token, end)})
			case "h1", "h2", "h3", "h4", "h5", "h6":
				end = findEnd(token, idx, &htmlTokens)
				titles = append(titles, &Title{Part: *newPart(token, end)})
			}
		}
	}
	return &Document{
		Tokens:     tokens,
		Paragraphs: paragraphs,
		Titles:     titles,
		Raw:        string(b),
		Checksum:   fmt.Sprintf("%x", hashTokens(&tokens)),
	}, nil
}

func findEnd(t *htmlSearch, from int, tokens *[]*htmlSearch) *htmlSearch {
	var res *htmlSearch
	for _, token := range (*tokens)[from:] {
		if token.isEnd && token.tag == t.tag {
			res = token
			break
		}
	}
	return res
}

func newPart(start *htmlSearch, end *htmlSearch) *Part {
	return &Part{
		TokenStart: start.Idx,
		TokenEnd:   end.Idx,
		ByteStart:  start.Start,
		ByteEnd:    end.End,
	}
}

func hashTokens(tokens *[]*Token) [16]byte {
	bin, err := json.Marshal(*tokens)
	utils.FailOnError(err, "failed to hash tokens")
	return md5.Sum(bin)
}
