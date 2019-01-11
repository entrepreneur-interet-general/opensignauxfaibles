train_light_gradient_boosting <- function(
  database,
  last_batch,
  training_date_inf,
  training_date_sup,
  algo,
  min_effectif,
  fields,
  x_fields_model,
  reexport_csv
  ){

  browser()

  raw_data <- connect_to_database(
    database,
    "Features",
    last_batch,
    date_inf = training_date_inf,
    date_sup = training_date_sup,
    algo = algorithm,
    min_effectif = min_effectif,
    fields = fields,
    type = "csv",
    reexport_csv = reexport_csv
    )




  # train <- as.h2o(raw_data)
  train <- raw_data
  # FIX ME !
  train["etat_proc_collective"] <- h2o.asfactor(train["etat_proc_collective"])
  rm(raw_data)

  # train["outcome"] <- h2o.relevel(x = train["outcome"], y = "false")

  #
  # Target Encoding de differents groupes sectoriels
  #

  te_map <- h2o.target_encode_create(
    train,
    x = list(c("code_naf"), c("code_ape_niveau2"), c("code_ape_niveau3"), c("code_ape_niveau4"), c("code_ape")),
    y = "outcome")

  train <- h2o_target_encode(
    te_map,
    train,
    "train")

  #
  # Train the model
  #


  y <- "outcome"

  model <- h2o.xgboost(
    model_id = "Model_train",
    x = x_fields_model,
    y = y,
    training_frame = train,
    tree_method = "hist",
    grow_policy = "lossguide",
    learn_rate = 0.1,
    max_depth = 4,
    ntrees = 60,
    seed = 123
    )

  save_h2o_object(model, "lgb", "model") #lgb like Light Gradient Boosting
  save_h2o_object(te_map, "te_map", "temap")

  return(list(model = model, te_map = te_map))
}
