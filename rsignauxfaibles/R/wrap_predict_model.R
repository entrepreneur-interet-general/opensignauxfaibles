#' Predict the model on last available data
#'
#' @param model
#' @param database
#' @param collection
#' @param te_map
#' @param last_batch
#' @param periods Date, or vector of two Dates defining lower and upper time limit
#' @param min_effectif
#' @param fields
#'
#' @return
#' @export
#'
#' @examples
predict_on_last_batch <- function(
    model,
    database,
    collection,
    te_map,
    last_batch,
    periods,
    min_effectif,
    fields
  ){

  current_data <- get_last_batch(
    database,
    collection,
    last_batch,
    periods,
    fields,
    min_effectif
    )

  current_data <- prepare_frame(
    data_to_prepare = current_data,
    test_or_train = "test",
    te_map = te_map,
    save_or_load_map = FALSE)[["data"]]

  prediction  <- predict_model(
    model,
    new_data = current_data)

  pred_data <- prediction %>%
    group_by(siret) %>%
    arrange(siret, periode) %>%
    mutate(last_prob = lag(prob)) %>%
    ungroup() %>%
    mutate(diff = prob - last_prob) %>%
    filter(periode >= min(periods), periode <= max(periods))

  return(pred_data)
}

#' Make prediction on new data
#'
#' @param model An H2O model
#' @param new_data Prepared data
#'
#' @return
#' @export
#'
#' @examples
predict_model  <- function(model, new_data) {

  prediction <- h2o::h2o.cbind(
    new_data[, c("siret", "periode")],
    h2o::h2o.predict(model, new_data) %>%
      .[,c(1,3)] %>%
      setNames(list("predicted_outcome", "prob"))
    ) %>%
  tibble::as.tibble()

prediction <- prediction %>%
  mutate(
    # H2O bug ??
    periode =  as.Date(structure(periode / 1000,
        class = c("POSIXct", "POSIXt")))
    )
return(prediction)
}
