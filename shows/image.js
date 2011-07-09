function(doc, req) {
    var mustache = require("vendor/couchapp/lib/mustache");
    var path = require("vendor/couchapp/lib/path").init(req);
    var markdown = require("vendor/couchapp/lib/markdown");

    doc.imageLink = path.attachment(doc._id, 'full');
    doc.state = doc.state || 'boring';
    doc.boring_display = doc.state === 'boring' ? "block" : "none";
    doc.fave_display = doc.state === 'fave' ? "block" : "none";
    return mustache.to_html(this.templates.pic, doc);

}