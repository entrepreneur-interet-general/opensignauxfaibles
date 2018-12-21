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
    let periode_dispo = DateAddMonth(periode_arrete_bilan, 8) // 01/09 pour un bilan le 31/12
    let series = generatePeriodSerie(
      periode_dispo,
      DateAddMonth(periode_dispo, 13) // periode de validité d'un bilan auprès de la Banque de France: 21 mois (13+8)
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

        // Passé

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
    
    // FONCTIONS POUR CALCULER LES RATIOS BDF
    function taux_marge(diane) {
      if  (("excedent_brut_d_exploitation" in diane) && (diane["excedent_brut_d_exploitation"] !== null) &&
        ("valeur_ajoutee" in diane) && (diane["valeur_ajoutee"] !== null) &&
        (diane["excedent_brut_d_exploitation"] != 0)){
        return diane["excedent_brut_d_exploitation"]/diane["valeur_ajoutee"] * 100
      } else {
        return null
      }
    }
    function financier_court_terme(diane) {
      if  (("concours_bancaire_courant" in diane) && (diane["concours_bancaire_courant"] !== null) &&
        ("ca" in diane) && (diane["ca"] !== null) &&
        (diane["ca"] != 0)){
        return diane["concours_bancaire_courant"]/diane["ca"] * 100
      } else {
        return null
      }
    }
    function poids_frng(diane){
      if  (("couverture_ca_fdr" in diane) && (diane["couverture_ca_fdr"] !== null)){
        return diane["couverture_ca_fdr"]/360 * 100
      } else {
        return null
      }
    }
    function dette_fiscale(diane){
      if  (("dette_fiscale_et_sociale" in diane) && (diane["dette_fiscale_et_sociale"] !== null) &&
          ("valeur_ajoutee" in diane) && (diane["valeur_ajoutee"] !== null) &&
          (diane["valeur_ajoutee"] != 0)){
        return diane["dette_fiscale_et_sociale"]/ diane["valeur_ajoutee"] * 100
      } else {
        return null
      }
    }
    function frais_financier(diane){
      if  (("interets" in diane) && (diane["interets"] !== null) &&
        ("excedent_brut_d_exploitation" in diane) && (diane["excedent_brut_d_exploitation"] !== null) &&
        ("produits_financiers" in diane) && (diane["produits_financiers"] !== null) &&
        ("charges_financieres" in diane) && (diane["charges_financieres"] !== null) &&
        ("charge_exceptionnelle" in  diane) && (diane["charge_exceptionnelle"] !== null) &&
        ("produit_exceptionnel" in diane) && (diane["produit_exceptionnel"] !== null) &&
        diane["excedent_brut_d_exploitation"] + diane["produits_financiers"] + diane["produit_exceptionnel"] - diane["charge_exceptionnelle"] - diane["charges_financieres"] != 0 ){
        return diane["interets"] / (diane["excedent_brut_d_exploitation"] + diane["produits_financiers"] + diane["produit_exceptionnel"] -
          diane["charge_exceptionnelle"] - diane["charges_financieres"] ) * 100
      } else {
        return null
      }
    }

    // function delai_fournisseur(diane){
    //   if (("dette_fournisseur" in diane) &&("achat_marchandises" in diane) && ("achat_matieres_premieres" in diane) && ("autres_achats_charges_externes" in diane) &&
    //       (diane["achat_marchandises"] + diane["achat_matieres_premieres"] + diane["autres_achats_charges_externes"] != 0)){
    //     return diane["dette_fournisseur"] * 360 / (diane["achat_marchandises"] + diane["achat_matieres_premieres"] + diane["autres_achats_charges_externes"])
    //   } else {
    //     return null
    //   }
    // }

    series.forEach(periode => {
      if (periode.getTime() in output_indexed){
        // Recalcul BdF si ratios bdf sont absents
        if (!("taux_marge" in output_indexed[periode.getTime()]) && (taux_marge(v.diane[hash]) !== null)){
          output_indexed[periode.getTime()].taux_marge = taux_marge(v.diane[hash])
        }
        if (!("financier_court_terme" in output_indexed[periode.getTime()]) && (financier_court_terme(v.diane[hash]) !== null)){
          output_indexed[periode.getTime()].financier_court_terme = financier_court_terme(v.diane[hash])
        }
        if (!("poids_frng" in output_indexed[periode.getTime()]) && (poids_frng(v.diane[hash]) !== null)){
          output_indexed[periode.getTime()].poids_frng = poids_frng(v.diane[hash])
        }
        if (!("dette_fiscale" in output_indexed[periode.getTime()]) && (dette_fiscale(v.diane[hash]) !== null)){
          output_indexed[periode.getTime()].dette_fiscale = dette_fiscale(v.diane[hash])
        }
        if (!("frais_financier" in output_indexed[periode.getTime()]) && (frais_financier(v.diane[hash]) !== null)){
          output_indexed[periode.getTime()].frais_financier = frais_financier(v.diane[hash])
        }
        // if (!("delai_fournisseur" in output_indexed[periode.getTime()]) && (delai_fournisseur(v.diane[hash]) !== null)){
        //   output_indexed[periode.getTime()].delai_fournisseur = delai_fournisseur(v.diane[hash])
        // }
        var bdf_vars = ["taux_marge", "poids_frng", "dette_fiscale", "financier_court_terme", "frais_financier"]
        let past_year_offset = [1,2]
        bdf_vars.forEach(k =>{
          if (k in output_indexed[periode.getTime()]){
            past_year_offset.forEach(offset =>{
              let periode_offset = DateAddMonth(periode, 12 * offset)
              let variable_name =  k + "_past_" + offset

              if (periode_offset.getTime() in output_indexed){
                output_indexed[periode_offset.getTime()][variable_name] = output_indexed[periode.getTime()][k]
              }
            })
          }
        })
      }
    })


    //        var EBE = (output_indexed[periode].valeur_ajoutee - output_indexed[periode].charges_sociales - output_indexed[periode].salaires_et_traitements)
    //        var achats_ht =  output_indexed[periode].marchandises + output_indexed[periode].matieres_prem_approv
    //
    //        if ("taux_marge" in output_indexed[periode]){
    //          output_indexed[periode].taux_marge = EBE / output_indexed[periode].valeur_ajoutee
    //        }
    //        if ("financier_court_terme" in output_indexed[periode]){
    //
    //          output_indexed[periode].financier_court_terme = output_indexed[periode].concours_bancaire_courant / output_indexed[periode].CA
    //        }
    //        if ("delai_fournisseur" in output_indexed[periode]){
    //          output_indexed[periode].delai_fournisseur = output_indexed[periode].dettes_fourn_et_cptes_ratt / achats_ht
    //        }
    //        if ("dette_fiscale" in output_indexed[periode]){
    //          output_indexed[periode].dette_fiscale = output_indexed[periode].dettes_fiscales_et_sociales / output_indexed[periode].valeur_ajoutee
    //        }
    //        if ("frais_financier" in output_indexed[periode]){
    //          output_indexed[periode].frais_financier = output_indexed[periode].total_des_charges_fin / (EBE + output_indexed[periode].total_des_produits_fin - output_indexed[periode].total_des_charges_fin) 
    //        }
    //        if ("poids_frng" in output_indexed[periode]){
    //          output_indexed[periode].poids_frng = output_indexed[periode].fonds_de_roul_net_global / output_indexed[periode].CA
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
