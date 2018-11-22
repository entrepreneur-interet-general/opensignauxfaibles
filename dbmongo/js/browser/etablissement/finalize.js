function finalize(k, v) {
 
    // var offset_effectif = (date_fin_effectif.getUTCFullYear() - date_fin.getUTCFullYear()) * 12 + date_fin_effectif.getUTCMonth() - date_fin.getUTCMonth()
    // liste_periodes = generatePeriodSerie(date_debut, date_fin)
    
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
    }, { "siret": k })

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
                v.debit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
            }
        })
    })
    
    var value_cotisation = {}
    
    Object.keys(v.cotisation).forEach(function (h) {
        var cotisation = v.cotisation[h]
        var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
        periode_cotisation.forEach(function (date_cotisation) {
            value_cotisation[date_cotisation.getTime()] = (value_cotisation[date_cotisation.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
        })
    })
    
    var value_dette = {}
    
    liste_periodes = generatePeriodSerie(date_debut, date_fin)

    var output_array = liste_periodes.map(function (e) {
        return {
            "siret": v.siret,
            "periode": e,
            //"lag_effectif_missing": true,
            "effectif": null,
            "date_effectif": null,
            "apart_heures_consommees": 0,
            "apart_motif_recours": 0,
            //"effectif_history": {},
            //"outcome_0_12": "non_default",
            //"date_defaillance": null,
            //"cotisation_due_periode": {},
            "debit_array": [],
            "etat_proc_collective": "in_bonis"
        }
    });

    var output_indexed = output_array.reduce(function (periode, val) {
        periode[val.periode.getTime()] = val
        return periode
    }, {});

    Object.keys(v.debit).forEach(function (h) {
        var debit = v.debit[h]
        if (debit.part_ouvriere + debit.part_patronale > 0) {

            var debit_suivant = (v.debit[debit.debit_suivant] || {"date_traitement" : date_fin})
            date_limite = date_fin//new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))
            date_traitement_debut = new Date(
                Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth())
            )
            
            date_traitement_fin = new Date(
                Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth())
            )
            
            periode_debut = (date_traitement_debut.getTime() >= date_limite.getTime() ? date_limite : date_traitement_debut)
            periode_fin = (date_traitement_fin.getTime() >= date_limite.getTime() ? date_limite : date_traitement_fin)
            
            generatePeriodSerie(periode_debut, periode_fin).map(function (date) {
                time = date.getTime()
                value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale }])
            })
        }
    })    

    Object.keys(output_indexed).forEach(function (time) {
        if (time in value_cotisation){
            output_indexed[time].cotisation = value_cotisation[time].reduce((a,cot) => a + cot,0)
        }
        
        if (time in value_dette) {
            output_indexed[time].debit_array = value_dette[time]
        }
    })

    output_array.forEach(function (val) {
        
        val.montant_dette = val.debit_array.reduce(function (m, dette) {
            m.part_ouvriere += dette.part_ouvriere
            m.part_patronale += dette.part_patronale
            return m
        }, { "part_ouvriere": 0, "part_patronale": 0 })
        
        val.montant_part_ouvriere = val.montant_dette.part_ouvriere
        val.montant_part_patronale = val.montant_dette.part_patronale

        delete val.montant_dette
        delete val.debit_array
        
    })

    v = Object.keys(v).filter(k => k !== 'siret' && k !== 'compact').reduce((accu, type) => {
        accu[type] = Object.keys(v[type]).map(k => v[type][k])
        return accu
    }, {})

    v.array_debit = output_array 

    return v

//     v.apconso = (v.apconso || {})
//     v.apdemande = (v.apdemande || {})
//     v.effectif = (v.effectif || {})
//     v.altares = (v.altares || {})
//     v.cotisation = (v.cotisation || {})
//     v.debit = (v.debit || {})
//     v.delai = (v.delai || {})
    
    
//     // relier les débits
//     var ecn = Object.keys(v.debit).reduce((m, h) => {
//         var d = [h, v.debit[h]]
//         var start = d[1].periode.start
//         var end = d[1].periode.end
//         var num_ecn = d[1].numero_ecart_negatif
//         var compte = d[1].numero_compte
//         var key = start + "-" + end + "-" + num_ecn + "-" + compte
//         m[key] = (m[key] || []).concat([{
//             "hash": d[0],
//             "numero_historique": d[1].numero_historique,
//             "date_traitement": d[1].date_traitement
//         }])
//         return m
//     }, {})

