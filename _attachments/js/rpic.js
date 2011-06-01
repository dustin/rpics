function rpic_recent_feed(app, target) {
    var path = app.require("vendor/couchapp/lib/path").init(app.req);
    var Mustache = app.require("vendor/couchapp/lib/mustache");
    var maxItems = 100;

    var template = '<a href="{{full}}"><img title="{{title}}" src="{{thumb}}" /></a>';
    app.db.info({success: function(dbi) {
        var since = Math.max(0, dbi.update_seq - (maxItems + (maxItems * 0.5)));
        console.log("Starting fetch at", since, dbi);

        var changeFeed = app.db.changes(since, {"include_docs": true});
        changeFeed.onChange(function(data) {
            var nItems = 0;
            data.results.forEach(function(row) {
                if (row.doc['_attachments'] && row.doc._attachments['thumb']) {
                    var tdata = {
                        title: row.doc.title,
                        thumb: path.attachment(row.id, 'thumb'),
                        full: path.attachment(row.id, 'full')
                    };
                    target.prepend(Mustache.to_html(template, tdata));
                    target.find("a:gt(" + maxItems + ")").remove();
                }
            });
        });
    }});
}