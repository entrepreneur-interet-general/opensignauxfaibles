shapley_plot <- function(
  mes_sirets,
  my_data,
  model,
  batch,
  dir_out = find_rstudio_root_file("..", "output", "shapley", batch)
  ) {

  assert_that(inherits(my_data, "data.frame"))
  max_periode <- max(my_data$periode)

  h2o::h2o.no_progress()
  pred <- function(model, newdata)  {
    results <- as.data.frame(h2o::h2o.predict(model, h2o::as.h2o(newdata)))
    return(results[[3L]])
  }


  x_medium <- c(
    "montant_part_patronale",
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
    "benefice_ou_perte",
    "taux_rotation_stocks",
    "nombre_etab_secondaire",
    "nbr_etablissements_connus",
    "CA",
    "chiffre_affaires_net_lie_aux_exportations",
    "ratio_dette_delai",
    "ratio_marge_operationnelle_distrib_APE1",
    "poids_frng"
    )



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

  # x_medium_names_test <- c(
  #   "Dette_URSSAF",
  #   "Dette_URSSAF",
  #   "Dette_URSSAF",
  #   "Procédure_collective_en_cours",
  #   "Taux_de_defaillance_dans_le_secteur_d_activité",
  #   "Taille_de_l_entreprise",
  #   "Endettement_par_code_NAF",
  #   "Solvabilité_par_code_NAF",
  #   "Dette_URSSAF",
  #   "Solvabilité_par_code_NAF",
  #   "Dette_fiscale_et_sociale",
  #   "Délais_de_paiement_par_code_NAF",
  #   "Dette_URSSAF",
  #   "Dette_URSSAF",
  #   "Rentabilité",
  #   "Rentabilité",
  #   "Robustesse",
  #   "Délais_de_paiement_par_code_NAF",
  #   "Comparaison_du_taux_de_rotation_des_stocks_par_code_NAF",
  #   "Taille_de_l_entreprise",
  #   "Rentabilité_par_code_NAF",
  #   "Ratio_d_exportation",
  #   "Taux_de_defaillance_dans_le_secteur_d_activité",
  #   "Variation_de_taille",
  #   "Dette_URSSAF",
  #   "Solvabilité_par_code_NAF",
  #   #"effectif_entreprise",
  #   "Age_de_l_entreprise",
  #   "Solvabilité",
  #   "Rentabilité_par_code_NAF",
  #   "Endettement",
  #   "Solvabilité",
  #   "Délais",
  #   "Taux_de_defaillance_dans_le_secteur_d_activité",
  #   "Résultat_net_consolidé",
  #   "Taux_de_rotation_des_stocks",
  #   "Taille_de_l_entreprise",
  #   "Taille_de_l_entreprise",
  #   "Taille_de_l_entreprise",
  #   "Taille_de_l_entreprise",
  #   "Délai_URSSAF",
  #   "Rentabilité_par_code_NAF",
  #   "Robustesse"
  #   )

  # names(x_medium_names_test) <- x_medium
  names(x_medium_names) <- x_medium

  features  <- my_data[, x_medium]

  response <- pred(model, features)

  predictor.xgb <- iml::Predictor$new(
    model = model,
    data = features,
    y = response,
    predict.fun = pred,
    class = "classification"
    )

  for (i in seq_along(mes_sirets)){
    etablissement <- my_data %>%
      filter(siret == mes_sirets[i]) %>%
      filter(periode == max_periode)
    etablissement <- etablissement[, x_medium]
    shap.xgb <- iml::Shapley$new(predictor.xgb, x.interest = etablissement)

    shap_plot <- shap.xgb %>%
      plot()


    # Changing for more informative names
    # names <- levels(shap_plot$data$feature)

    shap_plot$data <- shap_plot$data %>%
      mutate(category = unname(x_medium_names[feature])) %>%
      group_by(category) %>%
      summarize(
        feature = unique(category),
        phi = sum(phi),
        phi.var = sum(phi.var)
        ) %>%
      mutate(feature.value = as.factor(category))

    thresh <- 5e-3
    to_remove <- abs(shap_plot$data[, "phi"]) < thresh
    shap_plot$data <- shap_plot$data[!to_remove, ]

    # labels
    shap_plot$labels$x <- ""
    dir.create(dir_out, showWarnings = FALSE)
    ggsave(filename = paste0(mes_sirets[i], "_", max_periode,".png"),
      plot = shap_plot,
      path = dir_out)
  }
  h2o::h2o.show_progress()
  return(shap_plot)
}
