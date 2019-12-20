package crawl

import (
	"io/ioutil"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFindContent(t *testing.T) {

	crawler := NewCrawler()

	tests := []struct {
		path     string
		content  string
		hasError bool
	}{
		{"testdata/simple.html", "Sample text", false},
		{"testdata/complex.html", "<p>Paragraph 1</p><p>Paragraph 2</p>", false},
		{"testdata/empty.html", "", true},
	}
	for _, test := range tests {
		data := loadHTML(test.path)

		content, err := crawler.FindContent(data)

		if test.hasError {
			assert.NotNil(t, err)
			assert.Nil(t, content)
		} else {
			assert.NotNil(t, content)
			assert.Nil(t, err)
			assert.Equal(t, test.content, content)
		}
	}
}

func loadHTML(path string) []byte {
	dat, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return dat
}
