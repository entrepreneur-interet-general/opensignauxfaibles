function finalize(k, v) {

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
    }, { "siren": k })

    v = Object.keys(v).filter(k => k !== 'siren' && k !== 'compact').reduce((accu, type) => {
        accu[type] = Object.keys(v[type]).map(k => v[type][k])
        return accu
    }, {})

    return v
    // liste_periodes = generatePeriodSerie(date_debut, date_fin)

    // var value_array = liste_periodes.map(function (e) {
    //     return {
    //         "siren": v.siren,
    //         "periode": e,
    //         "arrete_bilan_bdf": new Date(0),
    //         "annee_bilan_diane": 0
    //     }
    // });

    // v.bdf = (v.bdf || {})
    // v.diane = (v.diane || {})

    // Object.keys(v.bdf).forEach(hash => {
    //     value_array.forEach(siren_periode => {
    //         if (siren_periode.arrete_bilan_bdf.getTime() < v.bdf[hash].arrete_bilan.getTime()
    //             && siren_periode.periode.getTime() > v.bdf[hash].arrete_bilan.getTime()) {
    //             siren_periode.arrete_bilan_bdf = v.bdf[hash].arrete_bilan
    //             //periode.secteur = v.bdf[hash].secteur
    //             siren_periode.poids_frng = v.bdf[hash].poids_frng
    //             siren_periode.taux_marge = v.bdf[hash].taux_marge
    //             siren_periode.delai_fournisseur = v.bdf[hash].delai_fournisseur
    //             siren_periode.dette_fiscale = v.bdf[hash].dette_fiscale
    //             siren_periode.financier_court_terme = v.bdf[hash].financier_court_terme
    //             siren_periode.frais_financier = v.bdf[hash].frais_financier
    //         }
    //     })
    // })

    // Object.keys(v.diane).forEach(hash => {
    //     value_array.forEach(siren_periode => {
    //         annee_bilan = v.diane[hash].annee 
    //         if (siren_periode.annee_bilan_diane < annee_bilan
    //             && siren_periode.periode.getTime() > Date.UTC(annee_bilan+1,0,0,0,0,0,0)) {
    //                 siren_periode.annee_bilan_diane = annee_bilan 
    //                 siren_periode.bilan_diane_absent = v.diane[hash].CA == null

    //                 Object.keys(v.diane[hash]).filter( k => {
    //                     var omit = ["annee","marquee", "nom_entreprise","numero_siren",
    //                         "statut_juridique", "procedure_collective"]
    //                     return (v.diane[hash][k] != null &&  !(omit.includes(k)))
    //                 }
    //             ).forEach( k => {       

    //                      siren_periode[k] = v.diane[hash][k]
    //                 }                   
    //             )           
                
    //             // siren_periode.procedure_collective_diane = v.diane[hash].procedure_collective 
    //             // siren_periode.CA = v.diane[hash].CA 
    //             // siren_periode.valeur_ajoutee = v.diane[hash].valeur_ajoutee 
    //             // siren_periode.resultat_net_consolide = v.diane[hash].resultat_net_consolide 
    //             // siren_periode.capacite_autofinanc_avant_repartition = v.diane[hash].capacite_autofinanc_avant_repartition 
    //             // siren_periode.capital_social_ou_individuel = v.diane[hash].capital_social_ou_individuel 
    //             // siren_periode.capitaux_propres_du_groupe = v.diane[hash].capitaux_propres_du_groupe 
    //             // siren_periode.fonds_de_roul_net_global = v.diane[hash].fonds_de_roul_net_global 
    //             // siren_periode.endettement_pourcent = v.diane[hash].endettement_pourcent 
    //             // siren_periode.liquidite_reduite = v.diane[hash].liquidite_reduite 
    //             // siren_periode.rentabilite_nette_pourcent = v.diane[hash].rentabilite_nette_pourcent 
    //             // siren_periode.rend_des_capitaux_propres_nets_pourcent = v.diane[hash].rend_des_capitaux_propres_nets_pourcent 
    //             // siren_periode.rend_des_ress_durables_nettes_pourcent = v.diane[hash].rend_des_ress_durables_nettes_pourcent 
    //             // siren_periode.effectif_consolide = v.diane[hash].effectif_consolide 
    //             // siren_periode.total_actif_immob = v.diane[hash].total_actif_immob 
    //             // siren_periode.total_immob_fin = v.diane[hash].total_immob_fin 
    //             // siren_periode.total_immob_corp = v.diane[hash].total_immob_corp 
    //             // siren_periode.total_immob_incorp = v.diane[hash].total_immob_incorp 
    //             // siren_periode.stocks = v.diane[hash].stocks 
    //             // siren_periode.produits_intermed_et_finis = v.diane[hash].produits_intermed_et_finis 
    //             // siren_periode.marchandises = v.diane[hash].marchandises 
    //             // siren_periode.en_cours_de_prod_de_biens = v.diane[hash].en_cours_de_prod_de_biens 
    //             // siren_periode.matières_prem_approv = v.diane[hash].matières_prem_approv 
    //             // siren_periode.creances_expl = v.diane[hash].creances_expl 
    //             // siren_periode.clients_et_cptes_ratt = v.diane[hash].clients_et_cptes_ratt 
    //             // siren_periode.av_et_ac_sur_commandes = v.diane[hash].av_et_ac_sur_commandes 
    //             // siren_periode.disponibilites = v.diane[hash].disponibilites 
    //             // siren_periode.total_actif_circ_ch_const_av = v.diane[hash].total_actif_circ_ch_const_av 
    //             // siren_periode.total_actif = v.diane[hash].total_actif 
    //             // siren_periode.capitaux_propres_groupe = v.diane[hash].capitaux_propres_groupe 
    //             // siren_periode.resultat_consolide_part_du_groupe = v.diane[hash].resultat_consolide_part_du_groupe 
    //             // siren_periode.total_dettes_fin = v.diane[hash].total_dettes_fin 
    //             // siren_periode.total_dette_expl_et_divers = v.diane[hash].total_dette_expl_et_divers 
    //             // siren_periode.dettes_fourn_et_cptes_ratt = v.diane[hash].dettes_fourn_et_cptes_ratt 
    //             // siren_periode.dettes_fiscales_et_sociales = v.diane[hash].dettes_fiscales_et_sociales 
    //             // siren_periode.dettes_sur_immob_cptes_ratt = v.diane[hash].dettes_sur_immob_cptes_ratt 
    //             // siren_periode.total_du_passif = v.diane[hash].total_du_passif 
    //             // siren_periode.chiffre_affaires_net = v.diane[hash].chiffre_affaires_net 
    //             // siren_periode.chiffre_affaires_net_en_france = v.diane[hash].chiffre_affaires_net_en_france 
    //             // siren_periode.chiffre_affaires_net_lie_aux_exportations = v.diane[hash].chiffre_affaires_net_lie_aux_exportations 
    //             // siren_periode.salaires_et_traitements = v.diane[hash].salaires_et_traitements 
    //             // siren_periode.charges_sociales = v.diane[hash].charges_sociales 
    //             // siren_periode.total_des_charges_expl = v.diane[hash].total_des_charges_expl 
    //             // siren_periode.resultat_expl = v.diane[hash].resultat_expl 
    //             // siren_periode.total_des_produits_fin = v.diane[hash].total_des_produits_fin 
    //             // siren_periode.total_des_charges_fin = v.diane[hash].total_des_charges_fin 
    //             // siren_periode.resultat_courant_avant_impots = v.diane[hash].resultat_courant_avant_impots 
    //             // siren_periode.resultat_financier = v.diane[hash].resultat_financier 
    //             // siren_periode.resultat_exceptionnel = v.diane[hash].resultat_exceptionnel 
    //             // siren_periode.total_des_charges = v.diane[hash].total_des_charges 
    //             // siren_periode.total_des_produits = v.diane[hash].total_des_produits 
    //             // siren_periode.frais_de_RetD = v.diane[hash].frais_de_RetD 
    //             // siren_periode.conces_brev_et_droits_sim = v.diane[hash].conces_brev_et_droits_sim 
    //             // siren_periode.note_preface = v.diane[hash].note_preface 
    //             // siren_periode.nombre_etab_secondaire = v.diane[hash].nombre_etab_secondaire 
    //         }
    //     })
    // })

    // value_array.forEach((periode, index) => {
    //     if ((periode.arrete_bilan_bdf||new Date(0)).getTime() == 0 && (periode.annee_bilan_diane || 0) == 0) {
    //         delete value_array[index]
    //     }
    // })

    // return {"siren": k, "entreprise": value_array}
}