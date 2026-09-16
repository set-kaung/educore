package openlibrary

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
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
	return &Client{
		http:      &http.Client{Timeout: 10 * time.Second},
		userAgent: userAgent,
	}
}

func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", c.userAgent)
	return c.http.Do(req)
}

func (c *Client) Search(query string) (SearchResult, error) {
	u, _ := url.Parse(baseURL + "/search.json")
	u.RawQuery = url.Values{"q": {query}}.Encode()

	var lastErr error
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequest("GET", u.String(), nil)
		if err != nil {
			return SearchResult{}, fmt.Errorf("openlibrary request failed: %w", err)
		}

		resp, err := c.do(req)
		if err != nil {
			lastErr = fmt.Errorf("openlibrary request failed: %w", err)
			if attempt == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return SearchResult{}, lastErr
		}

		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			lastErr = fmt.Errorf("openlibrary read body failed: %w", err)
			if attempt == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return SearchResult{}, lastErr
		}

		if resp.StatusCode != http.StatusOK {
			lastErr = fmt.Errorf("openlibrary returned status %d: %s", resp.StatusCode, string(body))
			if attempt == 0 {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			return SearchResult{}, lastErr
		}

		var result SearchResult
		if err := json.Unmarshal(body, &result); err != nil {
			return SearchResult{}, fmt.Errorf("openlibrary decode failed: %w", err)
		}
		return result, nil
	}
	return SearchResult{}, lastErr
}
