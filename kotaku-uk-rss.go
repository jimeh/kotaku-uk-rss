package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/gorilla/feeds"
)

/*
  Structs
*/

// Article holds various info about an article.
type Article struct {
	Title    string
	Date     time.Time
	Author   string
	Summary  string
	MediaURL string
	URL      string
	Href     string
}

// Articles is a collection of Article structs.
type Articles []Article

// SectionsByArticleLength enables sorting section nodes by number of article
// child nodes.
type SectionsByArticleLength []*goquery.Selection

func (s SectionsByArticleLength) Len() int {
	return len(s)
}

func (s SectionsByArticleLength) Less(i int, j int) bool {
	return s[i].Find("article").Length() > s[j].Find("article").Length()
}

func (s SectionsByArticleLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// RssCache is a string cache with mutex locks.
type RssCache struct {
	sync.RWMutex
	body string
}

func (r *RssCache) Set(value string) {
	r.Lock()
	defer r.Unlock()
	r.body = value
}

func (r *RssCache) Get() string {
	r.RLock()
	defer r.RUnlock()
	return r.body
}

/*
  Functions
*/

func fetchDocument(url string) *goquery.Document {
	doc, err := goquery.NewDocument(url)
	if err != nil {
		log.Fatal(err)
	}

	return doc
}

func extractArticleSection(doc *goquery.Document) *goquery.Selection {
	var sections = SectionsByArticleLength{}
	doc.Find("section").Each(func(i int, s *goquery.Selection) {
		sections = append(sections, s)
	})

	sort.Sort(sections)
	return sections[0]
}

func parseArticleSection(section *goquery.Selection) Articles {
	result := Articles{}

	section.Find("article").Each(func(i int, s *goquery.Selection) {
		result = append(result, parseArticle(s))
	})

	return result
}

func parseArticle(s *goquery.Selection) Article {
	href, _ := s.Find(".media__body h2 a").Attr("href")
	if href == "" {
		href, _ = s.Find("figure a").Attr("href")
	}
	url := rootURL + href

	summary := s.Find(".media__body p").Text()
	title := s.Find(".media__body h2").Text()
	if title == "" {
		title = truncateString(summary, 60) + "..."
	}

	mediaURL, _ := s.Find("figure").Attr("data-media992")
	timeString, _ := s.Find(".meta__limited time").Attr("datetime")
	parsedTime, _ := time.Parse(time.RFC3339, timeString)

	return Article{
		Title:    title,
		Date:     parsedTime,
		Author:   s.Find(".meta__full a.is-author").Text(),
		Summary:  summary,
		MediaURL: mediaURL,
		URL:      url,
		Href:     href,
	}
}

func truncateString(s string, l int) string {
	end := len(s)
	if end > l {
		end = l
	}
	return s[:end]
}

func getArticlesFromUrl(url string) Articles {
	doc := fetchDocument(url)
	section := extractArticleSection(doc)
	return parseArticleSection(section)
}

func buildFeed(articles Articles) *feeds.Feed {
	feed := &feeds.Feed{
		Title:       "Kotaku UK",
		Link:        &feeds.Link{Href: "http://www.kotaku.co.uk/"},
		Description: "Kotaku UK is the UK version of Kotaku",
		Created:     time.Now(),
	}

	feed.Items = []*feeds.Item{}

	for _, article := range articles {
		feed.Items = append(feed.Items, buildFeedItem(article))
	}

	return feed
}

func buildFeedItem(article Article) *feeds.Item {
	description := article.Summary

	if article.MediaURL != "" {
		description = fmt.Sprintf("<img href=\"%s\" /> %s",
			article.MediaURL, description)
	}

	return &feeds.Item{
		Id:          article.Href,
		Title:       article.Title,
		Link:        &feeds.Link{Href: article.URL},
		Description: description,
		Author:      &feeds.Author{Name: article.Author},
		Created:     article.Date,
	}
}

func serveRss(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, rssCache.Get())
}

func updateRssLoop() {
	for {
		articles := Articles{}
		for _, url := range pageUrls {
			fmt.Printf("fetching and parsing articles from: %s\n", url)
			for _, article := range getArticlesFromUrl(url) {
				articles = append(articles, article)
			}
		}

		fmt.Printf("building feed cache... ")
		feed, _ := buildFeed(articles).ToRss()
		rssCache.Set(feed)
		fmt.Println("done")
		fmt.Println("taking a nap for 60 seconds ^_^")
		time.Sleep(60 * time.Second)
	}
}

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

func main() {
	port := "80"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}

	http.HandleFunc("/rss", serveRss)

	go updateRssLoop()

	fmt.Println("Listing on port " + port + "...")
	err := http.ListenAndServe(":"+port, nil)
	if err != nil {
		panic(err)
	}
}
