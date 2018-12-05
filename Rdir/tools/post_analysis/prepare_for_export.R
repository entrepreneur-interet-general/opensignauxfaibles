prepare_for_export <- function(donnees, export_fields, database, last_batch, algorithm){

  last_period  <- max(donnees$periode, na.rm = TRUE)
  cat("Préparation à l'export ... \n")
  cat(paste0("Dernière période connue: ",
      last_period, "\n"))

  full_data <- connect_to_database(
    database,
    "Features",
    last_batch,
    date_inf = last_period,
    date_sup = last_period %m+% months(1),
    algo = algorithm,
    min_effectif = 10,
    fields = export_fields[!export_fields %in% c("connu", "diff", "prob")]
  )

  donnees <- donnees %>%
    mutate(siret = as.character(siret)) %>%
    left_join(full_data %>% mutate(siret = as.character(siret)), by = c("siret", "periode")) %>%
    dplyr::mutate(CCSF = date_ccsf) %>%
    dplyr::arrange(dplyr::desc(prob))

  # Report des dernières infos financieres connues

  donnees <- donnees %>%
    mark_known_sirets(name = "sirets_connus.csv") %>%
    select(export_fields)

  all_names <- names(donnees)
  cat("Les variables suivantes sont absentes du dataframe:", "\n")
  cat(export_fields[!(export_fields %in% all_names)])
  export_fields <- export_fields[export_fields %in% all_names]


  #if (is.emp)
  to_export <- donnees %>%
    dplyr::select(one_of(export_fields))
}
