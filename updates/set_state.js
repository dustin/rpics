function(doc, req) {
    doc['state'] = req.form.new_state;
    return [doc, 'State set to ' + doc['state']];
}
