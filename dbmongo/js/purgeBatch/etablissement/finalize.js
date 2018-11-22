function finalize(k, o) {
    delete o.batch[currentBatch]
    return o
}