package scraper

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestCrawl(t *testing.T) {
	flag.Parse()

	if len(os.Args) < 2 {
		fmt.Printf("usage: %s [url]\n", os.Args[0])
		os.Exit(1)
	}

	transport, err := NewTransport(http.DefaultTransport)
	if err != nil {
		log.Fatal(err)
	}

	client := &http.Client{Transport: transport}
	makeRequest(client, flag.Arg(0))
}

func TestTransport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, err := ioutil.ReadFile("_examples/challenge.html")
		if err != nil {
			t.Fatal(err)
		}
		w.Header().Set("Content-Type", "text/html")
		w.Header().Set("Server", "cloudflare-nginx")
		w.WriteHeader(503)
		w.Write(b)
	}))
	defer ts.Close()

	scraper, err := NewTransport(http.DefaultTransport)
	if err != nil {
		t.Fatal(err)
	}

	c := http.Client{
		Transport: scraper,
	}

	res, err := c.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}

	_, err = ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}
}

func makeRequest(c *http.Client, url string) {
	t := time.Now()

	log.Printf("Requesting %s", url)
	resp, err := c.Get(url)
	if err != nil {
		log.Fatal(err)
	}

	body, _ := ioutil.ReadAll(resp.Body)
	log.Printf("Fetched %s in %s, %d bytes (status %d)",
		url, time.Now().Sub(t), len(body), resp.StatusCode)

	if resp.StatusCode != http.StatusOK {
		log.Fatal("Invalid response code")
	}
}
