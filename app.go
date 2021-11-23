package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli/v2"
)

type myURL struct {
	url string
	body                    []byte
}
// Parse & check provided URL string
func url_parse(s string) (*url.URL, error) {
	uri, err := url.Parse(s)

	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return uri, nil
}

// Create output directory for fetched content
func createDstDir(p string) {
	// Check if directory do exist
	if _, err := os.Stat(p);
	// Not exists - create entire path
	os.IsNotExist(err) {
		// TODO: Error handler?
		os.MkdirAll(p, os.ModePerm)
	}
}

// Fetch one URL per time and send page body to next step
func (u *myURL) fetch() {
	log.Printf("Fetching %v", u.url)
	resp, err := http.Get(u.url)
	if err != nil {
		log.Fatalln(err)
	}
	// Read response body and close stream
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("error while reading %s: %v", u.url, err)
		return
	}
	u.body = b
}

// Parse HTML body and extrac HREFs
func (u *myURL) getURLs() []string {
	var hrefs []string
	// Load HTML content from recieved bytes
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(u.body))
	if err != nil {
		log.Fatal(err)
	}
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		link, _ := s.Attr("href")
		hrefs = append(hrefs, link)
	})
	// Send extracted HREFs to next step
	return hrefs
}

// TODO:How to deal with invalid SSL certs?

// Process input arguments
func parse_arguments(c *cli.Context) (string, string, error) {
			if c.Args().Len() < 2 {
		log.Printf("Please provide all required arguments: %v\n", c.App.UsageText)
		return "", "", nil
			}
	base_url := c.Args().Get(0)

			dst := c.Args().Get(1)
			if dst == "" {
		log.Printf("Please provide non empty output path")
		return "", "", nil
	}

	_, err := url_parse(base_url)
	if err != nil {
		log.Fatal("Cannot parse URL, please provide correct URL")
		return "", "", err
	}

	return base_url, dst, nil
}

// Main URL-processing function
func crawl(url string, ch chan []string) {
	var extracted []string
	u := myURL{url: url}
	u.fetch()
	extracted = u.getURLs()
	log.Printf("Saving data from %v", u.url)
	ch <- extracted
}

// Validate URL string against rules
func check_url(base_url *string, link *string) error {
	log.Printf("LINK: %v. HTTP? %v", link, strings.HasPrefix("http", *link))
	if !strings.HasPrefix(*base_url, *link) && !strings.HasPrefix("http", *link) {
		*link = *base_url + *link
	}
	_, err := url_parse(*link)
	return err
	}
func crawl_cmd(c *cli.Context) error {
	base_url, dst, err := parse_arguments(c)
	if err != nil {
		panic("Error. Please check input params")
			}
	ch := make(chan []string, 4)

	log.Printf("Base URL: %q", base_url)
	log.Printf("Output dir: %q", dst)

			createDstDir(dst)
	// Setup Ctrl-C handler
	signals_ch := make(chan os.Signal, 1)
	signal.Notify(signals_ch, os.Interrupt)
	// Send base_url for processing
	// TODO: Looks ugly, need to find better solution
	go func(b string) {
		log.Printf("Sending base URL %v", b)
		ch <- []string{b}
	}(base_url)

	for {
		select {
		case links := <-ch:
			for _, link := range links {
				err = check_url(&base_url, &link)
				if err != nil {
					log.Printf("Incorrect URL: %v, %v", link, err)
					continue
				}

				go crawl(link, ch)
			}
		case <-signals_ch:
			log.Println("Shutting down...")
			return nil
		}
	}
	}

func main() {
	// Setup CLI commands
	app := &cli.App{
		Name:      "krawler",
		Usage:     "It takes URL string, processes it, and stores allowed content into DST directory.",
		UsageText: "krawler <URL> <Output dir>",
		Action:    crawl_cmd,
	}
	// Launch app
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
