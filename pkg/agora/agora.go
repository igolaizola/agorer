package agora

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type client struct {
	host   string
	token  string
	client *http.Client
	logDir string
}

func New(host, token string, logDir string) *client {
	return &client{
		host:  host,
		token: token,
		client: &http.Client{
			Timeout: 1 * time.Minute,
		},
		logDir: logDir,
	}
}

func (c *client) ExportDay(ctx context.Context, date time.Time) (*Day, error) {
	var day Day
	businessDay := date.Format("2006-01-02")
	path := fmt.Sprintf("export?business-day=%s", businessDay)
	if err := c.do(ctx, path, &day); err != nil {
		return nil, err
	}
	return &day, nil
}

func (c *client) ExportMaster(ctx context.Context, filters ...string) (*Master, error) {
	var master Master
	var filter string
	if len(filters) > 0 {
		filter = fmt.Sprintf("?filter=%s", strings.Join(filters, ","))
	}
	path := fmt.Sprintf("export-master/%s", filter)
	if err := c.do(ctx, path, &master); err != nil {
		return nil, err
	}
	return &master, nil
}

func (c *client) do(ctx context.Context, path string, out any) error {
	u := fmt.Sprintf("%s/api/%s", c.host, path)
	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return fmt.Errorf("agora: couldn't create request: %w", err)
	}
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Api-Token", c.token)

	log.Println("agora: request", u)
	resp, err := c.client.Do(req)
	if err != nil {
		return fmt.Errorf("agora: couldn't make request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("agora: couldn't read response body: %w", err)
	}

	// Write response body to log file
	f := strings.Trim(path, "/")
	f = strings.ReplaceAll(f, "/", "-")
	f = strings.ReplaceAll(f, "?", "-")
	f = strings.ReplaceAll(f, "=", "-")
	f = filepath.Join(c.logDir, fmt.Sprintf("%s_%s.json", f, time.Now().Format("20060102_150405")))
	if err := os.WriteFile(f, body, 0644); err != nil {
		log.Println(fmt.Errorf("agora: couldn't write file: %w", err))
	}

	if resp.StatusCode != http.StatusOK {
		log.Println(string(body))
		return fmt.Errorf("agora: unexpected status code: %d", resp.StatusCode)
	}

	if err := json.Unmarshal(body, out); err != nil {
		return fmt.Errorf("agora: couldn't unmarshal json: %w", err)
	}
	return nil
}
