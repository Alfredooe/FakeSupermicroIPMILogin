package main

import (
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s", r.RemoteAddr, r.Method, r.URL)

		if r.Method == http.MethodPost {
			body, err := io.ReadAll(r.Body)
			if err != nil {
				log.Printf("Error reading body: %v", err)
			} else {
				log.Printf("%s", body)
			}
			r.Body.Close()

			http.Redirect(w, r, "https://www.youtube.com/watch?v=dQw4w9WgXcQ", http.StatusFound)
			return
		}

		delay := time.Duration(2+rand.Intn(5)) * time.Second
		time.Sleep(delay)

		crw := &customResponseWriter{ResponseWriter: w}
		next.ServeHTTP(crw, r)
	})
}

type customResponseWriter struct {
	http.ResponseWriter
}

func (crw *customResponseWriter) WriteHeader(statusCode int) {

	crw.Header().Set("Server", "GoAhead-Webs")
	crw.Header().Set("Date", time.Now().Format("Mon Jan 2 15:04:05 2006"))
	crw.Header().Set("Pragma", "no-cache")
	crw.Header().Set("Cache-Control", "no-cache")
	crw.Header().Set("Content-Type", "text/html")
	crw.Header().Set("Set-Cookie", "(null)")

	crw.Header().Del("Last-Modified")
	crw.Header().Del("Accept-Ranges")
	crw.Header().Del("Content-Length")

	crw.ResponseWriter.WriteHeader(statusCode)
}

func main() {
	rand.Seed(time.Now().UnixNano())

	server := http.FileServer(http.Dir("./static"))

	http.Handle("/", handler(server))

	log.Println("Starting server on :80")
	http.ListenAndServe(":80", nil)
}
