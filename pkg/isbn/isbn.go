package isbn

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type Client struct {
	dataFile  string
	errFile   string
	dataCache map[string]string
	errCache  map[string]string
	lck       sync.Mutex
	rateLimit sync.Mutex
	client    *http.Client
}

func New(dataFile, errFile string) (*Client, error) {
	var dataCache, errCache map[string]string

	if dataFile != "" {
		if errFile == "" {
			return nil, fmt.Errorf("isbn: missing err file")
		}
		var err error
		dataCache, err = loadCache(dataFile)
		if err != nil {
			return nil, err
		}
		errCache, err = loadCache(errFile)
		if err != nil {
			return nil, err
		}
	} else {
		log.Println("isbn: warning, cache is disabled")
	}

	return &Client{
		dataFile:  dataFile,
		errFile:   errFile,
		dataCache: dataCache,
		errCache:  errCache,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}, nil
}

func loadCache(file string) (map[string]string, error) {
	cache := map[string]string{}
	if _, err := os.Stat(file); err != nil {
		if _, err := os.Create(file); err != nil {
			return nil, fmt.Errorf("isbn: couldn't create file %s: %w", file, err)
		}
	}
	f, err := os.ReadFile(file)
	if err != nil {
		return nil, fmt.Errorf("isbn: couldn't read file %s: %w", file, err)
	}

	// Parse cache file
	if len(f) > 0 {
		if err := json.Unmarshal(f, &cache); err != nil {
			return nil, fmt.Errorf("isbn: couldn't parse file %s: %w", file, err)
		}
	}
	return cache, nil
}

func saveItem(file string, cache map[string]string, lck *sync.Mutex, k, v string) error {
	// Save to cache
	if cache != nil {
		lck.Lock()
		defer lck.Unlock()
		cache[k] = v
		// Write cache to file
		js, err := json.MarshalIndent(cache, "", "  ")
		if err != nil {
			return fmt.Errorf("isbn: couldn't marshal cache %s: %w", file, err)
		}
		if err := os.WriteFile(file, js, 0644); err != nil {
			return fmt.Errorf("isbn: couldn't write cache %s: %w", file, err)
		}
	}
	return nil
}

func (c *Client) Hyphenate(ctx context.Context, raw string, msgs ...string) (string, error) {
	if c.dataCache != nil {
		if hyphenated, ok := c.dataCache[raw]; ok {
			return hyphenated, nil
		}
		if msg, ok := c.errCache[raw]; ok {
			return "", fmt.Errorf("isbn: cached error %s: %w", msg, ErrNotFound)
		}
	}

	// Rate limit
	c.rateLimit.Lock()
	defer func() {
		go func() {
			time.Sleep(250 * time.Millisecond)
			c.rateLimit.Unlock()
		}()
	}()

	// Hyphenate
	log.Println("isbn: hyphenating", raw)
	isbnHyphenated, err := ttlHyphenate(ctx, raw)
	if err != nil {
		ttlErr := err
		isbnHyphenated, err = gobHyphenate(ctx, raw)
		if err != nil {
			err = fmt.Errorf("isbn: %w, %w", ttlErr, err)
		}
	}

	// Save error to cache
	if errors.Is(err, ErrNotFound) {
		msg := strings.Join(msgs, ",")
		if err := saveItem(c.errFile, c.errCache, &c.lck, raw, msg); err != nil {
			return "", err
		}
	}
	if err != nil {
		return "", err
	}

	// Save value to cache
	if err := saveItem(c.dataFile, c.dataCache, &c.lck, raw, isbnHyphenated); err != nil {
		return "", err
	}

	return isbnHyphenated, nil
}

