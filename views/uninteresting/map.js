function(doc) {
    if (doc.state !== 'fave') {
        emit(doc.updated, doc._rev);
    }
}