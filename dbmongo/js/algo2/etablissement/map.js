function map() {
    v = this.value

    v = Object.keys((v.batch || {})).sort().filter(batch => batch <= actual_batch).reduce((m, batch) => {
        Object.keys(v.batch[batch]).forEach((type) => {
            m[type] = (m[type] || {})
            var  array_delete = (v.batch[batch].compact.delete[type]||[])
            if (array_delete != {}) {array_delete.forEach(hash => {
                delete m[type][hash]
            })
            }
            Object.assign(m[type], v.batch[batch][type])
        })
        return m
    }, { "siret": this.value.siret })

    // filtre des établissements qui n'ont pas dépassé 20 employés dans les 6 derniers mois
    // FIX ME: la présence ou non des etablissements dans l'historique depend des derniers mois. 
    var filter = (Object.keys(v.effectif).sort(
        (eff1,eff2) => v.effectif[eff1].periode <= v.effectif[eff2].periode
    ).slice(0,6).filter(
        eff => v.effectif[eff].effectif >=20
    ) || []).length == 0

    if( !filter ){
        emit(this.value.siret, v)
    }

}
