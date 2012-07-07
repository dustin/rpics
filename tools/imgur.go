package main

import (
	"encoding/json"
	"io"
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
