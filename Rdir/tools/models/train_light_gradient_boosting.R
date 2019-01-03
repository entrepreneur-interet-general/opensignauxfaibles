train_light_gradient_boosting <- function(
  database,
  last_batch,
  training_date_inf,
  training_date_sup,
  algo,
  min_effectif,
  fields,
  x_fields_model
  ){

  raw_data <- connect_to_database(
    database,
    "Features",
    last_batch,
    date_inf = training_date_inf,
    date_sup = training_date_sup,
    algo = algorithm,
    min_effectif = min_effectif,
    fields = fields
    )

  train <- as.h2o(raw_data)
  rm(raw_data)

  train["outcome"] <- h2o.relevel(x = train["outcome"], y = "non_default")
  
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

  # FIX ME save model
  # FIX ME save te_map

  return(TRUE)
}
