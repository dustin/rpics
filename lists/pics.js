function(head, req) {
	// !json templates.head
	// !json templates.tail

	provides("html", function() {
		var row;

		var data = {
			title: "A Bunch of Pictures",
            mainid: "thumblist"
		};

        var Mustache = require("vendor/couchapp/lib/mustache");
        var path = require("vendor/couchapp/lib/path").init(req);

		send(Mustache.to_html(templates.head, data));

        var template = '<a id="img/{{id}}" class="state-{{state}}" href="{{link}}"><img title="{{title}}" src="{{thumb}}" /></a>';

		while( (row = getRow()) ) {
			send(Mustache.to_html(template, {
                title: row.doc.title,
                thumb: path.attachment(row.id, 'thumb'),
                link: path.show('image', row.id),
                state: row.doc.state || 'boring',
                id: row.id + '/' + row.doc._rev
			}));
		}
		send(Mustache.to_html(templates.tail, data));
	});
}
