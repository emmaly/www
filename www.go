package main

import (
	"fmt"
	"net/http"
	"os"
	"path"
	"strings"

	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	cwd, _ = os.Getwd()
	dir    = kingpin.Arg("path", "Serve path").Default(cwd).String()
	host   = kingpin.Flag("host", "Listen host").Short('h').Default("0.0.0.0").String()
	port   = kingpin.Flag("port", "Listen port").Short('p').Default("8080").String()
)

func main() {
	kingpin.Parse()
	fmt.Printf("Listening on %s:%s and serving %s\n", *host, *port, *dir)
	http.Handle("/", headers(http.FileServer(http.Dir(*dir))))
	panic(http.ListenAndServe(*host+":"+*port, nil))
}

type headersMiddleware struct {
	handler http.Handler
}

func headers(handler http.Handler) http.Handler {
	return &headersMiddleware{handler}
}

// ServeHTTP fulfills http.Handler
func (hm headersMiddleware) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch strings.ToLower(path.Ext(r.URL.Path)) {
	case "wasm":
		w.Header().Set("Content-type", "application/wasm")
	}
	hm.handler.ServeHTTP(w, r)
}
