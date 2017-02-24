package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"gopkg.in/alecthomas/kingpin.v2"
)

// Version gets populated with version at build-time.
var Version string
var defaultPort = "8080"

var (
	port = kingpin.Flag("port", "Port to listen to.").Short('p').
		Default(defaultPort).String()
	bind = kingpin.Flag("bind", "Bind address.").Short('b').
		Default("0.0.0.0").String()
	version = kingpin.Flag("version", "Print version info.").
		Short('v').Bool()
)

/*
  Configuration
*/

var rootURL = "http://www.kotaku.co.uk"

var pageUrls = []string{
	rootURL + "/",
	rootURL + "/page/2/",
	rootURL + "/page/3/",
}

/*
  Main...
*/

var rssCache = RssCache{}

func printVersion() {
	fmt.Println("kotaku-uk-rss " + Version)
}

func serveRoot(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}

func startServer() {
	if *port == defaultPort {
		if envPort := os.Getenv("PORT"); envPort != "" {
			*port = envPort
		}
	}

	http.HandleFunc("/", serveRoot)
	http.HandleFunc("/rss", serveRss)

	address := *bind + ":" + *port
	fmt.Println("Listening on " + address)
	log.Fatal(http.ListenAndServe(address, nil))
}

func main() {
	kingpin.Parse()

	if *version {
		printVersion()
	} else {
		go updateRssLoop()
		startServer()
	}
}
