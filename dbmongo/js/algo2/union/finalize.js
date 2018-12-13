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
      if (v[siret]){
        //    var time = v[siret].periode.getTime()
        Object.assign(v[siret],entreprise) 
        v[siret].effectif_entreprise = (v[siret].effectif_entreprise || 0) + v[siret].effectif
        v[siret].apart_entreprise = (v[siret].apart_entreprise || 0)  + v[siret].apart_heures_consommees
        v[siret].debit_entreprise = (v[siret].debit_entreprise || 0) + v[siret].montant_part_patronale + v[siret].montant_part_ouvriere   
      }
    }
  })

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
