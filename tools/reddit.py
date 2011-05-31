#!/usr/bin/env python

import os
import sys
import time
import urllib
import base64
import hashlib
import traceback

import couchdb
import BeautifulSoup

import feedparser

SERVER = os.getenv("COUCHDB") or 'http://127.0.0.1:5984/'

DB = couchdb.Server(SERVER)['rpics']

ISO8601 = "%Y-%m-%dT%H:%M:%S"

def getImage(url, recurse=True):
    o = urllib.urlopen(url)
    ctype = o.headers['content-type']
    content = o.read()
    if not ctype.startswith('image/'):
        # Maybe I can find the image...
        s = BeautifulSoup.BeautifulSoup(content)
        src = s.find('link', rel='image_src')
        if src:
            print "\tFound original at", src['href']
            return getImage(src['href'], False)

        if not recurse:
            raise NotFound()
        else:
            return getImage(url + '.jpg', False)
    return ctype, content

def handle(sub, e):
    s = BeautifulSoup.BeautifulSoup(e.summary_detail.value)
    u = s.findAll(lambda x: x.name == 'a' and x.findAll(text='[link]'))[0]['href']

    shortlink = s.findAll("img")[0]
    su = shortlink['src']
    st = shortlink['title']

    md5 = hashlib.md5()
    md5.update(st.encode('utf-8'))
    docid = md5.hexdigest()

    if DB.get(docid):
        print "-", st
        return

    print "+", st

    thumbtype, thumb = getImage(su)
    fulltype, full = getImage(u)

    doc = {'_id': docid,
           'sub': sub,
           'updated': time.strftime("%Y-%m-%dT%H:%M:%S", e.updated_parsed),
           'title': st,
           'url': e.link,
           '_attachments': {
            'thumb': {'content_type': thumbtype, 'data': base64.b64encode(thumb)},
            'full': {'content_type': fulltype, 'data': base64.b64encode(full)}
            }}

    DB.save(doc)

if __name__ == '__main__':
    for sub in sys.argv[1:]:
        url = 'http://www.reddit.com/r/' + sub + '/.rss'
        f = feedparser.parse(url)

        for e in f.entries:
            try:
                handle(sub, e)
            except:
                traceback.print_exc()
