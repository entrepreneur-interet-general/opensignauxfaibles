function finalize(k, v) {

  v = Object.keys((v.batch || {})).sort().reduce((m, batch) => {
    Object.keys(v.batch[batch]).forEach((type) => {
      m[type] = (m[type] || {})
      Object.assign(m[type], v.batch[batch][type])
    })
    return m
  }, { "siren": v.siren })

  liste_periodes = generatePeriodSerie(date_debut, date_fin)

  var output_array = liste_periodes.map(function (e) {
    return {
      "siren": v.siren,
      "periode": e,
      "exercice_bdf": 0,
      "arrete_bilan_bdf": new Date(0),
      "exercice_diane": 0,
      "arrete_bilan_diane": new Date(0)
    }
  });

  var output_indexed = output_array.reduce(function (periode, val) {
    periode[val.periode.getTime()] = val
    return periode
  }, {});

  v.bdf = (v.bdf || {})
  v.diane = (v.diane || {})

  Object.keys(v.bdf).forEach(hash => {
    let periode_arrete_bilan = new Date(Date.UTC(v.bdf[hash].arrete_bilan_bdf.getUTCFullYear(), v.bdf[hash].arrete_bilan_bdf.getUTCMonth() +1, 1, 0, 0, 0, 0));
    let series = generatePeriodSerie(
      periode_arrete_bilan,
      DateAddMonth(periode_arrete_bilan, 12)
    )

    series.forEach(periode => {
      Object.keys(v.bdf[hash]).filter( k => {
        var omit = ["raison_sociale","secteur", "siren"]
        return (v.bdf[hash][k] != null &&  !(omit.includes(k)))
      }).forEach(k => {
        if (periode.getTime() in output_indexed){
          output_indexed[periode.getTime()][k] = v.bdf[hash][k]
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

    // Tant que la date exacte du bilan n'est pas exportee, on part sur le 31 decembre de l'exercice courant. 
    v.diane[hash].arrete_bilan_diane = new Date(Date.UTC(v.diane[hash].exercice_diane, 12, 31, 0, 0, 0, 0))
    let periode_arrete_bilan = new Date(Date.UTC(v.diane[hash].arrete_bilan_diane.getUTCFullYear(), v.diane[hash].arrete_bilan_diane.getUTCMonth() +1, 1, 0, 0, 0, 0));

    let series = generatePeriodSerie(
      periode_arrete_bilan,
      DateAddMonth(periode_arrete_bilan, 12)
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
          let periode_offset = DateAddMonth(periode, 12* offset)
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
  })


  output_array.forEach((periode, index) => {
    if ((periode.arrete_bilan_bdf||new Date(0)).getTime() == 0 && (periode.arrete_bilan_diane || new Date(0)).getTime() == 0) {
      delete output_array[index]
    }
  })

  return {"siren": k, "entreprise": output_array}
}
