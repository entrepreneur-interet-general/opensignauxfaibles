function finalize(k, v) {


  //var empty = (v.entreprise||[]).reduce((accu, siren_periode) => {
  //  if (siren_periode){
  //    Object.keys(siren_periode).forEach(key => {
  //      accu[key] = null
  //    })
  //  }
  //  return(accu)
  //},{})
  //
  ///
  ///////////////////////////////////////////////
  // consolidation a l'echelle de l'entreprise //
  ///////////////////////////////////////////////
  ///
  //

  let etablissements_connus = []
  let entreprise = (v.entreprise || {})

  Object.keys(v).forEach(siret =>{
    if (siret != "entreprise" && siret != "siren" ) {
      etablissements_connus[siret] = true
      //if (v[siret]){  // always TRUE
      //    var time = v[siret].periode.getTime()
      entreprise.effectif_entreprise = (entreprise.effectif_entreprise || 0) + v[siret].effectif // initialized to null
      entreprise.apart_entreprise = (entreprise.apart_entreprise || 0) + v[siret].apart_heures_consommees // initialized to 0
      entreprise.debit_entreprise = (entreprise.debit_entreprise || 0) +
        (v[siret].montant_part_patronale || 0) +
        (v[siret].montant_part_ouvriere || 0) 
      // not initialized
      //}
    }
  })

  Object.keys(v).forEach(siret =>{ 
      Object.assign(v[siret], entreprise) 
  })

  //
  ///
  //////////////////////////////
  /// Objectif entrainement ///
  //////////////////////////////
  ///
  //


  //une fois que les comptes sont faits...
  let output = []
  Object.keys(v).forEach(siret =>{
    if (siret != "entreprise" && siret != "siren" && v[siret]) {
      v[siret].nbr_etablissements_connus = Object.keys(etablissements_connus).length
      output.push(v[siret])
    }
  })

  if (output.length > 0)
    return output
}
