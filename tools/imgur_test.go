package main

import (
	"strings"
	"testing"
)

var imgurlsample = `{"image":{"image":{"title":null,"caption":null,"hash":"bW1mm","datetime":"2012-07-06 19:27:22","type":"image\/jpeg","animated":"false","width":2072,"height":1434,"size":183453,"views":489794,"bandwidth":89854178682},"links":{"original":"http:\/\/i.imgur.com\/bW1mm.jpg","imgur_page":"http:\/\/imgur.com\/bW1mm","small_square":"http:\/\/i.imgur.com\/bW1mms.jpg","large_thumbnail":"http:\/\/i.imgur.com\/bW1mml.jpg"}}}`

func TestImgURParsing(t *testing.T) {
	d, err := parseImgUr(strings.NewReader(imgurlsample))
	if err != nil {
		t.Fatalf("Error parsing IMG UR data: %v", err)
	}
	t.Logf("%#v", d)
}
