#!/usr/bin/env python

import os
import sys
import json
import couchdb

db = couchdb.Server(os.getenv("COUCHDB", "http://127.0.0.1:5984/"))['rpics']

todo = [{'_id': d.id, '_rev': d.value, '_deleted': True} for d in db.view('rp/uninteresting', limit=100)]

db.update(todo)
print "Deleted", len(todo), "images"
