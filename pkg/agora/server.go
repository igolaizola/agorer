package agora

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
)

func MockServe(ctx context.Context, addr string, masterFile string) (int, error) {
	b, err := os.ReadFile(masterFile)
	if err != nil {
		return 0, fmt.Errorf("couldn't read %s: %w", masterFile, err)
	}

	// Listen on addr
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return 0, fmt.Errorf("couldn't listen on %s: %w", addr, err)
	}

	// Serve masterFile json file on /export-master endpoint
	http.HandleFunc("/api/export-master/", func(w http.ResponseWriter, r *http.Request) {
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
