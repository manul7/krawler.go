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
	url  url.URL
	body []byte
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
func (u *myURL) fetch() error {
	log.Printf("Fetching %v", u.url.String())
	resp, err := http.Get(u.url.String())
	if err != nil {
		log.Print(err)
		return err
	}

	// Read response body and close stream
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		log.Fatalf("error while reading %s: %v", u.url.String(), err)
		return err
	}

	u.body = b
	return nil
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
func parse_arguments(c *cli.Context) (url.URL, string, error) {
			if c.Args().Len() < 2 {
		log.Fatalf("Please provide all required arguments: %v\n", c.App.UsageText)
			}
	inputURL := c.Args().Get(0)

			dst := c.Args().Get(1)
			if dst == "" {
		log.Fatal("Please provide non empty output path")
	}

	baseURL, err := buildBaseUrl(&inputURL)
	if err != nil {
		log.Fatal("Cannot parse URL, please provide correct URL")
	}

	return baseURL, dst, nil
}

// Build BaseURL for crawler
func buildBaseUrl(s *string) (url.URL, error) {
	uri, err := url.Parse(*s)
	if err != nil {
		log.Fatal(err)
		return *uri, err
	}
	// In case when no scheme provided - add default "HTTP"
	if uri.Scheme == "" {
		uri.Scheme = "http"
	}
	return *uri, nil
	}

// Save page content into file
func (u *myURL) saveContent(dst *string) {
	resultPath := filepath.Join(*dst, u.url.Path)
	createDstDir(resultPath)
	log.Printf("Saving to %v", resultPath)
	// write the whole body at once
	outFilePath := filepath.Join(resultPath, "index.html")
	err := ioutil.WriteFile(outFilePath, u.body, 0644)
	if err != nil {
		panic(err)
	}
}

// Main URL-processing function
func crawl(dst *string, url *url.URL, ch chan []string, wg *sync.WaitGroup) {
	// Will decrement when all crawling functions finished
	defer wg.Done()
	u := myURL{url: *url}
	err := u.fetch()

	if err != nil {
		log.Printf("Error while fetching: %v", err)
		return
	}

	log.Printf("Saving data from %v", u.url.String())
	u.saveContent(dst)
	ch <- u.getURLs()
}

// Validate URL string against rules
func isInScope(baseUrl *url.URL, link *url.URL) bool {
	if link.Host == "" {
		link.Host = baseUrl.Host
		link.Scheme = baseUrl.Scheme
	}

	if link.Host != baseUrl.Host {
		log.Printf("Error: URL does not belong to starting host: %v", link.Host)
		return false
	}
	if !strings.HasPrefix(link.String(), baseUrl.String()) {
		log.Printf("Error: URL does not belong to base URL: %v | %v", baseUrl.String(), link.String())
		return false
	}
	return true
}

func crawl_cmd(c *cli.Context) error {
	base_url, dst, err := parse_arguments(c)
	if err != nil {
		panic("Error. Please check input params")
			}
	ch := make(chan []string, 2)
	var wg sync.WaitGroup

	log.Printf("Base URL: %q", base_url.String())
	log.Printf("Output dir: %q", dst)

			createDstDir(dst)
	// Setup Ctrl-C handler
	signals_ch := make(chan os.Signal, 1)
	signal.Notify(signals_ch, os.Interrupt)
	// Send base_url for processing
	// TODO: Looks ugly, need to find better solution
	go func(b string, wg sync.WaitGroup) { ch <- []string{b} }(base_url.String(), wg)

	for {
		select {
		// Incoming URLs handler
		case links := <-ch:
			for _, link := range links {
				uri, err := url.Parse(link)
				if err != nil {
					log.Printf("Error: Incorrect URI: %v", link)
					continue
				}

				inScope := isInScope(&base_url, uri)
				if !inScope {
					continue
			}
				wg.Add(1)
				go crawl(&dst, uri, ch, &wg)
			}
			// wg.Wait()
		// Ctrl-C handler
		case <-signals_ch:
			log.Println("Shutting down...")
			close(ch)
			wg.Wait()
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
