package couch

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// http://docs.couchdb.org/en/stable/intro/api.html

// Client ...
type Client struct {
	*http.Client
	config Config
}

// Result ...
type Result struct {
	Ok     bool   `json:"ok"`
	Error  string `json:"error"`
	Reason string `json:"reason"`
}

// Document ...
type Document struct {
	ID  string `json:"_id"`
	Rev string `json:"_rev"`
}

// Row ...
type Row struct {
	ID  string                 `json:"id"`
	Key string                 `json:"key"`
	Doc map[string]interface{} `json:"doc"`
}

// DocumentList ...
type DocumentList struct {
	TotalRows int   `json:"total_rows"`
	Offset    int   `json:"offset"`
	Rows      []Row `json:"rows"`
}

// New ...
func New(c Config) (*Client, error) {
	return &Client{
		config: c,
		Client: &http.Client{},
	}, nil
}

// CreateDatabase ...
func (c Client) CreateDatabase() (*Result, error) {
	var result *Result
	url := c.config.asURL()

	body, err := c.factory(url, "PUT", []byte{})
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(body, &result)

	return result, nil
}

// ListDocuments ...
func (c Client) ListDocuments(includeDoc bool) (*DocumentList, error) {
	var docs DocumentList

	url := c.config.asURL() + "_all_docs?"

	if includeDoc {
		url += fmt.Sprintf("include_docs=true&")
	}

	body, err := c.factory(url, "GET", []byte{})
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(body, &docs)
	return &docs, nil
}

// Insert ...
func (c Client) Insert(doc interface{}) (*Result, error) {
	var result *Result
	js, _ := json.Marshal(doc)

	url := c.config.asURL() + c.generateID(js)

	body, err := c.factory(url, "PUT", js)
	if err != nil {
		return nil, err
	}

	_ = json.Unmarshal(body, &result)

	if !result.Ok {
		return result, errors.New(result.Reason)
	}
	return result, err
}

func (c Client) setAuth(req *http.Request) error {
	req.SetBasicAuth(c.config.Username, c.config.Password)
	return nil
}

func (c Client) factory(url string, method string, content []byte) ([]byte, error) {
	req, err := http.NewRequest(method, url, bytes.NewBuffer(content))
	if err != nil {
		return nil, err
	}

	if err := c.setAuth(req); err != nil {
		return nil, err
	}

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	return body, err
}

func (c Client) generateID(data []byte) string {
	return fmt.Sprintf("%x", md5.Sum(data))
}
