function finalize(k, v) {

  v = Object.keys((v.batch || {})).sort().reduce((m, batch) => {
    Object.keys(v.batch[batch]).forEach(type => {
      m[type] = (m[type] || {})
      Object.assign(m[type], v.batch[batch][type])
    })
    return m
  }, { "siren": v.siren })

  var liste_periodes = generatePeriodSerie(date_debut, date_fin)

  var output_array = liste_periodes.map(function (e) {
    return {
      "siren": v.siren,
      "periode": e,
      "exercice_bdf": 0,
      "arrete_bilan_bdf": new Date(0),
      "exercice_diane": 0,
      "arrete_bilan_diane": new Date(0)
    }
  })

  var output_indexed = output_array.reduce(function (periode, val) {
    periode[val.periode.getTime()] = val
    return periode
  }, {})

  v.bdf = (v.bdf || {})
  v.diane = (v.diane || {})

  Object.keys(v.bdf).forEach(hash => {
    let periode_arrete_bilan = new Date(Date.UTC(v.bdf[hash].arrete_bilan_bdf.getUTCFullYear(), v.bdf[hash].arrete_bilan_bdf.getUTCMonth() +1, 1, 0, 0, 0, 0))
    let periode_dispo = DateAddMonth(periode_arrete_bilan, 8)
    let series = generatePeriodSerie(
      periode_dispo,
      DateAddMonth(periode_dispo, 13)
    )

    series.forEach(periode => {
      Object.keys(v.bdf[hash]).filter( k => {
        var omit = ["raison_sociale","secteur", "siren"]
        return (v.bdf[hash][k] != null &&  !(omit.includes(k)))
      }).forEach(k => {
        if (periode.getTime() in output_indexed){
          output_indexed[periode.getTime()][k] = v.bdf[hash][k]
          output_indexed[periode.getTime()].exercice_bdf = output_indexed[periode.getTime()].annee_bdf - 1
        }

        let past_year_offset = [1,2]
        past_year_offset.forEach( offset =>{
          let periode_offset = DateAddMonth(periode, 12* offset)
          let variable_name =  k + "_past_" + offset
          if (periode_offset.getTime() in output_indexed && 
            k != "arrete_bilan_bdf" &&
            k != "exercice_bdf"){
            output_indexed[periode_offset.getTime()][variable_name] = v.bdf[hash][k]  
          }
        })
      }
      )
    })
  })

  Object.keys(v.diane).forEach(hash => {

    //v.diane[hash].arrete_bilan_diane = new Date(Date.UTC(v.diane[hash].exercice_diane, 11, 31, 0, 0, 0, 0))
    let periode_arrete_bilan = new Date(Date.UTC(v.diane[hash].arrete_bilan_diane.getUTCFullYear(), v.diane[hash].arrete_bilan_diane.getUTCMonth() +1, 1, 0, 0, 0, 0))
    let periode_dispo = DateAddMonth(periode_arrete_bilan, 8)
    let series = generatePeriodSerie(
      periode_dispo,
      DateAddMonth(periode_dispo, 13)
    )

    series.forEach(periode => {
      Object.keys(v.diane[hash]).filter( k => {
        var omit = ["marquee", "nom_entreprise","numero_siren",
          "statut_juridique", "procedure_collective"]
        return (v.diane[hash][k] != null &&  !(omit.includes(k)))
      }).forEach(k => {       
        if (periode.getTime() in output_indexed){
          output_indexed[periode.getTime()][k] = v.diane[hash][k]
        }
        let past_year_offset = [1,2]
        past_year_offset.forEach(offset =>{
          let periode_offset = DateAddMonth(periode, 12 * offset)
          let variable_name =  k + "_past_" + offset

          if (periode_offset.getTime() in output_indexed && 
            k != "arrete_bilan_diane" &&
            k != "exercice_diane"){
            output_indexed[periode_offset.getTime()][variable_name] = v.diane[hash][k]
          }
        })
      }                   
      )           
    })


    //    series.forEach(periode => {
    //      if (periode.getTime() in output_indexed){
    //
    //        var EBE = (output_indexed[periode].valeur_ajoutee - output_indexed[periode].charges_sociales - output_indexed[periode].salaires_et_traitements)
    //        var achats_ht =  output_indexed[periode].marchandises + output_indexed[periode].matieres_prem_approv
    //
    //        if ("taux_marge" in output_indexed[periode]){
    //          output_indexed[periode].taux_marge = EBE / output_indexed[periode].valeur_ajoutee
    //          print("output_indexed[periode].taux_marge") 
    //        }
    //        if ("financier_court_terme" in output_indexed[periode]){
    //
    //          output_indexed[periode].financier_court_terme = output_indexed[periode].concours_bancaire_courant / output_indexed[periode].CA
    //          print("output_indexed[periode].financier_court_terme") 
    //        }
    //        if ("delai_fournisseur" in output_indexed[periode]){
    //          output_indexed[periode].delai_fournisseur = output_indexed[periode].dettes_fourn_et_cptes_ratt / achats_ht
    //          print("output_indexed[periode].delai_fournisseur") 
    //        }
    //        if ("dette_fiscale" in output_indexed[periode]){
    //          output_indexed[periode].dette_fiscale = output_indexed[periode].dettes_fiscales_et_sociales / output_indexed[periode].valeur_ajoutee
    //          print("output_indexed[periode].dette_fiscale") 
    //        }
    //        if ("frais_financier" in output_indexed[periode]){
    //          output_indexed[periode].frais_financier = output_indexed[periode].total_des_charges_fin / (EBE + output_indexed[periode].total_des_produits_fin - output_indexed[periode].total_des_charges_fin) 
    //          print("output_indexed[periode].frais_financier") 
    //        }
    //        if ("poids_frng" in output_indexed[periode]){
    //          output_indexed[periode].poids_frng = output_indexed[periode].fonds_de_roul_net_global / output_indexed[periode].CA
    //          print("output_indexed[periode].poids_frng")
    //        }
    //      }
    //    })
  })


  output_array.forEach((periode, index) => {
    if ((periode.arrete_bilan_bdf||new Date(0)).getTime() == 0 && (periode.arrete_bilan_diane || new Date(0)).getTime() == 0) {
      delete output_array[index]
    }
    if ((periode.arrete_bilan_bdf||new Date(0)).getTime() == 0){
      delete periode.arrete_bilan_bdf
    }
    if ((periode.arrete_bilan_diane||new Date(0)).getTime() == 0){
      delete periode.arrete_bilan_diane
    }
  })

  return {"siren": k, "entreprise": output_array}
}
