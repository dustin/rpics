function(doc) {
    if (doc.state === 'fave') {
        emit(doc.updated, null);
    }
}