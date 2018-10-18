feature_engineering <- function(train_set,
                                ...,
                                oversampling = FALSE,
                                ratios_urssaf = TRUE,
                                avec_ratios_financiers = TRUE,
                                imputation = FALSE,
                                number_imputations = NULL,
                                number_iterations = NULL,
                                past_trend = TRUE,
                                past_trend_vars = NULL,
                                past_trend_lookbacks = NULL,
                                past_trend_vars_years = NULL,
                                past_trend_lookbacks_years = NULL,
                                quantile = TRUE,
                                quantile_vars = NULL,
                                quantile_levels = NULL,
                                failure_rate_APE = FALSE) {
  arguments <- list(...)
  cat('Data preproccessing ...', '\n')

  ratios_financiers = c(
    "CA",
    "taux_marge",
    "resultat_net_consolide",
    "resultat_expl",
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

  if (past_trend) {
    if (is.null(past_trend_vars)) {
      past_trend_vars <- c(
        'apart_heures_consommees',
        'montant_part_ouvriere',
        'montant_part_patronale'
      )
      past_trend_vars_years <- ratios_financiers
      cat("Taking default value for past trends variables", '\n')
    }
    if (is.null(past_trend_lookbacks)) {
      past_trend_lookbacks <-  c(1, 2, 3, 6, 12)
      past_trend_lookbacks_years <- c(1)
      cat("Taking default value for past trends lookbacks", '\n')
    }
    past_trend_lookbacks_ym <- past_trend_lookbacks_years * 12
  }

  if (quantile) {
    if (is.null(quantile_vars)) {
      quantile_vars <- c(ratios_financiers, 'effectif')

      cat("Taking default value for quantile_APE variables", '\n')
    }
    if (is.null(quantile_levels)) {
      quantile_levels <-  c(1, 2)
      cat("Taking default value for quantile_APE levels", '\n')
    }
  }

  ##
  ###
  ############################################
  ## Auxiliary function applied to all sets ##
  ############################################
  ###
  ##

  aux_fun <- function(data, n_impute, n_iter) {
    assertthat::assert_that(all(c('siret', 'periode') %in% names(data)))

    data <- data %>%
      mutate_if(is.POSIXct, as.Date)

    ##################
    ##  PreFiltering   ##
    ##################

    wrong_sirets <- c('34117000900037',
                      '38520230400023')

    data <- data %>%
      filter(
        is.na(effectif_consolide) | effectif_consolide < 500,
        is.na(effectif_entreprise) | effectif_entreprise < 500,
        is.na(CA) | CA < 100000,!siret %in% wrong_sirets
      )

    ##################
    ## Shaping data ##
    ##################

    # Removing errors in etat_proc_collective
    data <- data %>% group_by(siret) %>%
      arrange(siret, periode) %>%
      mutate(
        etat_proc_collective_lag = lag(etat_proc_collective),
        etat_proc_collective = ifelse(
          etat_proc_collective_lag == 'in_bonis' &
            etat_proc_collective == 'continuation',
          'plan_redressement',
          etat_proc_collective
        ),
        etat_proc_collective = ifelse(
          etat_proc_collective_lag == 'in_bonis' &
            etat_proc_collective == 'sauvegarde',
          'plan_sauvegarde',
          etat_proc_collective
        )
      ) %>%
      ungroup() %>%
      select(-etat_proc_collective_lag)


    data <- data %>%
      mutate(
        activite_saisonniere = ifelse(activite_saisonniere == "S", 1, 0),
        productif = ifelse(productif == "O", 1, 0),
        etat_proc_collective = factor(
          etat_proc_collective,
          levels = c(
            'in_bonis',
            'cession',
            'sauvegarde',
            'continuation',
            'plan_sauvegarde',
            'plan_redressement',
            'liquidation'
          )
        ),
        etat_proc_collective_num = as.numeric(etat_proc_collective)
      )

    data <- data %>%
      # select(
      #    -stocks,
      #    -total_immob_fin,
      #    -total_immob_incorp,
      #    -total_immob_corp,
      #    -creances_expl,
      #    -total_dettes_fin,
      #    -total_dette_expl_et_divers
      # ) %>%
      rename(
        ratio_liquidite_reduite = liquidite_reduite,
        ratio_rentabilite_nette = rentabilite_nette_pourcent,
        ratio_endettement = endettement_pourcent,
        ratio_rend_capitaux_propres = rend_des_capitaux_propres_nets_pourcent,
        ratio_rend_des_ress_durables = rend_des_ress_durables_nettes_pourcent
      ) %>%
      mutate(ratio_liquidite_reduite = 100 * ratio_liquidite_reduite)


    #########################
    ## New computed fields ##
    #########################



    # AGE
    assertthat::assert_that(all(c('debut_activite') %in% names(data)))
    data <- data %>%
      mutate(age = lubridate::year(as.POSIXct.Date(periode)) - debut_activite)

  # APE2
    data <- data  %>%
      mutate(code_ape_niveau2 = as.factor(substr(code_ape,1,2)),
             code_ape_niveau3 = as.factor(substr(code_ape,1,3)),
             code_ape_niveau4 = as.factor(substr(code_ape,1,4)))

    # REPLACE NA (DELAIS, COTISATION)
    assertthat::assert_that(all(
      c(
        'montant_echeancier',
        'delai',
        'duree_delai',
        'cotisation'
      ) %in% names(data)
    ))
    data <- replace_na_by("montant_echeancier", data, 0)
    data <- replace_na_by("delai", data, 0)
    data <- replace_na_by("duree_delai", data, 0)
    data <- replace_na_by("cotisation", data, 0)

    # SIMPLIFIED NAF
    assertthat::assert_that("libelle_naf_niveau1" %in% names(data))
    libelle_naf_simplifie <- data$libelle_naf_niveau1
    libelle_naf_simplifie[libelle_naf_simplifie %in% c(
      'Enseignement',
      'Activités extra-territoriales',
      'Administration publique',
      'Agriculture, sylviculture et pêche'
    )] <-
      'autre'
    data <- data %>%
      mutate(libelle_naf_simplifie = libelle_naf_simplifie)

    # Ratios URSSAF

    if (ratios_urssaf) {
      assertthat::assert_that(all(
        c(
          "apart_heures_consommees",
          "effectif",
          "montant_part_patronale",
          "montant_part_ouvriere",
          "cotisation"
        ) %in% names(data)
      ))

      # Debits
      #FIX ME toutes les périodes sont elles contigues ??

      # Correction cotisations déconnantes
      data <- data %>%
        group_by(siret) %>%
        arrange(siret, periode) %>%
        mutate(cotisation_moy12m = average_12m(cotisation)) %>%
        ungroup()



      data <-  data %>%
        mutate(
          ratio_apart = apart_heures_consommees / (effectif * 157.67),
          ratio_dette = base::ifelse(
            !is.na(cotisation_moy12m) & cotisation_moy12m > 0,
            (montant_part_patronale + montant_part_ouvriere) / cotisation_moy12m,
            NA
          ),
          ratio_dette_delai = base::ifelse(
            !is.na(duree_delai) & duree_delai > 0,
            (
              montant_part_patronale + montant_part_ouvriere - montant_echeancier * delai / (duree_delai /
                                                                                               30)
            ) / montant_echeancier,
            NA
          )
        ) %>%
        group_by(siret) %>%
        arrange(siret, periode) %>%
        mutate(ratio_dette_moy12m = average_12m(ratio_dette),
               dette_any_12m = (ratio_dette_moy12m > 0)) %>%
        ungroup()
    }

    if (avec_ratios_financiers) {
      # Liquidites et solvabilite
      data <- data %>%
        mutate(# Suspect # BFR = total_actif_circ_ch_const_av - total_des_charges_expl,
          # Suspect # tresorerie_nette = fonds_de_roul_net_global - BFR,
          ratio_CAF = 100 * capacite_autofinanc_avant_repartition / CA)

      # Rentabilite
      data <- data %>%
        mutate(
          ratio_marge_operationnelle =  100 * resultat_expl / CA,
          taux_marge_neg = taux_marge < 0,
          taux_marge_extr_neg = taux_marge < -100,
          taux_marge_extr_pos = taux_marge > 100
        )

      # Stocks
      data <- data %>%
        mutate(
          stocks = produits_intermed_et_finis + marchandises + en_cours_de_prod_de_biens + matières_prem_approv ,
          taux_rotation_stocks =  CA / stocks
        )

      # Autre
      data <- data %>%
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
    }


    ##################
    ## PAST TRENDS ###
    ##################

    if (past_trend) {
      assertthat::assert_that(all(past_trend_vars %in% names(data)))
      data <- data %>% add_past_trends(past_trend_vars,
                                       past_trend_lookbacks,
                                       type = 'lag')

      data <- data %>% add_past_trends(past_trend_vars_years,
                                       past_trend_lookbacks_ym,
                                       type = 'mean_unique')

      names_with_na <- names(data %>% select(contains('variation')))
      for (name in names_with_na)
        data <- replace_na_by(name, data, 0)

      data <- data %>%
        group_by(siret) %>%
        arrange(siret, periode) %>%
        mutate(
          effectif_diff = c(NA, diff(effectif)),
          effectif_diff_moy12 = average_12m(effectif_diff)
        ) %>%
        select(-effectif_diff) %>%
        ungroup()
    }




    ################
    ## IMPUTATION ##
    ################

    if (imputation) {
      seed <- 1234
      data <-
        impute_missing_data_BdF(data, n_impute, n_iter, seed) %>%
        as_tbl_time(periode)
    }

    ###############
    ## TRIM #######
    ###############

    # No trim for tree based methods.


    ###################
    ## POST FILTER ####
    ###################

    data <- data  %>% filter(etat_proc_collective != 'cession')

    return(data)

  }

  ##
  ###
  #######################
  ## Apply to all sets ## #######################################################################################
  #######################
  ###
  ##

  # OVERSAMPLING
  if (oversampling) {
    train_set <- oversample(train_set)
  }


  #####################
  ## APE COMPARISONS ##
  #####################


  # if (failure_rate_APE) {
  #   out <- failure_APE(train_set, ..., ape_levels = c(2, 5))
  #   train_set <- out[[1]]
  #   arguments <- out[[-1]]
  # }


  data_list <-
    c(
      list(aux_fun(train_set, 1, 5)),
      lapply(arguments, aux_fun, number_imputations, number_iterations)
    )

  train_set <- data_list[[1]]
  arguments <- data_list[-1]

  if (quantile) {

    assertthat::assert_that(all(quantile_vars %in% names(train_set)))

    out <- do.call(function(...) quantile_APE(train_set, ..., variable_names = quantile_vars, levels = quantile_levels, noise = 0.05), arguments)

  }




  return(out)

}
