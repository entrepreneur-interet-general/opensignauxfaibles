light_gradient_boosting <- function(actual_period, last_batch, algorithm, min_effectif = 10){

  fields <- c(
    "siret",
    "siren",
    "periode",
    # URSSAF
    "montant_part_patronale",
    "montant_part_ouvriere",
    "effectif",
    "effectif_entreprise",
    "effectif_consolide",
    "montant_echeancier",
    "delai",
    "duree_delai",
    # URSSAF 12m, 
    "ratio_dette",
    "ratio_dette_moy12m",
    "cotisation_moy12m",
    # URSSAF past
    "montant_part_patronale_past_1",
    "montant_part_ouvriere_past_1",
    "montant_part_patronale_past_2",
    "montant_part_ouvriere_past_2",
    "montant_part_patronale_past_3",
    "montant_part_ouvriere_past_3",
    "montant_part_patronale_past_6",
    "montant_part_ouvriere_past_6",
    "montant_part_patronale_past_12",
    "montant_part_ouvriere_past_12",
    "effectif_past_6",
    "effectif_past_12",
    "effectif_past_18",
    "effectif_past_24",
    # ALTARES
    "etat_proc_collective",
    # DIANE
    "CA",
    "resultat_net_consolide",
    "resultat_expl",
    "taux_marge",
    "delai_fournisseur",
    "poids_frng",
    "frais_financier",
    "financier_court_terme",
    "dette_fiscale",
    "capacite_autofinanc_avant_repartition",
    "ratio_marge_operationnelle",
    "taux_rotation_stocks",
    "ratio_productivite",
    "ratio_export",
    "ratio_delai_client",
    "liquidite_reduite",
    "rentabilite_nette_pourcent",
    "endettement_pourcent",
    "rend_des_capitaux_propres_nets_pourcent",
    "rend_des_ress_durables_nettes_pourcent",
    "clients_et_cptes_ratt",
    "frais_de_RetD",
    "produits_intermed_et_finis",
    "marchandises",
    "en_cours_de_prod_de_biens",
    "matieres_prem_approv",
    "chiffre_affaires_net_lie_aux_exportations",
    # DIANE_past
    "CA_past_1",
    "resultat_net_consolide_past_1",
    "resultat_expl_past_1",
    "taux_marge_past_1",
    "delai_fournisseur_past_1",
    "poids_frng_past_1",
    "frais_financier_past_1",
    "financier_court_terme_past_1",
    "dette_fiscale_past_1",
    "capacite_autofinanc_avant_repartition_past_1",
    "ratio_marge_operationnelle_past_1",
    "taux_rotation_stocks_past_1",
    "ratio_productivite_past_1",
    "ratio_export_past_1",
    "ratio_delai_client_past_1",
    "liquidite_reduite_past_1",
    "rentabilite_nette_pourcent_past_1",
    "endettement_pourcent_past_1",
    "rend_des_capitaux_propres_nets_pourcent_past_1",
    "rend_des_ress_durables_nettes_pourcent_past_1",
    "clients_et_cptes_ratt_past_1",
    "frais_de_RetD_past_1",
    "produits_intermed_et_finis_past_1",
    "marchandises_past_1",
    "en_cours_de_prod_de_biens_past_1",
    "matieres_prem_approv_past_1",
    "chiffre_affaires_net_lie_aux_exportations_past_1",
    "CA_past_2",
    "resultat_net_consolide_past_2",
    "resultat_expl_past_2",
    "taux_marge_past_2",
    "delai_fournisseur_past_2",
    "poids_frng_past_2",
    "frais_financier_past_2",
    "financier_court_terme_past_2",
    "dette_fiscale_past_2",
    "capacite_autofinanc_avant_repartition_past_2",
    "ratio_marge_operationnelle_past_2",
    "taux_rotation_stocks_past_2",
    "ratio_productivite_past_2",
    "ratio_export_past_2",
    "ratio_delai_client_past_2",
    "liquidite_reduite_past_2",
    "rentabilite_nette_pourcent_past_2",
    "endettement_pourcent_past_2",
    "rend_des_capitaux_propres_nets_pourcent_past_2",
    "rend_des_ress_durables_nettes_pourcent_past_2",
    "clients_et_cptes_ratt_past_2",
    "frais_de_RetD_past_2",
    "produits_intermed_et_finis_past_2",
    "marchandises_past_2",
    "en_cours_de_prod_de_biens_past_2",
    "matieres_prem_approv_past_2",
    "chiffre_affaires_net_lie_aux_exportations_past_2",
    "nombre_etab_secondaire",
    "nbr_etablissements_connus",
    # APART
    "apart_heures_consommees",
    "apart_heures_autorisees",
    # SIRENE
    "code_ape",
    "code_naf",
    "debut_activite",
    "activite_saisonniere",
    "productif"
    )

  date_inf <- as.Date("2015-01-01")
  date_sup <- as.Date("2017-01-01")

  raw_data <- connect_to_database(
    "Features",
    last_batch,
    date_inf = date_inf,
    date_sup = date_sup,
    algo = algorithm,
    min_effectif = min_effectif,
    fields = fields)

  current_data <- connect_to_database(
    "Features",
    last_batch,
    date_inf = actual_period %m-% months(1),
    date_sup = actual_period %m+% months(1),
    algo = algorithm,
    min_effectif = min_effectif,
    fields = fields)

  raw_data <- raw_data %>%
    objective_default_or_failure(n_months = 3, threshold = 1, lookback = 18) %>%
    set_objective("default")

  out <- feature_engineering_std(raw_data, current_data)
  raw_data <- out[[1]]
  current_data <- out[[2]]

  ref_f_e <- feature_engineering_create(raw_data)

  out <- feature_engineering_apply(ref_f_e, raw_data, current_data)

  raw_data <- out[[1]]
  current_data <- out[[2]]

  rm(out)

  h2o.init(ip = "localhost", port = 4444)

  train <- as.h2o(raw_data)
  current <- as.h2o(current_data)

  train["outcome"] <- h2o.relevel(x = train["outcome"], y = "non_default")

  te_map <- h2o.target_encode_create(
    train,
    x = list(c("code_naf"), c("code_ape_niveau2"), c("code_ape_niveau3"), c("code_ape_niveau4"), c("code_ape")),
    y = "outcome")

  train <- h2o.target_encode_apply(
    train,
    x = list(c("code_naf"), c("code_ape_niveau2"), c("code_ape_niveau3"), c("code_ape_niveau4"), c("code_ape")),
    y = "outcome",
    target_encode_map = te_map,
    holdout_type = "LeaveOneOut",
    blended_avg = TRUE,
    seed = 1234)

  current <- h2o.target_encode_apply(
    current,
    x = list(c("code_naf"), c("code_ape_niveau2"), c("code_ape_niveau3"), c("code_ape_niveau4"), c("code_ape")),
    y = "outcome",
    target_encode_map = te_map,
    holdout_type = "None",
    blended_avg = FALSE,
    fold_column = "fold_column",
    noise_level = 0)

  x_medium <- c("montant_part_patronale",
    "ratio_dette",
    "ratio_dette_moy12m",
    "etat_proc_collective_num",
    "TargetEncode_code_ape_niveau3",
    "cotisation_moy12m",
    "frais_financier_distrib_APE1",
    "taux_marge_distrib_APE1",
    "montant_part_patronale_past_3",
    "ratio_liquidite_reduite_distrib_APE1",
    "dette_fiscale",
    "ratio_delai_client_distrib_APE1",
    "montant_part_patronale_past_1",
    "montant_part_patronale_past_2",
    "ratio_rend_capitaux_propres",
    "taux_marge",
    "poids_frng_distrib_APE1",
    "delai_fournisseur_distrib_APE1",
    "taux_rotation_stocks_distrib_APE1",
    "effectif",
    "ratio_rentabilite_nette_distrib_APE1",
    "ratio_export_distrib_APE1",
    "TargetEncode_code_ape_niveau2",
    "effectif_past_12",
    "montant_part_ouvriere",
    "financier_court_terme_distrib_APE1",
    #"effectif_entreprise",
    "age",
    "ratio_liquidite_reduite",
    "ratio_productivite_distrib_APE1",
    "frais_financier",
    "financier_court_terme",
    "ratio_delai_client",
    "TargetEncode_code_naf",
    "resultat_net_consolide",
    "taux_rotation_stocks",
    "nombre_etab_secondaire",
    "nbr_etablissements_connus",
    "CA",
    "chiffre_affaires_net_lie_aux_exportations",
    "ratio_dette_delai",
    "ratio_marge_operationnelle_distrib_APE1",
    "poids_frng")

  y <- "outcome"

  model <- h2o.xgboost(
    model_id = "Model_train",
    x = x_medium,
    y = y,
    training_frame = train,
    tree_method = "hist",
    grow_policy = "lossguide",
    learn_rate = 0.1,
    max_depth = 4,
    ntrees = 60,
    seed = 123
    )
  prediction <- as.tibble(h2o.cbind(current, h2o.predict(model, current)))

  prediction <- prediction %>% mutate(
    # H2O bug ??
    periode =  as.Date(structure(periode / 1000,
        class = c("POSIXct", "POSIXt")))
    ) %>%
  rename(prob = default) %>%
  select(predict, prob, siret, periode)

pred_data <- prediction %>%
  group_by(siret) %>%
  arrange(siret,periode) %>%
  mutate(last_prob = lag(prob)) %>%
  ungroup() %>%
  mutate(diff = prob - last_prob)


export_fields <-  c(
  "siret",
  "periode",
  "raison_sociale",
  "departement",
  "region",
  "prob",
  "diff",
  "date_ccsf",
  "etat_proc_collective",
  "date_proc_collective",
  "connu",
  "interessante_urssaf",
  #"default_urssaf",
  "effectif",
  "libelle_naf",
  "libelle_ape5",
  "code_ape",
  "montant_part_ouvriere",
  "montant_part_patronale",
  "CA",
  "CA_past_1",
  "resultat_net_consolide",
  "resultat_net_consolide_past_1",
  "resultat_expl",
  "resultat_expl_past_1",
  "poids_frng",
  "taux_marge",
  "frais_financier",
  "financier_court_terme",
  "delai_fournisseur",
  "dette_fiscale",
  "apart_heures_consommees",
  "apart_heures_autorisees",
  "cotisation_moy12m",
  "montant_majorations",
  "numero_compte_urssaf",
  "exercice_bdf",
  "exercice_diane"
  )

# Export
pred_data %>%
  filter(periode == actual_period) %>%
  prepare_for_export(export_fields = export_fields, last_batch = last_batch, algorithm = algorithm) %>%
  export(batch = "1810")

# Returns H2O frames and model
return(list(train_data = train, current_data = current, model = model))
}
