package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
)

type ImgURData struct {
	Image struct {
		Type     string
		Width    int
		Height   int
		Animated bool `json:",string"`
	}
	Links struct {
		Original string
	}
}

func parseImgUr(r io.Reader) (rv ImgURData, err error) {
	d := json.NewDecoder(r)
	crap := struct {
		Image ImgURData
	}{}
	err = d.Decode(&crap)
	return crap.Image, err
}

func getImageImgur(p Posting) (Image, error) {
	u, err := url.Parse(p.URL)
	if err != nil {
		return Image{}, err
	}
	cleanp := filepath.Base(u.Path)
	iuh := strings.Split(cleanp, ".")

	imguru := fmt.Sprintf("http://api.imgur.com/2/image/%v.json", iuh[0])
	resp, err := http.Get(imguru)
	if err != nil {
		return Image{}, err
	}
	defer resp.Body.Close()

	imgu, err := parseImgUr(resp.Body)
	if err != nil {
		return Image{}, err
	}

	rv, err := getImageRaw(imgu.Links.Original)
	if err == nil {
		rv.ContentType = imgu.Image.Type
	}
	return rv, err

}
