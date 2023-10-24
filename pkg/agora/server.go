package agora

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

func MockServe(ctx context.Context, addr string, masterFile string) (int, error) {
	master, err := os.ReadFile(masterFile)
	if err != nil {
		return 0, fmt.Errorf("couldn't read %s: %w", masterFile, err)
	}

	days := map[string][]byte{}

	// Loop all files in master file directory
	dir := filepath.Dir(masterFile)
	err = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return fmt.Errorf("couldn't walk %s: %w", path, err)
		}
		base := filepath.Base(path)
		// Skip directories
		if info.IsDir() {
			return nil
		}
		// Check if file json
		if filepath.Ext(path) != ".json" {
			return nil
		}
		// Check if file name is a date + .json
		date := base[:len(base)-len(".json")]
		if _, err := time.Parse("2006-01-02", date); err != nil {
			return nil
		}
		// Read file
		b, err := os.ReadFile(path)
		if err != nil {
			return fmt.Errorf("couldn't read %s: %w", path, err)
		}
		// Add to days map
		days[date] = b
		return nil
	})
	if err != nil {
		return 0, fmt.Errorf("couldn't walk %s: %w", dir, err)
	}

	// Listen on addr
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("couldn't listen on %s: %w", addr, err)
	}

	// Serve masterFile json file on /export-master endpoint
	http.HandleFunc("/api/export-master/", func(w http.ResponseWriter, r *http.Request) {
		w.Write(master)
	})

	// Serve day json files on /api/export endpoint
	http.HandleFunc("/api/export", func(w http.ResponseWriter, r *http.Request) {
		// Read ?business-day=%s query param
		q := r.URL.Query()
		date := q.Get("business-day")
		if date == "" {
			http.Error(w, "missing business-day query param", http.StatusBadRequest)
			return
		}
		b, ok := days[date]
		if !ok {
			http.Error(w, "no data for business-day", http.StatusNotFound)
			return
		}
		w.Write(b)
	})

	// Serve on addr until ctx is done
	go func() {
		http.Serve(l, nil)
	}()
	go func() {
		<-ctx.Done()
		l.Close()
	}()
	log.Println("Mocking agora server on", l.Addr().String())
	return l.Addr().(*net.TCPAddr).Port, nil
}
