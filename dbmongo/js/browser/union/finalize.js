function finalize(k, v) {

    // indexation periode entreprise
    var entreprise = (v.entreprise||[]).reduce((m, siren_periode) => {
        if (siren_periode) {
            m[siren_periode.periode.getTime()] = siren_periode
        }
        return m
    }, {})

/*     empty = {
        "arrete_bilan": null,
        "secteur": null,
        "poids_frng": null,
        "taux_marge": null,
        "delai_fournisseur": null,
        "dette_fiscale": null,
        "financier_court_terme": null,
        "frais_financier": null
    }
 */
    empty = (v.entreprise||[]).reduce((accu,siren_periode)=> {
        if (siren_periode){
        Object.keys(siren_periode).forEach(key => {
            accu[key] = null
        })
        }
        return(accu)
    },{})

    //
    ///
    ///////////////////////////////////////////////
    // Consolidation a l'echelle de l'entreprise //
    ///////////////////////////////////////////////
    ///
    //

    var etablissements_connus = []

    Object.keys(v).forEach(siret =>{
    
        if (siret != "entreprise" && siret != "siren" ) {
            etablissements_connus[siret] = true;
            
            v[siret].forEach(siret_periode => {
                if (siret_periode){
                    time = siret_periode.periode.getTime()
                    entreprise[time] = (entreprise[time] || {})
                    entreprise[time].effectif_entreprise = (entreprise[time].effectif_entreprise || 0) + siret_periode.effectif
                    entreprise[time].apart_entreprise = (entreprise[time].apart_entreprise || 0)  + siret_periode.apart_heures_consommees
                    entreprise[time].debit_entreprise = (entreprise[time].debit_entreprise || 0) + siret_periode.montant_part_patronale + siret_periode.montant_part_ouvriere   

                }
            })
        }
    })

    Object.keys(entreprise).forEach(k => {
        entreprise[k].nbr_etablissements_connus = Object.keys(etablissements_connus).length
    })

    
    // reduce sur tous les établissements
    output = Object.keys(v).reduce((accu, k) => {
        if (k != "entreprise" && k != "siren") {

            v[k].forEach(siret_periode => {
                if (siret_periode) {
                    Object.assign(siret_periode, (entreprise[(siret_periode || { "periode": new Date(0) }).periode.getTime()] || empty))
                    accu.push(siret_periode)
                }
            })

        }
        return accu

    }, [])

    return output
}