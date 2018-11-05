shapley_plot <- function(siret, my_data, model) {
  pred <- function(model, newdata)  {
    results <- as.data.frame(h2o::h2o.predict(model, h2o::as.h2o(newdata)))
    return(results[[3L]])
  }


  x_medium <- c("montant_part_patronale",
                "ratio_dette",
                "ratio_dette_moy12m",
                "etat_proc_collective_num",
                "TargetEncode_code_ape_niveau3",
                "cotisation",
                "frais_financier_distrib_APE1",
                "taux_marge_distrib_APE1",
                "montant_part_patronale_variation_3",
                "ratio_liquidite_reduite_distrib_APE1",
                "dette_fiscale",
                "ratio_delai_client_distrib_APE1",
                "montant_part_patronale_variation_1",
                "montant_part_patronale_variation_2",
                "ratio_rend_capitaux_propres",
                "taux_marge",
                "poids_frng_distrib_APE1",
                "delai_fournisseur_distrib_APE1",
                "taux_rotation_stocks_distrib_APE1",
                "effectif",
                "ratio_rentabilite_nette_distrib_APE1",
                "ratio_export_distrib_APE1",
                "TargetEncode_code_ape_niveau2",
                "effectif_diff_moy12",
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


  x_medium_names <- c(
    "Montant part patronale",
    "Ratio dette / cotisation",
    "Moyenne dette/cotisation (12 mois)",
    "Procédure collective en cours",
    "Taux de défaillance dans le secteur d'activité (code APE 3)",
    "Cotisations URSSAF",
    "Comparaison des frais financiers par code NAF",
    "Comparaison du taux de marge par code NAF",
    "Montant part patronale 3 mois en arrière",
    "Comparaison des liquidités réduites par code NAF",
    "Dette fiscale et sociale",
    "Comparaison du délai client par code NAF",
    "Montant part patronale 1 mois en arrière",
    "Montant part patronale 2 mois en arrière",
    "Rendement des capitaux propres",
    "Taux de marge",
    "Poids du frng",
    "Comparaison du délai fournisseur par code NAF",
    "Comparaison du taux de rotation des stocks par code NAF",
    "Effectif salarié",
    "Comparaison de la rentabilité nette par code NAF",
    "ratio_export_distrib_APE1",
    "Taux de défaillance dans le secteur d'activité (code APE 2)",
    "Variation mensuelle d'effectif moyenne sur 12 mois",
    "Montant de la part ouvrière",
    "Comparaison financier court terme par code NAF",
    #"effectif_entreprise",
    "Age de l'entreprise",
    "Ratio des liquidités réduites",
    "Comparaison de la productivité par code NAF",
    "Frais financiers",
    "Frais financiers court terme",
    "Délai client",
    "Taux de défillance dans le secteur d'activité (code NAF)",
    "Résultat net consolidé",
    "Taux de rotation des stocks",
    "Nombre d'établissements secondaires",
    "Nombre d'établissements connus",
    "Chiffre d'affaire",
    "Chiffre d'affaire net lié aux exportations",
    "Décroissance de la dette pendant un délai URSSAF",
    "Comparaison de la marge opérationnelle par code NAF",
    "Poids du frng"
  )

  names(x_medium_names) <- x_medium


  features  <- my_data[, x_medium]

  response <- 2 - as.numeric(as.factor(as.vector(my_data$outcome)))

  predictor.xgb <- Predictor$new(
    model = model,
    data = features,
    y = response,
    predict.fun = pred,
    class = "classification"
  )
  etablissement <- my_data %>%
    filter(siret == mon_siret) %>%
    filter(periode == max(periode))
  etablissement <- etablissement[, x_medium]
  shap.xgb <- iml::Shapley$new(predictor.xgb, x.interest = etablissement)

  shap_plot <- shap.xgb %>%
    plot()
  thresh <- 5e-3
  to_remove <- abs(shap_plot$data[, "phi"]) < thresh
  shap_plot$data <- shap_plot$data[!to_remove, ]

  # Changing for more informative names
  current_names <- levels(shap_plot$data$feature.value)
  chopped_names <- sapply(
    stringr::str_split(current_names, pattern = "="),
    FUN = function(x) x[1]
    )
  new_levels <- x_medium_names[chopped_names]

  levels(shap_plot$data$feature.value) <- new_levels

  # labels
  shap_plot$labels$x <- ""
  return(shap_plot)
}
