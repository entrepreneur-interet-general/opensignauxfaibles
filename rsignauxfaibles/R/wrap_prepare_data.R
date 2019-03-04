#' Préparation des données pour l'entraînement ou la prédiction
#'
#' Inclut la conversion en H2O, et le target encoding.
#' La map pour le target encoding peut être sauvegardé, afin de permettre de
#' l'appliquer à de nouvelles données.
#'
#' @param data_to_prepare
#' @param test_or_train "train" if a new te_map should be computed, and applied noise; "test" otherwise
#' @param te_map Only relevant for test: specify used te_map.
#' @param save_or_load_map if test_or_train="test", then loading the map ? (no
#' te_map should be specified). if test_or_train="train", the saving the map ?
#' @param outcome to perform the target encode map on. Only needed if test_or_train = "train"
#'
#' @return
#' @export
#'
#' @examples
prepare_frame <- function(
  data_to_prepare,
  test_or_train = "train",
  te_map = NULL,
  save_or_load_map = TRUE,
  outcome = "outcome"
  ) {

  if ((!is.null(te_map) && test_or_train == "train") ||
    (!is.null(te_map) && save_or_load_map == TRUE))
    error('te_map should not be specified if it is computed (test_or_train =
      "train") or if map is loaded (save_or_load_map = TRUE)')

  h2o_data <- convert_to_h2o(data_to_prepare)

  #
  # Target Encoding de differents groupes sectoriels
  #

  if (test_or_train == "train") {
    te_map <- h2o::h2o.target_encode_create(
      h2o_data,
      x = list(c("code_naf"),
        c("code_ape_niveau2"),
        c("code_ape_niveau3")),
      y = outcome)

    if (save_or_load_map) {
      save_h2o_object(te_map, "te_map")
    }
  } else if (test_or_train == "test" && save_or_load_map) {
    te_map <- load_h2o_object("te_map", "temap", last = TRUE)
  }

  h2o_data <- h2o_target_encode(
    te_map,
    h2o_data,
    test_or_train)

  res <- list(
    data = h2o_data,
    te_map = te_map)
  return(res)
}

#' Get newest data
#'
#' @return
#' @export
#'
#' @examples
get_last_batch <- function(
  database,
  collection,
  last_batch,
  periods,
  fields,
  min_effectif,
  rollback_months = 1){

  current_data <- connect_to_database(
    database,
    collection,
    last_batch,
    date_inf = min(periods) %m-% months(rollback_months),
    date_sup = max(periods) %m+% months(1),
    min_effectif = min_effectif,
    fields = fields
    )

  if ("periode" %in% fields && max(current_data$periode) != max(periods))
    warning("Data is missing at actual period !!")

  return(current_data)
}