//     Object.keys(ecn).forEach(i => {
//         ecn[i].sort(compareDebit)
//         var l = ecn[i].length
//         ecn[i].forEach((e, idx) => {
//             if (idx <= l - 2) {
//                 v.debit[e.hash].debit_suivant = ecn[i][idx + 1].hash;
//             }
//         })
//     })
    
//     var output_array = liste_periodes.map(function (e) {
//         return {
//             "siret": v.siret,
//             "periode": e,
//             //"lag_effectif_missing": true,
//             "effectif": null,
//             "date_effectif": null,
//             "apart_heures_consommees": 0,
//             "apart_motif_recours": 0,
//             //"effectif_history": {},
//             //"outcome_0_12": "non_default",
//             //"date_defaillance": null,
//             //"cotisation_due_periode": {},
//             "debit_array": [],
//             "etat_proc_collective": "in_bonis"
//         }
//     });
    
//     var output_indexed = output_array.reduce(function (periode, val) {
//         periode[val.periode.getTime()] = val
//         return periode
//     }, {});
    
    
    
//     //
//     ///
//     ///////////////
//     // Effectifs //
//     ///////////////
//     ///
//     //
    
//     map_effectif = Object.keys(v.effectif).reduce(function (map_effectif, hash) {
//         var effectif = v.effectif[hash];
//         if (effectif == null) {
//             return map_effectif
//         }
//         var effectifTime = effectif.periode.getTime()
//         map_effectif[effectifTime] = (map_effectif[effectifTime] || 0) + effectif.effectif
//         return map_effectif
//     }, {})

    
    
//     Object.keys(map_effectif).forEach(time =>{
//         time_d = new Date(parseInt(time))
//         time_offset = DateAddMonth(time_d, -offset_effectif -1)
//         if (time_offset.getTime() in output_indexed){
//             output_indexed[time_offset.getTime()].effectif = map_effectif[time]
//             output_indexed[time_offset.getTime()].date_effectif = time_d
//         }
//     }
// )

//  output_array.forEach(function (val, index) {
//      if (val.effectif == null) {
//          delete output_indexed[val.periode.getTime()]
//          delete output_array[index]
//     }
//  })


// // inscription des effectifs dans les périodes
// // value_array.map(function (val) {
// //     var effectifDate = DateAddMonth(val.periode, offset_effectif)
// //     var historyDate = DateAddMonth(val.periode, offset_effectif - 12)
// //     var historyPeriods = generatePeriodSerie(historyDate, effectifDate)
// //     val.effectif_date = effectifDate
// //     val.effectif = map_effectif[effectifDate.getTime()]
// //     val.lag_effectif = map_effectif[historyDate.getTime()]
// //     historyPeriods.map(function (p) {
// //         val.effectif_history[p.getTime()] = map_effectif[p.getTime()]
// //     })
// // })

// //
// ///
// ////////////////////////
// // activite partielle // 
// //////////////////////// 
// ///
// //

// var apart = Object.keys(v.apdemande).reduce((apart, hash) => {
//     apart[v.apdemande[hash].id_demande.substring(0, 10)] = {
//         "demande": hash,
//         "consommation": []
//     }
//     return apart
// }, {})



// Object.keys(v.apconso).forEach(hash => {
//     var valueap = v.apconso[hash]
//     if (valueap.id_conso.substring(0, 10) in apart) {
//         apart[valueap.id_conso.substring(0, 10)].consommation.push(hash)
//     }
// })

// Object.keys(apart).forEach(k => {
//     v.apdemande[apart[k].demande].hash_consommation = apart[k].consommation
//     for (j in apart[k].consommation) {
//         v.apconso[apart[k].consommation[j]].hash_demande = apart[k].demande;
//     }
// })

// Object.keys(v.apconso).forEach(
//     function (h) {
//         var conso = v.apconso[h]
//         if (conso.hash_demande) {
//             var time = conso.periode.getTime()
//             if (time in output_indexed){
//                 output_indexed[time].apart_heures_consommees = output_indexed[time].apart_heures_consommees + conso.heure_consomme;
//                 output_indexed[time].apart_motif_recours = v.apdemande[conso.hash_demande].motif_recours_se;
//             }
//         }
//     })
    
