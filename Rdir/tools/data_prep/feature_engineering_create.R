feature_engineering_create <- function(
  train_set,
  quantile_vars = NULL,
  quantile_levels = NULL
  ){

  out  <-  list()
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
