package main

import (
	"os"
	"testing"
)

func TestFlickrImage(t *testing.T) {
	f, err := os.Open("flickrexample.html")
	if err != nil {
		t.Fatalf("Error opening flickr example: %v", err)
	}
	u, err := parseFlickr(f)
	if err != nil {
		t.Fatalf("Error in flickr example: %v", err)
	}

	exp := "http://farm3.staticflickr.com/2413/5769871946_e69f2734f4_z.jpg"
	if u != exp {
		t.Fatalf("Expected %v for URL, got %v", exp, u)
	}
}