//     //
//     ///
//     ////////////
//     // delais //
//     ////////////
//     ///
//     //
    
//     Object.keys(v.delai).map(
//         function (hash) {
//             var delai = v.delai[hash]
//             var date_creation = new Date(Date.UTC(delai.date_creation.getUTCFullYear(), delai.date_creation.getUTCMonth(), 1, 0, 0, 0, 0))
//             var date_echeance = new Date(Date.UTC(delai.date_echeance.getUTCFullYear(), delai.date_echeance.getUTCMonth(), 1, 0, 0, 0, 0))
//             var pastYearTimes = generatePeriodSerie(date_creation, date_echeance).map(function (date) { return date.getTime() })
//             pastYearTimes.map(
//                 function(time){
//                     if (time in output_indexed) {
//                         var remaining_months = (date_echeance.getUTCMonth() - new Date(time).getUTCMonth()) +
//                         12*(date_echeance.getUTCFullYear() - new Date(time).getUTCFullYear()) 
//                         output_indexed[time].delai = remaining_months;
//                         output_indexed[time].duree_delai = delai.duree_delai
//                         output_indexed[time].montant_echeancier = delai.montant_echeancier
                        
//                     }
//                 }
//             )
//         }
//     ) 
    
//     //
//     ///
//     //////////////////
//     // defaillance  //
//     //////////////////
//     ///
//     //
    
//     // On filtre altares pour ne garder que les codes qui nous intéressents
//     var altares_codes  =  Object.keys(v.altares).reduce(function(events,hash) {
//         var altares_event = v.altares[hash]
        
        
//         var etat = altaresToHuman(altares_event.code_evenement)
        
//         if (etat != null)
//         events.push({"etat": etat, "date_proc_col": new Date(altares_event.date_effet)})
        
//         return(events)
//     },[{"etat" : "in_bonis", "date_proc_col" : new Date(0)}]).sort(
//         function(a,b) {return(a.date_proc_col.getTime() > b.date_proc_col.getTime())}
//     )
    
    
    
//     altares_codes.forEach(
//         function (event) {
//             var periode_effet = new Date(Date.UTC(event.date_proc_col.getFullYear(), event.date_proc_col.getUTCMonth(), 1, 0, 0, 0, 0))
//             var time_til_last = Object.keys(output_indexed).filter(val => {return (val >= periode_effet)})
//             time_til_last.forEach(time => {
//                 if (time in output_indexed) {
//                     output_indexed[time].etat_proc_collective = event.etat
//                     output_indexed[time].date_proc_collective = event.date_proc_col
//                 }
//             })
//         }
//     )
    
    
//     // Object.keys(v.altares).forEach(
//     //     function (hash) {
//     //         var altares = v.altares[hash]
//     //         var periode_effet = new Date(Date.UTC(altares.date_effet.getUTCFullYear(), altares.date_effet.getUTCMonth(), 1, 0, 0, 0, 0))
//     //         var periode_outcome = new Date(Date.UTC(altares.date_effet.getUTCFullYear() - 1, altares.date_effet.getUTCMonth(), 1, 0, 0, 0, 0))
//     //         var pastYearTimes = generatePeriodSerie(periode_outcome, periode_effet).map(function (date) { return date.getTime() })
//     //         pastYearTimes.map(
//     //             function (time) {
//     //                 if (time in output) {
//     //                     output[time].date_defaillance = altares.date_effet
//     //                     output[time].outcome_0_12 = "default";
//     //                 }
//     //             }
//     //         )
//     //     }
//     // )
    
    
    
    
    
//     //
//     ///
//     ////////////////////////////
//     // Cotisation et débits  ///
//     ////////////////////////////
//     ///
//     //
    
//     var value_cotisation = {}
    
//     Object.keys(v.cotisation).forEach(function (h) {
//         var cotisation = v.cotisation[h]
//         var periode_cotisation = generatePeriodSerie(cotisation.periode.start, cotisation.periode.end)
//         periode_cotisation.forEach(function (date_cotisation) {
//             value_cotisation[date_cotisation.getTime()] = (value_cotisation[date_cotisation.getTime()] || []).concat(cotisation.du / periode_cotisation.length)
//         })
//     })
    
//     var value_dette = {}
    
