package main

import (
	"os"
	"testing"
)

func TestSampleResponse(t *testing.T) {
	r, err := os.Open("sample.json")
	if err != nil {
		t.Fatalf("Error opening sample file: %v", err)
	}
	defer r.Close()

	rs, err := parseResponse(r)
	if err != nil {
		t.Fatalf("Error parsing response: %v", err)
	}

	if len(rs) != 25 {
		t.Fatalf("Expected 25 results, got %v", len(rs))
	}

	posting := Posting{
		Domain:    "imgur.com",
		Permalink: "/r/pics/comments/w568e/its_103_out_so_my_daughter_and_niece_waited_for/",
		Thumbnail: "http://a.thumbs.redditmedia.com/OoulK88mGUsZMp-B.jpg",
		Title: "It's 103+ out, so my daughter and niece waited for our garbage " +
			"men to come by so they could run out and give them some ice cold Gatorades ",
		URL:       "http://imgur.com/bW1mm",
		Subreddit: "pics",
	}

	if posting != rs[0] {
		t.Fatalf("Expected %v for the first entry, got %v", posting, rs[0])
	}
}
