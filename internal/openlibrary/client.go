package openlibrary

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

const baseURL = "https://openlibrary.org"

type SearchResult struct {
	NumFound int       `json:"numFound"`
	Docs     []BookDoc `json:"docs"`
}

type BookDoc struct {
	Title     string   `json:"title"`
	Author    []string `json:"author_name"`
	ISBN      []string `json:"isbn"`
	Key       string   `json:"key"`
	CoverI    int      `json:"cover_i"`
	FirstYear int      `json:"first_publish_year"`
}

type Client struct {
	http      *http.Client
	userAgent string
}

func NewClient(userAgent string) *Client {
	return &Client{http: &http.Client{}, userAgent: userAgent}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", c.userAgent)
	return c.http.Do(req)
}

func (c *Client) Search(query string) (SearchResult, error) {
	u, _ := url.Parse(baseURL + "/search.json")
	u.RawQuery = url.Values{"q": {query}}.Encode()

	req, err := http.NewRequest("GET", u.String(), nil)
	if err != nil {
		return SearchResult{}, fmt.Errorf("openlibrary request failed: %w", err)
	}

	resp, err := c.do(req)
	if err != nil {
		return SearchResult{}, fmt.Errorf("openlibrary request failed: %w", err)
	}
	defer resp.Body.Close()

	var result SearchResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return SearchResult{}, fmt.Errorf("openlibrary decode failed: %w", err)
	}
	return result, nil
}
