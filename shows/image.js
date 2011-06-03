function(doc, req) {
    var mustache = require("vendor/couchapp/lib/mustache");
    var path = require("vendor/couchapp/lib/path").init(req);
    var markdown = require("vendor/couchapp/lib/markdown");

    doc.imageLink = path.attachment(doc._id, 'full');
    return mustache.to_html(this.templates.pic, doc);

}