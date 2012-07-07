package main

import (
	"errors"
	"exp/html"
	"io"
	"net/http"
)

func parseFlickr(r io.Reader) (rv string, err error) {
	doc, err := html.Parse(r)
	if err != nil {
		return "", err
	}

	var f func(n *html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "meta" {
			content := ""
			isImage := false
			for _, a := range n.Attr {
				if a.Key == "property" && a.Val == "og:image" {
					isImage = true
				} else if a.Key == "content" {
					content = a.Val
				}
			}
			if isImage {
				rv = content
				return
			}
		}
		for _, c := range n.Child {
			if rv == "" {
				f(c)
			}
		}
	}
	f(doc)
	if rv == "" {
		err = errors.New("Flickr URL not found")
	}
	return rv, err
}

func getImageFlickr(p Posting) (Image, error) {
	resp, err := http.Get(p.URL)
	if err != nil {
		return Image{}, err
	}
	defer resp.Body.Close()

	u, err := parseFlickr(resp.Body)
	if err != nil {
		return Image{}, err
	}

	return getImageRaw(u)
}
