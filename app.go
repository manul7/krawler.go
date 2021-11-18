package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"

	"net/http"
	"net/url"
	"os"

	"github.com/PuerkitoBio/goquery"
	"github.com/urfave/cli/v2"
)

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


func fetch(url string, ch chan<- string) {
	resp, err := http.Get(url)
	if err != nil {
		ch <- fmt.Sprint(err)
	}
	// Read response body and close stream
	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()

	if err != nil {
		ch <- fmt.Sprintf("error while reading %s: %v", url, err)
		return
	}
	// Load HTML content from recieved bytes
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		log.Fatal(err)
	}

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		title := s.Text()
		link, _ := s.Attr("href")
		fmt.Printf("Post #%d: %s - %s\n", i, title, link)
	})

	ch <- fmt.Sprintf("Read: %7d %s", 0, url)
}

// TODO:How to deal with invalid SSL certs?

func main() {
	app := &cli.App{
		Name:      "krawler",
		Usage:     "It takes URL string, processes it, and stores allowed content into DST directory.",
		UsageText: "krawler <URL> <Output dir>",
		Action: func(c *cli.Context) error {
			if c.Args().Len() < 2 {
				fmt.Printf("Please provide all required arguments: %v\n", c.App.UsageText)
				return nil
			}
			// get URL from user
			url_string := c.Args().Get(0)
			// TODO: Add empty string handler
			// get DST dir from user
			dst := c.Args().Get(1)
			// TODO: Add empty string handler
			if dst == "" {
				fmt.Printf("SSS")
			}
			// validate/parse URL
			uri, err := url_parse(url_string)
			// validate/parse DST
			createDstDir(dst)
			// if OKx2 -> continue, else -> return error
			fmt.Printf("Input: %q\n", url_string)
			fmt.Printf("URL: %q \nDST: %q\n", uri.Host, dst)
			ch := make(chan string)
			go fetch(uri.String(), ch)
			fmt.Println(<-ch)
			return err
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

}
