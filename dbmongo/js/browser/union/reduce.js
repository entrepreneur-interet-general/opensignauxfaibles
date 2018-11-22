function reduce(key, values) {
    if (values.length > 40){
        return {}
    }
    return values.reduce((m, v) => {
        Object.assign(m, v)
        return m
    }, {})
}