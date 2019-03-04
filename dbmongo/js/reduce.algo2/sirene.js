function sirene (v, output_array) {
  var sireneHashes = Object.keys(v.sirene || {})

  output_array.forEach(val => {
    // geolocalisation

    if (sireneHashes.length != 0) {
      var sirene = v.sirene[sireneHashes[0]]
    }
    val.siren = val.siret.substring(0, 9)
    val.lattitude = (sirene || { "lattitude": null }).lattitude
    val.longitude = (sirene || { "longitude": null }).longitude
    val.region = (sirene || {"region": null}).region
    val.departement = (sirene || {"departement": null}).departement
    val.code_ape  = (sirene || { "ape": null}).ape
    val.raison_sociale = (sirene || {"raisonsociale": null}).raisonsociale
    val.activite_saisonniere = (sirene || {"activitesaisoniere": null}).activitesaisoniere
    val.productif = (sirene || {"productif": null}).productif
    val.date_creation = (sirene || {"creation": null}).creation
    val.date_creation = val.date_creation !== null ? val.date_creation.getFullYear() : val.date_creation
    val.age = val.periode.getFullYear() - val.date_creation
    val.tranche_ca = (sirene || {"tranche_ca":null}).tranche_ca
    val.indice_monoactivite = (sirene || {"indicemonoactivite": null}).indicemonoactivite  
  })
}
