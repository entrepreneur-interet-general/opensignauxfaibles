function finalize(k, v) {

  var offset_effectif = (date_fin_effectif.getUTCFullYear() - date_fin.getUTCFullYear()) * 12 + date_fin_effectif.getUTCMonth() - date_fin.getUTCMonth()
  var offset_cotisation = 1

  var liste_periodes = generatePeriodSerie(date_debut, date_fin)
 
  v = Object.keys((v.batch || {})).sort().filter(batch => batch <= actual_batch).reduce((m, batch) => {
    Object.keys(v.batch[batch]).forEach(type => {
      m[type] = (m[type] || {})
      var  array_delete = (v.batch[batch].compact.delete[type]||[])
      if (array_delete != {}) {array_delete.forEach(hash => {
        delete m[type][hash]
      })
      }
      Object.assign(m[type], v.batch[batch][type])
    })
    return m
  }, { "siret": k })


  v.apconso = (v.apconso || {})
  v.apdemande = (v.apdemande || {})
  v.effectif = (v.effectif || {})
  v.altares = (v.altares || {})
  v.procol = (v.procol || {})
  v.cotisation = (v.cotisation || {})
  v.debit = (v.debit || {})
  v.delai = (v.delai || {})


  // relier les débits
  var ecn = Object.keys(v.debit).reduce((m, h) => {
    var d = [h, v.debit[h]]
    var start = d[1].periode.start
    var end = d[1].periode.end
    var num_ecn = d[1].numero_ecart_negatif
    var compte = d[1].numero_compte
    var key = start + "-" + end + "-" + num_ecn + "-" + compte
    m[key] = (m[key] || []).concat([{
      "hash": d[0],
      "numero_historique": d[1].numero_historique,
      "date_traitement": d[1].date_traitement
    }])
    return m
  }, {})

  Object.keys(ecn).forEach(i => {
    ecn[i].sort(compareDebit)
    var l = ecn[i].length
    ecn[i].forEach((e, idx) => {
      if (idx <= l - 2) {
        v.debit[e.hash].debit_suivant = ecn[i][idx + 1].hash
      }
    })
  })

  var output_array = liste_periodes.map(function (e) {
    return {
      "siret": v.siret,
      "periode": e,
      //"lag_effectif_missing": true,
      "effectif": null,
      "date_effectif": null,
      "apart_heures_consommees": 0,
      "apart_heures_autorisees": 0,
      "apart_motif_recours": 0,
      //"effectif_history": {},
      //"outcome_0_12": "non_default",
      //"date_defaillance": null,
      //"cotisation_due_periode": {},
      "debit_array": [],
      "etat_proc_collective": "in_bonis",
      "interessante_urssaf": true
    }
  })

  var output_indexed = output_array.reduce(function (periode, val) {
    periode[val.periode.getTime()] = val
    return periode
  }, {})



  //
  ///
  ///////////////
  // Effectifs //
  ///////////////
  ///
  //
  var map_effectif = Object.keys(v.effectif).reduce(function (map_effectif, hash) {
    var effectif = v.effectif[hash]
    if (effectif == null) {
      return map_effectif
    }
    var effectifTime = effectif.periode.getTime()
    map_effectif[effectifTime] = (map_effectif[effectifTime] || 0) + effectif.effectif
    return map_effectif
  }, {})


  Object.keys(map_effectif).forEach(time =>{
    var time_d = new Date(parseInt(time))
    var time_offset = DateAddMonth(time_d, -offset_effectif -1)
    if (time_offset.getTime() in output_indexed){
      output_indexed[time_offset.getTime()].effectif = map_effectif[time]
      output_indexed[time_offset.getTime()].date_effectif = time_d
    }

    var past_month_offsets = [6,12,18,24]
    past_month_offsets.forEach( offset =>{
      var time_past_offset = DateAddMonth(time_offset, offset)
      var variable_name_effectif = "effectif_past_" + offset
      if (time_past_offset.getTime() in output_indexed){
        var val_offset = output_indexed[time_past_offset.getTime()]
        val_offset[variable_name_effectif] = map_effectif[time]
      }
    }
    )
  }
  )
  
  output_array.forEach(function (val, index) {
    if (val.effectif == null) {
      delete output_indexed[val.periode.getTime()]
      delete output_array[index]
    }
  })


  //
  ///
  ////////////////////////
  // activite partielle // 
  //////////////////////// 
  ///
  //

  var apart = Object.keys(v.apdemande).reduce((apart, hash) => {
    apart[v.apdemande[hash].id_demande.substring(0, 9)] = {
      "demande": hash,
      "consommation": []
    }
    return apart
  }, {})

  Object.keys(v.apdemande).forEach(hash => {
    var periode_deb = v.apdemande[hash].periode.start
    var periode_fin = v.apdemande[hash].periode.end
    var periode_deb_floor = new Date(Date.UTC(periode_deb.getUTCFullYear(), periode_deb.getUTCMonth(), 1, 0, 0, 0, 0))
    var periode_fin_ceil = new Date(Date.UTC(periode_fin.getUTCFullYear(), periode_fin.getUTCMonth() + 1, 1, 0, 0, 0, 0))
    var series = generatePeriodSerie(periode_deb_floor, periode_fin_ceil)
    series.forEach( date => {
      let time = date.getTime()
      if (time in output_indexed){
        output_indexed[time].apart_heures_autorisees = v.apdemande[hash].hta
      }  
    })
  })



  Object.keys(v.apconso).forEach(hash => {
    var valueap = v.apconso[hash]
    if (valueap.id_conso.substring(0, 9) in apart) {
      apart[valueap.id_conso.substring(0, 9)].consommation.push(hash)
    }
  })

  // relier apdemande et apconso
  Object.keys(apart).forEach(k => {
    v.apdemande[apart[k].demande].hash_consommation = apart[k].consommation
    for (let j in apart[k].consommation) {
      v.apconso[apart[k].consommation[j]].hash_demande = apart[k].demande
    }
  })

  Object.keys(v.apconso).forEach(
    function (h) {
      var conso = v.apconso[h]
      if (conso.hash_demande) {
        var time = conso.periode.getTime()
        if (time in output_indexed){
          output_indexed[time].apart_heures_consommees = output_indexed[time].apart_heures_consommees + conso.heure_consomme
          output_indexed[time].apart_motif_recours = v.apdemande[conso.hash_demande].motif_recours_se
        }
      }
    })

  output_array.forEach(val => {
    if (val !== null && val.effectif > 0){
      val.ratio_apart = val.apart_heures_consommees / (val.effectif * 157.67)
    }
  })


  // ratio_dette_delai = base::ifelse(
  //   !is.na(duree_delai) & duree_delai > 0,
  //   (
  //     montant_part_patronale + montant_part_ouvriere -
  //       montant_echeancier * delai / (duree_delai / 30)
  //   ) / montant_echeancier,
  //   NA
  // )
  
  //
  ///
  //////////////////
  // defaillance  //
  //////////////////
  ///
  //
  function deal_with_procols(data_source, altar_or_procol){
    var codes  =  Object.keys(data_source).reduce(function(events,hash) {
      var the_event = data_source[hash]

      if (altar_or_procol == "altares")
        var etat = altaresToHuman(the_event.code_evenement);
      else if (altar_or_procol == "procol")
        var etat = procolToHuman(the_event.action_procol, the_event.stade_procol);

      if (etat != null)
        events.push({"etat": etat, "date_proc_col": new Date(the_event.date_effet)})

      return(events)
    },[]).sort(
      function(a,b) {return(a.date_proc_col.getTime() > b.date_proc_col.getTime())}
    )

    codes.forEach(
      function (event) {
        let periode_effet = new Date(Date.UTC(event.date_proc_col.getFullYear(), event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0))
        var time_til_last = Object.keys(output_indexed).filter(val => {return (val >= periode_effet)})

        time_til_last.forEach(time => {
          if (time in output_indexed) {
            output_indexed[time].etat_proc_collective = event.etat
            output_indexed[time].date_proc_collective = event.date_proc_col
            if (event.etat != "in_bonis")
              output_indexed[time].tag_outcome = "failure"
          }
        })
      }
    )

    //output_array.forEach(periode => {
    //  if ((periode.date_proc_collective || new Date(0)).getTime() == 0){
    //    delete periode.date_proc_collective
    //  }
    //})
  }


  deal_with_procols(v.altares, "altares")
  deal_with_procols(v.procol, "procol")
  // // On filtre altares pour ne garder que les codes qui nous intéressents
  // var altares_codes  =  Object.keys(v.altares).reduce(function(events,hash) {
  //   var altares_event = v.altares[hash]

  //   var etat = altaresToHuman(altares_event.code_evenement)

  //   if (etat != null)
  //     events.push({"etat": etat, "date_proc_col": new Date(altares_event.date_effet)})

  //   return(events)
  // },[{"etat" : "in_bonis", "date_proc_col" : new Date(0)}]).sort(
  //   function(a,b) {return(a.date_proc_col.getTime() > b.date_proc_col.getTime())}
  // )

  // altares_codes.forEach(
  //   function (event) {
  //     var periode_effet = new Date(Date.UTC(event.date_proc_col.getFullYear(), event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0))
  //     var time_til_last = Object.keys(output_indexed).filter(val => {return (val >= periode_effet)})
  //     time_til_last.forEach(time => {
  //       if (time in output_indexed) {
  //         output_indexed[time].etat_proc_collective = event.etat
  //         output_indexed[time].date_proc_collective = event.date_proc_col
  //       }
  //     })
  //   }
  // )

  // output_array.forEach(periode => {
  //   if ((periode.date_proc_collective || new Date(0)).getTime() == 0){
  //     delete periode.date_proc_collective
  //   }
  // })

  // /// Procédures collectives
  // let procol_codes  =  Object.keys(v.procol).reduce(function(events,hash) {
  //   var procol_event = v.procol[hash]

  //   var etat = procolToHuman(procol_event.action_procol, procol_event.stade_procol)

  //   if (etat != null)
  //     events.push({"etat": etat, "date_proc_col": new Date(procol_event.date_effet)})

  //   return(events)
  // },[]).sort(
  //   function(a,b) {return(a.date_proc_col.getTime() > b.date_proc_col.getTime())}
  // )

  // procol_codes.forEach( my_event => {
  //   let periode_effet = new Date(Date.UTC(my_event.date_proc_col.getFullYear(), my_event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0))
  //   let time_til_last = Object.keys(output_indexed).filter(val => {return (val >= periode_effet)})

  //   time_til_last.forEach(time => {
  //     if (time in output_indexed) {
  //       output_indexed[time].etat_proc_collective = my_event.etat
  //       output_indexed[time].date_proc_collective = my_event.date_proc_col
  //     }
  //   })
  // })



  //
  ///
  ////////////////////////////
  // Cotisation et débits  ///
  ////////////////////////////
  ///
  //

  var value_cotisation = {}

  Object.keys(v.cotisation).forEach(function (h) {
    var cotisation = v.cotisation[h]
    var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
    periode_cotisation.forEach(date_cotisation => {
      let date_offset = DateAddMonth(date_cotisation, offset_cotisation)
      value_cotisation[date_offset.getTime()] = (value_cotisation[date_offset.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
    })
  })


  var value_dette = {}

  Object.keys(v.debit).forEach(function (h) {
    var debit = v.debit[h]

    var debit_suivant = (v.debit[debit.debit_suivant] || {"date_traitement" : date_fin})
    let date_limite = date_fin//new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))
    date_traitement_debut = new Date(
      Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth())
    )

    date_traitement_fin = new Date(
      Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth())
    )

    let periode_debut = (date_traitement_debut.getTime() >= date_limite.getTime() ? date_limite : date_traitement_debut)
    let periode_fin = (date_traitement_fin.getTime() >= date_limite.getTime() ? date_limite : date_traitement_fin)

    generatePeriodSerie(periode_debut, periode_fin).map(function (date) {
      let time = date.getTime()
      value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale, "montant_majorations": debit.montant_majorations}])
    })

  })    

  var numeros_compte = Array.from(new Set(
    Object.keys(v.cotisation).map(function (h) {
      return(v.cotisation[h].numero_compte)
    })))


  Object.keys(output_indexed).forEach(function (time) {
    output_indexed[time].numero_compte_urssaf = numeros_compte
    if (time in value_cotisation){
      output_indexed[time].cotisation = value_cotisation[time].reduce((a,cot) => a + cot,0)
    }

    if (time in value_dette) {
      output_indexed[time].debit_array = value_dette[time]
    }
  })

  Object.keys(output_indexed).forEach(time => {
    let time_d = new Date(parseInt(time))
    var val = output_indexed[time]

    val.montant_dette = val.debit_array.reduce(function (m, dette) {
      m.part_ouvriere += dette.part_ouvriere
      m.part_patronale += dette.part_patronale
      m.montant_majorations += dette.montant_majorations
      return m
    }, {"part_ouvriere": 0, "part_patronale": 0, "montant_majorations": 0})

    val.montant_part_ouvriere = val.montant_dette.part_ouvriere
    val.montant_part_patronale = val.montant_dette.part_patronale
    val.montant_majorations = val.montant_dette.montant_majorations

    let past_month_offsets = [1,2,3,6,12]
    past_month_offsets.forEach(offset => {
      let time_offset = DateAddMonth(time_d, offset)      
      let variable_name_part_ouvriere = "montant_part_ouvriere_past_" + offset
      let variable_name_part_patronale = "montant_part_patronale_past_" + offset
      if (time_offset.getTime() in output_indexed){
        let val_offset = output_indexed[time_offset.getTime()]
        val_offset[variable_name_part_ouvriere] = val.montant_part_ouvriere
        val_offset[variable_name_part_patronale] = val.montant_part_patronale
      }
    })

    let future_month_offsets = [0, 1, 2, 3, 4, 5]
    if (val.montant_part_ouvriere + val.montant_part_patronale > 0){
      future_month_offsets.forEach(offset => {
        let time_offset = DateAddMonth(time_d, offset)
        let val_offset = output_indexed[time_offset.getTime()]
        if (time_offset.getTime() in output_indexed){
          val_offset.interessante_urssaf = false    
        }
      })
    }
    delete val.montant_dette
    delete val.debit_array
  }
  )

  //
  ///
  ////////////
  // delais //
  ////////////
  ///
  //

  Object.keys(v.delai).map(
    function (hash) {
      var delai = v.delai[hash]
      var date_creation = new Date(Date.UTC(delai.date_creation.getUTCFullYear(), delai.date_creation.getUTCMonth(), 1, 0, 0, 0, 0))
      var date_echeance = new Date(Date.UTC(delai.date_echeance.getUTCFullYear(), delai.date_echeance.getUTCMonth(), 1, 0, 0, 0, 0))
      var pastYearTimes = generatePeriodSerie(date_creation, date_echeance).map(function (date) { return date.getTime() })
      pastYearTimes.map(
        function(time){
          if (time in output_indexed) {
            var remaining_months = (date_echeance.getUTCMonth() - new Date(time).getUTCMonth()) +
              12*(date_echeance.getUTCFullYear() - new Date(time).getUTCFullYear()) 
            output_indexed[time].delai = remaining_months
            output_indexed[time].duree_delai = delai.duree_delai
            output_indexed[time].montant_echeancier = delai.montant_echeancier

            if (delai.duree_delai > 0){
              output_indexed[time].ratio_dette_delai = (output_indexed[time].montant_part_patronale + 
                output_indexed[time].montant_part_ouvriere - delai.montant_echeancier * remaining_months * 30 / 
                (delai.duree_delai)) /
                delai.montant_echeancier
            }
          }
        }
      )
    }
  ) 

  //
  ///
  /////////
  // CCSF// 
  /////////
  ///
  //
  var ccsfHashes = Object.keys(v.ccsf || {}) 

  output_array.forEach(val => {        
    var optccsf = ccsfHashes.reduce( 
      function (accu, hash) { 
        var ccsf = v.ccsf[hash] 
        if (ccsf.date_traitement.getTime() < val.periode.getTime() && ccsf.date_traitement.getTime() > accu.date_traitement.getTime()) { 
          accu = ccsf 
        } 
        return(accu)
      }, 
      { 
        date_traitement: new Date(0) 
      } 
    )         

    if (optccsf.date_traitement.getTime() != 0) { 
      val.date_ccsf = optccsf.date_traitement 
    } 
  })


  //
  ///
  ////////////
  // Sirene //
  ////////////
  ///
  //

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
    val.tranche_ca = (sirene || {"trancheca":null}).trancheca
    val.indice_monoactivite = (sirene || {"indicemonoactivite": null}).indicemonoactivite  

  })


  //
  ///
  /////////
  // NAF //
  /////////
  ///
  //
  Object.keys(output_indexed).forEach(k =>{
    if (("code_ape" in output_indexed[k]) && (output_indexed[k].code_ape !== null)){
      output_indexed[k].code_naf = naf.n5to1[output_indexed[k].code_ape]
      output_indexed[k].libelle_naf = naf.n1[output_indexed[k].code_naf]
      output_indexed[k].libelle_ape2 = naf.n2[output_indexed[k].code_ape.substring(0,2)]
      output_indexed[k].libelle_ape3 = naf.n3[output_indexed[k].code_ape.substring(0,3)]
      output_indexed[k].libelle_ape4 = naf.n4[output_indexed[k].code_ape.substring(0,4)]
      output_indexed[k].libelle_ape5 = naf.n5[output_indexed[k].code_ape]
    }
  }
  )

  //
  ///
  /////////////////////
  // moyenne 12 mois //
  /////////////////////
  ///
  //

  // calcul de cotisation_moyenne sur 12 mois
  Object.keys(output_indexed).forEach(k => {
    let periode_courante = output_indexed[k].periode
    let periode_12_mois = DateAddMonth(periode_courante, 12)
    let series = generatePeriodSerie(periode_courante, periode_12_mois)
    series.forEach(periode => {
      if (periode.getTime() in output_indexed){
        if ("cotisation" in output_indexed[periode_courante.getTime()])
          output_indexed[periode.getTime()].cotisation_array = (output_indexed[periode.getTime()].cotisation_array || []).concat(output_indexed[periode_courante.getTime()].cotisation)

        output_indexed[periode.getTime()].montant_pp_array = 
          (output_indexed[periode.getTime()].montant_pp_array || []).concat( output_indexed[periode_courante.getTime()].montant_part_patronale)
        output_indexed[periode.getTime()].montant_po_array = 
          (output_indexed[periode.getTime()].montant_po_array || []).concat( output_indexed[periode_courante.getTime()].montant_part_ouvriere)
      }
    })
  }
  )

  output_array.forEach(val => {
    val.cotisation_array = (val.cotisation_array || [] )
    val.cotisation_moy12m = val.cotisation_array.reduce( (p, c) => p + c, 0) / (val.cotisation_array.length || 1) 
    if (val.cotisation_moy12m > 0) {
      val.ratio_dette = (val.montant_part_ouvriere + val.montant_part_patronale) / val.cotisation_moy12m
      let pp_average = (val.montant_pp_array || []).reduce((p, c) => p + c, 0) / (val.montant_pp_array.length || 1)
      let po_average =  (val.montant_po_array || []).reduce((p, c) => p + c, 0) / (val.montant_po_array.length || 1)
      val.ratio_dette_moy12m = (po_average + pp_average) / val.cotisation_moy12m
    }
    val.dette_any_12m = (val.montant_pp_array || []).reduce((p,c) => (c >= 100) || p, false) || (val.montant_po_array || []).reduce((p, c) => (c >= 100) || p, false)
    delete val.cotisation_array
    delete val.montant_pp_array
    delete val.montant_po_array
  })

  // Calcul des défauts URSSAF prolongés
  let counter = 0
  Object.keys(output_indexed).sort().forEach(k => {
    if (output_indexed[k].ratio_dette > 1){
      counter = counter + 1
      if (counter >= 3) 
        output_indexed[k].tag_outcome = "default"
    } else 
      counter = 0
  })

  //
  ///
  ///////////////////////////
  // Cible d'apprentissage //
  // ///////////////////////
  ///
  //
  counter = -1
  Object.keys(output_indexed).sort((a,b)=> a<=b).forEach( k => {
    if (counter >=0) counter = counter + 1 
    if (output_indexed[k].tag_outcome == "default" || output_indexed[k].tag_outcome == "failure"){
      counter = 0 
    }
    if (counter >= 0){
      output_indexed[k].time_til_outcome = counter
    }
  })

  //
  ///
  //////////////////////
  // Echantillon test //
  //////////////////////
  ///
  //


  var return_value = {"siren": k.substring(0, 9)}
  return_value[k] = output_array
  
  return return_value
}

