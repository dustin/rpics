function(data) {
    var app = $$(this).app;
    var path = app.require("vendor/couchapp/lib/path").init(app.req);

    var items = data.rows.map(function(r) {
        var d = r.doc;
        return {
            title: d.title,
            thumb: path.attachment(r.id, 'thumb'),
            full: path.attachment(r.id, 'full')
        };
        return d;
    });

    return {
        items: items
    };
};