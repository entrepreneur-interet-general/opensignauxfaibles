full_light_gradient_boosting <- function(
  database,
  collection,
  periods,
  last_batch,
  min_effectif = 10,
  retrain_model = FALSE,
  training_date_inf = as.Date("2015-01-01"),
  training_date_sup = as.Date("2017-01-01"),
  type = "dataframe") {

  fields <- get_fields(training = FALSE)
  x_fields_model <- get_fields(training = TRUE)

  if (retrain_model) {

    data_frame <- connect_to_database(
      database,
      collection,
      last_batch,
      siren = NULL,
      date_inf = training_date_inf,
      date_sup = training_date_sup,
      min_effectif = min_effectif,
      fields = fields,
      code_ape = NULL,
      type = type,
      subsample = 500000
      )

    out <- prepare_frame(
      data_to_prepare = data_frame,
      test_or_train = "train",
      save_or_load_map = TRUE, #save
      outcome = "outcome"
      )


    train_data <- out[["data"]]
    te_map <- out[["te_map"]]

    model <- train_light_gradient_boosting(
      h2o_train_data = train_data,
      x_fields_model = x_fields_model,
      save_results = TRUE
      )

  } else {
    model <- load_h2o_object("lgb", "model", last = TRUE)
    te_map <- load_h2o_object("te_map", "temap", last = TRUE)
  }

  pred_data <- predict_on_last_batch(
    model,
    database,
    collection,
    te_map,
    last_batch,
    periods,
    min_effectif = min_effectif,
    fields = fields
    )

  export_fields <-  c(
    "siret",
    "periode",
    "raison_sociale",
    "departement",
    "region",
    "prob",
    "diff",
    "connu",
    "date_ccsf",
    "etat_proc_collective",
    "date_proc_collective",
    "interessante_urssaf",
    #"default_urssaf",
    "effectif",
    "libelle_naf",
    "libelle_ape5",
    "code_ape",
    "montant_part_ouvriere",
    "montant_part_patronale",
    "ca",
    "ca_past_1",
    "benefice_ou_perte",
    "benefice_ou_perte_past_1",
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
    "exercice_bdf",
    "exercice_diane"
    )

  # Export
export_results(
  pred_data,
  export_fields,
  database,
  collection,
  last_batch)

# Returns H2O frames and model
return(TRUE)
}

export_results <- function(
  pred_data,
  export_fields,
  database,
  collection,
  last_batch){

  res <- pred_data %>%
    prepare_for_export(
      export_fields = export_fields,
      database = database,
      collection = collection,
      last_batch = last_batch)

  export(res, batch = last_batch)
  return(res)
}






