package main

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"code.google.com/p/dsallings-couch-go"
)

var couchURL = flag.String("couchdb", "http://localhost:5984/rpics2",
	"The CouchDB into which we store all the things.")
var syslogFlag = flag.Bool("syslog", false, "Log to syslog")
var watchdog = flag.Duration("watchdog", 5*time.Minute,
	"Maximum amount of time permitted to run")

type myRT struct {
	rt     http.RoundTripper
	delays map[string]time.Time
	l      sync.Mutex
}

func (m *myRT) RoundTrip(req *http.Request) (resp *http.Response, err error) {
	func() {
		m.l.Lock()
		defer m.l.Unlock()

		w, ok := m.delays[req.URL.Host]
		if ok && w.After(time.Now()) {
			st := w.Sub(time.Now())
			log.Printf("Waiting %v to request %v", st, req.URL)
			time.Sleep(st)
		}
		if strings.Contains(req.URL.Host, "reddit.com") {
			m.delays[req.URL.Host] = time.Now().Add(time.Second * 2)
		} else if strings.Contains(req.URL.Host, "imgur.com") {
			m.delays[req.URL.Host] = time.Now().Add(time.Millisecond * 200)
		}
	}()
	req.Header.Set("User-Agent", "Image coucher by /u/dlsspy")
	return m.rt.RoundTrip(req)
}

func init() {
	delays := map[string]time.Time{}
	oldrt := http.DefaultTransport
	http.DefaultTransport = &myRT{rt: oldrt, delays: delays}
}

type Posting struct {
	Domain    string
	Permalink string
	Thumbnail string
	Title     string
	Subreddit string
	URL       string
}

type Response struct {
	Data struct {
		Children []struct {
			Data Posting
		}
	}
}

type Image struct {
	ContentType string
	MD5         string
	Encoded     string
}

func (i *Image) MarshalJSON() ([]byte, error) {
	m := map[string]string{
		"content_type": i.ContentType,
		"data":         i.Encoded,
	}
	return json.Marshal(m)
}

type Document struct {
	Id          string            `json:"_id"`
	Updated     time.Time         `json:"updated"`
	Sub         string            `json:"sub"`
	Title       string            `json:"title"`
	URL         string            `json:"url"`
	Attachments map[string]*Image `json:"_attachments"`
}

func getImageRaw(u string) (Image, error) {
	resp, err := http.Get(u)
	if err != nil {
		return Image{}, err
	}
	defer resp.Body.Close()

	rv := Image{}

	buf := &bytes.Buffer{}
	b64 := base64.NewEncoder(base64.StdEncoding, buf)
	h := md5.New()
	io.Copy(b64, io.TeeReader(resp.Body, h))

	rv.ContentType = resp.Header.Get("Content-Type")
	rv.MD5 = fmt.Sprintf("%x", h.Sum(nil))
	rv.Encoded = buf.String()

	return rv, err
}

func getImage(p Posting) (Image, error) {
	if strings.Contains(p.Domain, "imgur.com") {
		return getImageImgur(p)
	} else if strings.Contains(p.Domain, "flickr.com") {
		return getImageFlickr(p)
	} else {
		return getImageRaw(p.URL)
	}
	return Image{}, nil
}

func parseResponse(r io.Reader) ([]Posting, error) {
	d := json.NewDecoder(r)
	rs := Response{}
	err := d.Decode(&rs)
	if err != nil {
		return []Posting{}, err
	}

	rv := make([]Posting, 0, len(rs.Data.Children))
	for _, p := range rs.Data.Children {
		rv = append(rv, p.Data)
	}
	return rv, nil
}

func hasDoc(db *couch.Database, docid string) bool {
	doc := map[string]interface{}{}
	err := db.Retrieve(docid, &doc)
	if err != nil {
		if !strings.Contains(err.Error(), "404") {
			log.Printf("Hasdoc %v:  %v", docid, err)
		}
	}
	return err == nil
}

func process(db *couch.Database, p Posting, wg *sync.WaitGroup) {
	defer wg.Done()

	m := md5.New()
	io.WriteString(m, p.Permalink)
	docid := fmt.Sprintf("%x", m.Sum(nil))

	if p.Thumbnail == "" {
		log.Printf("No thumbnail for %v", p.Title)
		return
	}

	if hasDoc(db, docid) {
		return
	}

	img, err := getImage(p)
	if err != nil {
		if err != noFlickrUrl {
			log.Printf("Error getting image %v", err)
		}
		return
	}
	if !strings.HasPrefix(img.ContentType, "image/") {
		log.Printf("Invalid content type in image from %v: %#v",
			p, img.ContentType)
		return
	}
	var thumb Image
	if p.Thumbnail == "nsfw" {
		thumb.Encoded = nsfw_png_b64
		thumb.ContentType = "image/png"
	} else {
		thumb, err = getImageRaw(p.Thumbnail)
		if err != nil {
			log.Printf("Error getting thumbnail from %#v, %v", p, err)
			return
		}
	}
	log.Printf("Got %v bytes of image and %v of thumbnail",
		len(img.Encoded), len(thumb.Encoded))

	doc := Document{
		Id:      docid,
		Updated: time.Now(),
		Sub:     p.Subreddit,
		Title:   p.Title,
		URL:     "http://www.reddit.com" + p.Permalink,
		Attachments: map[string]*Image{
			"thumb": &thumb,
			"full":  &img,
		},
	}
	_, _, err = db.Insert(&doc)
	if err != nil {
		log.Printf("Error storing %v: %v", doc.Title, err)
	}
}

func grabStuff(db *couch.Database, sub string) error {
	url := "http://www.reddit.com/r/" + sub + "/.json"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		log.Printf("Bad status: %v", resp.Status)
		return fmt.Errorf("Bad status: %v", err)
	}

	posting, err := parseResponse(resp.Body)
	if err != nil {
		return err
	}

	wg := sync.WaitGroup{}
	wg.Add(len(posting))
	for _, p := range posting {
		go process(db, p, &wg)
	}
	wg.Wait()
	return nil
}

func main() {
	flag.Parse()

	time.AfterFunc(*watchdog, func() {
		log.Fatalf("Watchdog firing, taking too long.")
	})

	setupLogging(*syslogFlag)

	db, err := couch.Connect(*couchURL)
	if err != nil {
		log.Fatalf("Error connecting to couchdb: %v", err)
	}

	wg := sync.WaitGroup{}
	for _, s := range flag.Args() {
		wg.Add(1)
		go func(which string) {
			defer wg.Done()
			err := grabStuff(&db, which)
			if err != nil {
				log.Printf("Got %v for %v", err, which)
			}
		}(s)
	}
	wg.Wait()
}
