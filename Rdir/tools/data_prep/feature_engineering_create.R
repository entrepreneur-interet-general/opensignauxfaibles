feature_engineering_create <- function(
  train_set,
  quantile_vars = NULL,
  quantile_levels = NULL
  ){

  out  <-  list()
  ratios_financiers <- c(
    "taux_marge",
    "delai_fournisseur",
    "poids_frng",
    "financier_court_terme",
    "frais_financier",
    "dette_fiscale"
    )
  #####################
  ## APE COMPARISONS ##
  #####################

  if (is.null(quantile_vars)) {
    quantile_vars <- c(ratios_financiers, "effectif")

    cat("Taking default value for quantile_APE variables", "\n")
  }

  if (is.null(quantile_levels)) {
    quantile_levels <-  c(1, 2)
    cat("Taking default value for quantile_APE levels", "\n")
  }

  assertthat::assert_that(all(quantile_vars %in% names(train_set)))

  out[["quantile"]] <- quantile_APE_create(
    train_set,
    variable_names = quantile_vars,
    ape_levels = quantile_levels)

  return(out)
}
