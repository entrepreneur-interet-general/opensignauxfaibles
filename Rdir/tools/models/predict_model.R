predict_model <- function(
    model,
    te_map,
    database,
    last_batch,
    actual_period,
    algo,
    min_effectif,
    fields
  ){

  current_data <- connect_to_database(
    database,
    "Features",
    last_batch,
    date_inf = actual_period %m-% months(1),
    date_sup = actual_period %m+% months(1),
    algo = algorithm,
    min_effectif = min_effectif,
    fields = fields
    )

  current_data <- as.h2o(current_data)
  # FIX ME !
  current_data[["etat_proc_collective"]] <- h2o.asfactor(current_data[["etat_proc_collective"]])

  current <- h2o_target_encode(
    te_map,
    current_data,
    "test")

  prediction <- h2o.cbind(current, h2o.predict(model, current)) %>%
      as.tibble()

  prediction <- prediction %>%
    mutate(
    # H2O bug ??
    periode =  as.Date(structure(periode / 1000,
        class = c("POSIXct", "POSIXt")))
    ) %>%
  rename(prob = TRUE.) %>% # Name automatically given to default probability
  select(predict, prob, siret, periode)



  pred_data <- prediction %>%
    group_by(siret) %>%
    arrange(siret, periode) %>%
    mutate(last_prob = lag(prob)) %>%
    ungroup() %>%
    mutate(diff = prob - last_prob) %>%
    filter(periode == actual_period)

  return(pred_data)
}
