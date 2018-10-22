feature_engineering_std <- function(...){

  cat("Data preproccessing ...", "\n")

  ratios_financiers <- c(
    "CA",
    "taux_marge",
    "delai_fournisseur",
    "poids_frng",
    "frais_financier",
    "financier_court_terme",
    "ratio_CAF",
    "ratio_marge_operationnelle",
    "taux_rotation_stocks",
    "ratio_productivite",
    "ratio_export",
    "ratio_delai_client",
    "ratio_liquidite_reduite",
    "ratio_rentabilite_nette",
    "ratio_endettement",
    "ratio_rend_capitaux_propres",
    "ratio_rend_des_ress_durables",
    "ratio_RetD"
  )

  ####################
  ## Default values ##
  ####################

  past_trend_vars <- c(
    "apart_heures_consommees",
    "montant_part_ouvriere",
    "montant_part_patronale"
  )
  past_trend_vars_years <- ratios_financiers
  cat("Taking default value for past trends variables", "\n")

  past_trend_lookbacks <-  c(1, 2, 3, 6, 12)
  past_trend_lookbacks_years <- c(1)

  cat("Taking default value for past trends lookbacks", "\n")

  past_trend_lookbacks_ym <- past_trend_lookbacks_years * 12

  aux_fun <- function(my_data){
    assertthat::assert_that(all(c("siret", "periode") %in% names(my_data)))

    my_data <- my_data %>%
      mutate_if(is.POSIXct, as.Date)

    ##################
    ##  PreFiltering   ##
    ##################

    wrong_sirets <- c("34117000900037",
      "38520230400023")

    my_data <- my_data %>%
      filter(
        is.na(effectif_consolide) | effectif_consolide < 500,
        is.na(effectif_entreprise) | effectif_entreprise < 500,
        is.na(CA) | CA < 100000, !siret %in% wrong_sirets
      )

    ##################
    ## Shaping my_data ##
    ##################

    # Removing errors in etat_proc_collective
    my_data <- my_data %>% group_by(siret) %>%
      arrange(siret, periode) %>%
      mutate(
        etat_proc_collective_lag = lag(etat_proc_collective),
        etat_proc_collective = ifelse(
          etat_proc_collective_lag == "in_bonis" &
            etat_proc_collective == "continuation",
          "plan_redressement",
          etat_proc_collective
          ),
        etat_proc_collective = ifelse(
          etat_proc_collective_lag == "in_bonis" &
            etat_proc_collective == "sauvegarde",
          "plan_sauvegarde",
          etat_proc_collective
        )
        ) %>%
      ungroup() %>%
      select(-etat_proc_collective_lag)


    my_data <- my_data %>%
      mutate(
        activite_saisonniere = ifelse(
          activite_saisonniere == "S",
          1,
          0),
        productif = ifelse(productif == "O", 1, 0),
        etat_proc_collective = factor(
          etat_proc_collective,
          levels = c(
            "in_bonis",
            "cession",
            "sauvegarde",
            "continuation",
            "plan_sauvegarde",
            "plan_redressement",
            "liquidation"
          )
          ),
        etat_proc_collective_num = as.numeric(etat_proc_collective)
      )

    my_data <- my_data %>%
      rename(
        ratio_liquidite_reduite = liquidite_reduite,
        ratio_rentabilite_nette = rentabilite_nette_pourcent,
        ratio_endettement = endettement_pourcent,
        ratio_rend_capitaux_propres =
          rend_des_capitaux_propres_nets_pourcent,
        ratio_rend_des_ress_durables =
          rend_des_ress_durables_nettes_pourcent
        ) %>%
    mutate(ratio_liquidite_reduite = 100 * ratio_liquidite_reduite)


  #########################
  ## New computed fields ##
  #########################



  # AGE
  assertthat::assert_that(all(c("debut_activite") %in% names(my_data)))
  my_data <- my_data %>%
    mutate(age = lubridate::year(as.POSIXct.Date(periode)) - debut_activite)

  # APE2
  my_data <- my_data  %>%
    mutate(code_ape_niveau2 = as.factor(substr(code_ape, 1, 2)),
      code_ape_niveau3 = as.factor(substr(code_ape, 1, 3)),
      code_ape_niveau4 = as.factor(substr(code_ape, 1, 4)))

  # REPLACE NA (DELAIS, COTISATION)
  assertthat::assert_that(all(
      c(
        "montant_echeancier",
        "delai",
        "duree_delai",
        "cotisation"
        ) %in% names(my_data)
      ))
  my_data <- replace_na_by("montant_echeancier", my_data, 0)
  my_data <- replace_na_by("delai", my_data, 0)
  my_data <- replace_na_by("duree_delai", my_data, 0)
  my_data <- replace_na_by("cotisation", my_data, 0)

  # SIMPLIFIED NAF
  assertthat::assert_that("libelle_naf_niveau1" %in% names(my_data))
  libelle_naf_simplifie <- my_data$libelle_naf_niveau1
  libelle_naf_simplifie[libelle_naf_simplifie %in% c(
    "Enseignement",
    "Activités extra-territoriales",
    "Administration publique",
    "Agriculture, sylviculture et pêche"
    )] <-
      "autre"
    my_data <- my_data %>%
      mutate(libelle_naf_simplifie = libelle_naf_simplifie)

    # Ratios URSSAF

    assertthat::assert_that(all(
        c(
          "apart_heures_consommees",
          "effectif",
          "montant_part_patronale",
          "montant_part_ouvriere",
          "cotisation"
          ) %in% names(my_data)
        ))

    # Debits
    #FIX ME toutes les périodes sont elles contigues ??

    # Correction cotisations déconnantes
    my_data <- my_data %>%
      group_by(siret) %>%
      arrange(siret, periode) %>%
      mutate(cotisation_moy12m = average_12m(cotisation)) %>%
      ungroup()



    my_data <-  my_data %>%
      mutate(
        ratio_apart = apart_heures_consommees / (effectif * 157.67),
        ratio_dette = base::ifelse(
          !is.na(cotisation_moy12m) & cotisation_moy12m > 0,
          (montant_part_patronale + montant_part_ouvriere) /
            cotisation_moy12m,
          NA
          ),
        ratio_dette_delai = base::ifelse(
          !is.na(duree_delai) & duree_delai > 0,
          (
            montant_part_patronale + montant_part_ouvriere -
              montant_echeancier * delai / (duree_delai / 30)
            ) / montant_echeancier,
          NA
        )
        ) %>%
    group_by(siret) %>%
    arrange(siret, periode) %>%
    mutate(ratio_dette_moy12m = average_12m(ratio_dette),
      dette_any_12m = (ratio_dette_moy12m > 0)) %>%
    ungroup()

  # Liquidites et solvabilite
  my_data <- my_data %>%
    mutate(
      # Suspect # BFR = total_actif_circ_ch_const_av - total_des_charges_expl,
      # Suspect # tresorerie_nette = fonds_de_roul_net_global - BFR,
      ratio_CAF = 100 * capacite_autofinanc_avant_repartition / CA)

  # Rentabilite
  my_data <- my_data %>%
    mutate(
      ratio_marge_operationnelle =  100 * resultat_expl / CA,
      taux_marge_neg = taux_marge < 0,
      taux_marge_extr_neg = taux_marge < -100,
      taux_marge_extr_pos = taux_marge > 100
    )

  # Stocks
  my_data <- my_data %>%
    mutate(
      stocks = produits_intermed_et_finis + marchandises +
        en_cours_de_prod_de_biens + matières_prem_approv + marchandises,
      taux_rotation_stocks =  CA / stocks
    )

  # Autre
  my_data <- my_data %>%
    mutate(
      ratio_productivite = 100 * CA / effectif,
      ratio_export = 100 * chiffre_affaires_net_lie_aux_exportations / CA,
      ratio_delai_client = (clients_et_cptes_ratt * 360) / CA,
      ratio_RetD = frais_de_RetD / CA
    )
  # mutate(
  #   EBE = valeur_ajoutee - salaires_et_traitements - charges_sociales,
  #   ratio_autonomie_financiere = capitaux_propres_du_groupe / total_actif,
  #  #liquidites_generales = total_actif_circ_ch_const_av, #/ dettes court terme
  #  ratio_productivite = CA / effectif,
  #  #ratio_solvabilite = total_actif / (total_dettes_fin + total_dette_expl_et_divers),
  #  #ratio_gearing = (total_dettes_fin + total_dette_expl_et_divers)/ capitaux_propres_du_groupe,
  #  #ratio_independence_financiere = capitaux_propres_du_groupe / total_dettes_fin,
  #  ratio_rendement_capitaux_propres = resultat_net_consolide / capitaux_propres_du_groupe,
  #  ratio_rentabilite_economique = resultat_net_consolide / total_actif,
  #  ratio_marge_net = resultat_net_consolide / CA,
  #  #ratio_BFR = stocks + creances_expl - (total_dettes_fin + total_dette_expl_et_divers),
  #  ratio_delai_client = (clients_et_cptes_ratt*360)/CA
  #  #ratio_tresorerie_nette = poids_frng - ratio_BFR
  # )


  ##################
  ## PAST TRENDS ###
  ##################

  assertthat::assert_that(all(past_trend_vars %in% names(my_data)))
  my_data <- my_data %>% add_past_trends(past_trend_vars,
    past_trend_lookbacks,
    type = "lag")

  my_data <- my_data %>% add_past_trends(past_trend_vars_years,
    past_trend_lookbacks_ym,
    type = "mean_unique")

  names_with_na <- names(my_data %>% select(contains("variation")))
  for (name in names_with_na)
    my_data <- replace_na_by(name, my_data, 0)

  my_data <- my_data %>%
    group_by(siret) %>%
    arrange(siret, periode) %>%
    mutate(
      effectif_diff = c(NA, diff(effectif)),
      effectif_diff_moy12 = average_12m(effectif_diff)
      ) %>%
    select(-effectif_diff) %>%
    ungroup()




  ###############
  ## TRIM #######
  ###############

  # No trim for tree based methods.


  ###################
  ## POST FILTER ####
  ###################

  my_data <- my_data %>% filter(etat_proc_collective != "cession")

  return(my_data)

  }

  return(
    lapply(list(...), aux_fun)
  )
}
