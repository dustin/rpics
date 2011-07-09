function rpic_recent_feed(app, target) {
    var path = app.require("vendor/couchapp/lib/path").init(app.req);
    var Mustache = app.require("vendor/couchapp/lib/mustache");
    var maxItems = 100;

    var template = '<a href="{{link}}"><img title="{{title}}" src="{{thumb}}" /></a>';
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
                        link: path.show('image', row.id)
                    };
                    target.prepend(Mustache.to_html(template, tdata));
                    target.find("a:gt(" + maxItems + ")").remove();
                }
            });
        });
    }});
}

function rpic_init_update_links(app) {
    var baseUri = app.db.uri;
    var ddoc = app.ddoc._id;

    $(".statechange").each(function(a, el) {
        var parts = el.id.split('-');
        $(el).click(function() {
            $.ajax({type: 'POST',
                    url: baseUri + ddoc + "/_update/set_state/" + parts[1],
                    data: 'new_state=' + encodeURIComponent(parts[0]),
                    dataType: "json",
                    complete: function(res) {
                        var new_state = parts[0];
                        var old_state = new_state === 'fave' ? 'boring' : 'fave';
                        $(".state-" + new_state).show();
                        $(".state-" + old_state).hide();
                        console.log("Result", res, old_state, "->", new_state);
                    }});

            return false;
        });
    });
}