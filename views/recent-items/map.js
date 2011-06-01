function(doc) {
    if (doc.updated) {
        emit(doc.updated, null);
    }
};