//     Object.keys(v.debit).forEach(function (h) {
//         var debit = v.debit[h]
//         if (debit.part_ouvriere + debit.part_patronale > 0) {

//             var debit_suivant = (v.debit[debit.debit_suivant] || {"date_traitement" : date_fin})
//             date_limite = date_fin//new Date(new Date(debit.periode.start).setFullYear(debit.periode.start.getFullYear() + 1))
//             date_traitement_debut = new Date(
//                 Date.UTC(debit.date_traitement.getFullYear(), debit.date_traitement.getUTCMonth())
//             )
            
//             date_traitement_fin = new Date(
//                 Date.UTC(debit_suivant.date_traitement.getFullYear(), debit_suivant.date_traitement.getUTCMonth())
//             )
            
//             periode_debut = (date_traitement_debut.getTime() >= date_limite.getTime() ? date_limite : date_traitement_debut)
//             periode_fin = (date_traitement_fin.getTime() >= date_limite.getTime() ? date_limite : date_traitement_fin)
            
//             generatePeriodSerie(periode_debut, periode_fin).map(function (date) {
//                 time = date.getTime()
//                 value_dette[time] = (value_dette[time] || []).concat([{ "periode": debit.periode.start, "part_ouvriere": debit.part_ouvriere, "part_patronale": debit.part_patronale }])
//             })
//         }
//     })    

//     Object.keys(output_indexed).forEach(function (time) {
//         if (time in value_cotisation){
//             output_indexed[time].cotisation = value_cotisation[time].reduce((a,cot) => a + cot,0)
//         }
        
//         if (time in value_dette) {
//             output_indexed[time].debit_array = value_dette[time]
//         }
//     })

//     output_array.forEach(function (val) {
        
//         val.montant_dette = val.debit_array.reduce(function (m, dette) {
//             m.part_ouvriere += dette.part_ouvriere
//             m.part_patronale += dette.part_patronale
//             return m
//         }, { "part_ouvriere": 0, "part_patronale": 0 })
        
//         val.montant_part_ouvriere = val.montant_dette.part_ouvriere
//         val.montant_part_patronale = val.montant_dette.part_patronale

//         delete val.montant_dette
//         delete val.debit_array
        
//     })

//     //
//     ///
//     /////////
//     // CCSF// 
//     /////////
//     ///
//     //
//     var ccsfHashes = Object.keys(v.ccsf || {}) 

//     output_array.forEach(val => {        
//         var optccsf = ccsfHashes.reduce( 
//             function (accu, hash) { 
//                 ccsf = v.ccsf[hash] 
//                 if (ccsf.date_traitement.getTime() < val.periode.getTime() && ccsf.date_traitement.getTime() > accu.date_traitement.getTime()) { 
//                     accu = ccsf 
//                 } 
//                 return(accu)
//             }, 
//             { 
//                 date_traitement: new Date(0) 
//             } 
//         )         
        
//         if (optccsf.date_traitement.getTime() != 0) { 
//             val.date_ccsf = optccsf.date_traitement 
//         } 
//     })
    
       
//     //
//     ///
//     ////////////
//     // Sirene //
//     ////////////
//     ///
//     //
    
//     var sireneHashes = Object.keys(v.sirene || {})

//     output_array.forEach(val => {
//         // geolocalisation
    
//         if (sireneHashes.length != 0) {
//             sirene = v.sirene[sireneHashes[0]]
//         }
        
//         val.lattitude = (sirene || { "lattitude": null }).lattitude
//         val.longitude = (sirene || { "longitude": null }).longitude
//         val.region = (sirene || {"region": null}).region
//         val.departement = (sirene || {"departement": null}).departement
//         val.code_ape  = (sirene || { "ape": null}).ape
//         val.raison_sociale = (sirene || {"raisonsociale": null}).raisonsociale
//         val.activite_saisonniere = (sirene || {"activitesaisoniere": null}).activitesaisoniere
//         val.productif = (sirene || {"productif": null}).productif
//         val.debut_activite = (sirene || {"debut_activite":null}).debut_activite.getFullYear()
//         val.tranche_ca = (sirene || {"trancheca":null}).trancheca
//         val.indice_monoactivite = (sirene || {"indicemonoactivite": null}).indicemonoactivite  
        
//     })


//     return_value = { "siren": k.substring(0, 9)}
//     return_value[k] = output_array
//     return return_value
}