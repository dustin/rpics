var rpicChangeFeed = undefined;
var rpicDisplayed = { };

function rpic_recent_feed(app, target) {
    var path = app.require("vendor/couchapp/lib/path").init(app.req);
    var Mustache = app.require("vendor/couchapp/lib/mustache");
    var maxItems = 100;

    var template = '<a href="{{link}}"><img title="{{title}}" src="{{thumb}}" /></a>';
    app.db.info({success: function(dbi) {
        var since = Math.max(0, dbi.update_seq - (maxItems + (maxItems * 0.5)));
        console.log("Starting fetch at", since, dbi);

        rpicChangeFeed = app.db.changes(since, {"include_docs": true});
        rpicChangeFeed.onChange(function(data) {
            var nItems = 0;
            data.results.forEach(function(row) {
                if (!rpicDisplayed[row.id] && row.doc['_attachments'] && row.doc._attachments['thumb']) {
                    rpicDisplayed[row.id] = true;
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

function rpic_feed_toggle(app, target) {
    setTimeout(function() {
        if (document.webkitHidden) {
            if (rpicChangeFeed) {
                console.log("Stopping the feed.");
                rpicChangeFeed.stop();
                rpicChangeFeed = undefined;
            }
        } else {
            if (!rpicChangeFeed) {
                console.log("Starting the feed.");
                rpic_recent_feed(app, target);
            }
        }
    }, 100);
}

function rpic_init_feed(app, target) {
    document.addEventListener('webkitvisibilitychange', function(e) {
        rpic_feed_toggle(app, target);
    }, false);
    rpic_feed_toggle(app, target);
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