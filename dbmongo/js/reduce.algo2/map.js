function map () {
  let v = f.flatten(this.value, actual_batch)

  if (v.scope == "etablissement") {
    let o = f.outputs(v, serie_periode)
    let output_array = o[0]
    let output_indexed = o[1]
    let periodes = Object.keys(output_indexed) 
    if (v.debit) { f.debits(v) }
    if (v.effectif) {
      let output_effectif = f.effectifs(v, output_indexed)
      f.add(output_effectif, output_indexed)
    }

    if (v.interim){
      let output_interim = f.interim(v.interim, output_indexed)
      f.add(output_interim, output_indexed)
    }
    //  Object.keys(output_interim).forEach(periode => {
    //    output_indexed[periode] = Object.assign(output_indexed[periode], output_interim[periode])
    //  }) 
    //}
    
    if (v.apconso && v.apdemande) {f.apart(v, output_indexed, output_array)}
    if (v.delai) {f.delais(v, output_indexed)}

    v.altares = v.altares || {}
    v.procol = v.procol || {}

    if (v.altares) {f.defaillances(v, output_indexed)}
    if (v.cotisation && v.debit) {f.cotisationsdettes(v, output_array, output_indexed)}

    if (v.ccsf) {f.ccsf(v, output_array)}
    if (v.sirene) {f.sirene(v, output_array)}

    f.naf(output_indexed, naf)

    f.cotisation(output_indexed, output_array)

    let output_cible = f.cibleApprentissage(output_indexed)
    f.add(output_cible, output_indexed)


    output_array.forEach(val => {
      let data = {}
      data[this._id] = val
      emit(    
        { 'siren': this._id.substring(0, 9),
          'batch': actual_batch,
          'periode': val.periode},
          data
      )
    })
  }

  if (v.scope == "entreprise") {
    var output_array = serie_periode.map(function (e) {
      return {
        "siren": v.key,
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
      let periode_dispo = f.dateAddMonth(periode_arrete_bilan, 8)
      let series = f.generatePeriodSerie(
        periode_dispo,
        f.dateAddMonth(periode_dispo, 13)
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
        })
      })
    })
  
    Object.keys(v.diane).forEach(hash => {
      //v.diane[hash].arrete_bilan_diane = new Date(Date.UTC(v.diane[hash].exercice_diane, 11, 31, 0, 0, 0, 0))
      let periode_arrete_bilan = new Date(Date.UTC(v.diane[hash].arrete_bilan_diane.getUTCFullYear(), v.diane[hash].arrete_bilan_diane.getUTCMonth() +1, 1, 0, 0, 0, 0))
      let periode_dispo = DateAddMonth(periode_arrete_bilan, 8) // 01/09 pour un bilan le 31/12
      let series = f.generatePeriodSerie(
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

      series.forEach(periode => {
        if (periode.getTime() in output_indexed){
          // Recalcul BdF si ratios bdf sont absents
          if (!("taux_marge" in output_indexed[periode.getTime()]) && (f.tauxMarge(v.diane[hash]) !== null)){
            output_indexed[periode.getTime()].taux_marge = f.tauxMarge(v.diane[hash])
          }
          if (!("financier_court_terme" in output_indexed[periode.getTime()]) && (f.financierCourtTerme(v.diane[hash]) !== null)){
            output_indexed[periode.getTime()].financier_court_terme = f.financierCourtTerme(v.diane[hash])
          }
          if (!("poids_frng" in output_indexed[periode.getTime()]) && (f.poidsFrng(v.diane[hash]) !== null)){
            output_indexed[periode.getTime()].poids_frng = f.poidsFrng(v.diane[hash])
          }
          if (!("dette_fiscale" in output_indexed[periode.getTime()]) && (f.detteFiscale(v.diane[hash]) !== null)){
            output_indexed[periode.getTime()].dette_fiscale = f.detteFiscale(v.diane[hash])
          }
          if (!("frais_financier" in output_indexed[periode.getTime()]) && (f.fraisFinancier(v.diane[hash]) !== null)){
            output_indexed[periode.getTime()].frais_financier = f.fraisFinancier(v.diane[hash])
          }

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
      emit(    
        { 'siren': this._id.substring(0, 9),
          'batch': actual_batch,
          'periode': periode.periode},
          {'entreprise': periode}
        )
    })
  }
}