func ttlHyphenate(ctx context.Context, raw string) (string, error) {
	// Create request
	u := fmt.Sprintf("https://www.todostuslibros.com/busquedas?titulo=&autor=&isbn=%s&editorial=&summary=", raw)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		return "", fmt.Errorf("isbn: couldn't create request: %w", err)
	}

	// Create client that doesn't follow redirects
	client := &http.Client{
		Timeout: 10 * time.Second,
		CheckRedirect: func(*http.Request, []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Send request
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("isbn: couldn't send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusFound {
		body, _ := io.ReadAll(resp.Body)
		text := strings.TrimSpace(string(body))
		if len(text) > 100 {
			text = text[:100] + "..."
		}
		if resp.StatusCode == http.StatusOK {
			return "", ErrNotFound
		}
		return "", fmt.Errorf("isbn: unexpected status code: %d (%s)", resp.StatusCode, text)
	}
	redirect := resp.Header.Get("Location")
	split := strings.Split(redirect, "_")

	// Obtain last part of URL
	isbnHyphenated := split[len(split)-1]
	isbn := strings.ReplaceAll(isbnHyphenated, "-", "")
	if !Valid(isbn) {
		return "", fmt.Errorf("isbn: invalid isbn: %s", isbn)
	}
	return isbnHyphenated, nil

}

var ErrNotFound = errors.New("isbn: not found")

func gobHyphenate(ctx context.Context, raw string) (string, error) {
	// Create a custom transport with insecure skip verification
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		return "", fmt.Errorf("isbn: couldn't create cookie jar: %w", err)
	}

	// Create an HTTP client with the custom transport
	client := &http.Client{
		Transport: tr,
		Timeout:   5 * time.Second,
		Jar:       jar,
	}

	resp, err := client.Get("https://www.culturaydeporte.gob.es/webISBN/tituloSimpleFilter.do?cache=init&prev_layout=busquedaisbn&layout=busquedaisbn&language=es")
	if err != nil {
		return "", fmt.Errorf("isbn: couldn't get: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return "", fmt.Errorf("isbn: unexpected status code: %d", resp.StatusCode)
	}

	// Create request
	u := "https://www.culturaydeporte.gob.es/webISBN/tituloSimpleDispatch.do"
	// Parameters x-www-form-urlencoded
	values := url.Values{}
	values.Set("params.forzaQuery", "N")
	values.Set("params.cdispo", "A")
	values.Set("params.cisbnExt", raw)
	values.Set("params.liConceptosExt[0].texto", "")
	values.Set("params.orderByFormId", "1")
	values.Set("action", "Buscar")
	values.Set("language", "es")
	values.Set("prev_layout", "busquedaisbn")
	values.Set("layout", "busquedaisbn")

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(values.Encode()))
	if err != nil {
		return "", fmt.Errorf("isbn: couldn't create request: %w", err)
	}
	req.Header.Set("Origin", "https://www.culturaydeporte.gob.es")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Send request
	resp, err = client.Do(req)
	if err != nil {
		return "", fmt.Errorf("isbn: couldn't send request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("isbn: unexpected status code: %d (%s)", resp.StatusCode, string(body))
	}

	// Search for div.isbnResultado a on the HTML
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return "", fmt.Errorf("isbn: couldn't parse body: %w", err)
	}
	isbn := doc.Find("div.isbnResultado a").First().Text()
	if isbn == "" {
		return "", fmt.Errorf("isbn: couldn't find isbn")
	}
	return isbn, nil
}

func Valid(code string) bool {
	// Remove hyphens
	code = strings.ReplaceAll(code, "-", "")

	// Calculate the checksum digit
	if len(code) != 13 {
		return false
	}

	checksum := 0
	for i, char := range code[:12] {
		digit, err := strconv.Atoi(string(char))
		if err != nil {
			return false
		}

		if i%2 == 0 {
			checksum += digit
		} else {
			checksum += digit * 3
		}
	}

	checksum = (10 - (checksum % 10)) % 10

	// Compare the calculated checksum with the 13th digit
	return strconv.Itoa(checksum) == string(code[12])
}
