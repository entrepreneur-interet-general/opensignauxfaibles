function finalize(k, o) {

    batches.reduce((m, batch) => {
        var deleteOld = new Set(completeTypes[batch])
        o.batch[batch] = (o.batch[batch]||{})
        o.batch[batch].compact = (o.batch[batch].compact||{})
        o.batch[batch].compact["status"] = (o.batch[batch].compact["status"]||false)

        types.map(type => {
            o.batch[batch][type] = (o.batch[batch][type]||{})
            m[type] = (m[type] || new Set())

            var keys = Object.keys(o.batch[batch][type])
           
            o.batch[batch].compact.delete = (o.batch[batch].compact.delete||{})
            if (deleteOld.has(type) && o.batch[batch].compact.status == false) {
                var discardKeys = [...m[type]].filter(key => !(new Set(keys).has(key)))
                o.batch[batch].compact.delete[type] = discardKeys;

                discardKeys.forEach(key => {
                    m[type].delete(key)
                })
            }

            keys.filter(key => (m[type].has(key))).forEach(key => delete o.batch[batch][type][key])
            m[type] = new Set([...m[type]].concat(keys))
            if (Object.keys(o.batch[batch][type]).length == 0) {delete o.batch[batch][type]}
        })

        o.batch[batch].compact = (o.batch[batch].compact||{})
        o.batch[batch].compact.status = true
        return m
    }, {})

    o.index = {"algo1": false,
               "algo2":false}
    Object.keys(o.batch).forEach(batch => {
        Object.keys((o.batch[batch].effectif||{})).forEach(effectif => {
            o.index.algo1 = true
            o.index.algo2 = true
        })      
    })

    return o
}