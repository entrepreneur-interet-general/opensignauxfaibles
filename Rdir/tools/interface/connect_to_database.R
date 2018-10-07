connect_to_database <- function(
  collection,
  batch,
  algo = "algo2",
  siren = NULL,
  min_effectif = 20){
  cat("Connexion à la collection mongodb ...")
  dbconnection <- mongo(
    collection = collection,
    db = "opensignauxfaibles",
    verbose = TRUE,
    url = "mongodb://localhost:27017")
  cat(" Fini.", "\n")

  cat("Import ...")
  if (is.null(siren)){
    requete  <- paste0('[{"$match":{"_id.batch":"', batch, '",
      "_id.algo":"', algo, '"}},
      {"$unwind":{"path": "$value"}}]')
} else {
  match_siren  <- '"$or":['
  for (i in seq_along(siren)){
    match_siren  <- paste0(
      match_siren,
      '{"_id.siren":"',
      siren[i],
      '"},'
      )
  }
  match_siren <- paste0(
    substr(match_siren, 1, nchar(match_siren) - 1),
    "]")

  requete  <- paste0('[{"$match":{"_id.batch":"', batch, '",
    "_id.algo":"', algo, '",',
    match_siren, '}},
    {"$unwind":{"path": "$value"}}]')
}

  data <- dbconnection$aggregate(requete)$value

  if (algo == "algo1") {
    # FIX ME: rien à faire ici
    cat("FIX ME ... feature engineering a effectuer ailleurs\n")
    data <- data %>%
      mutate(
        cut_growthrate = forcats::fct_relevel(
          cut_growthrate,
          c("moins_de_20p", "moins_20_a_5p", "stable",
            "plus_5_a_20p", "plus_20p", "manquant")),
        cut_effectif = forcats::fct_relevel(cut_effectif),
        outcome_0_12 = forcats::fct_relevel(outcome_0_12,
          c("default", "non_default")
          )
        )
  }

  cat(" Fini.", "\n")
  assertthat::assert_that(nrow(data)>0)
  assertthat::assert_that(all(c("periode", "siret", "effectif", "code_ape") %in% names(data)))
  assertthat::assert_that(anyDuplicated(data %>% select(siret, periode)) == 0)

  table_wholesample <- data %>%
    mutate(periode = as.Date(periode)) %>%
    arrange(periode) %>%
    tibbletime::as_tbl_time(periode)

  n_eta <- table_wholesample$siret %>% n_distinct()
  n_ent <- table_wholesample$siret %>% str_sub(1, 9) %>% n_distinct()
  cat("Import de", n_eta, "etablissements issus de", n_ent, "entreprises", "\n")

  if (is.null(siren)){
    cat("Filtrage des établissements avec des effectifs de moins de", min_effectif, "...")
    table_wholesample  <- table_wholesample %>%
      group_by(siret) %>%
      mutate(toKeep = max(effectif)> min_effectif) %>%
      filter(toKeep) %>%
      select(-toKeep)
  }

  if (algo == "algo1") {
    table_wholesample <- table_wholesample %>%
      dplyr::mutate(proc_collective = if_else(any(!is.na(date_defaillance)), max(date_defaillance, na.rm =TRUE), as.POSIXct(NA)))
  }

  table_wholesample <- table_wholesample %>%
    ungroup()
  cat(" Fini.", "\n")
  n_eta_filtered <- table_wholesample$siret %>% n_distinct()
  n_ent_filtered <-  table_wholesample$siret %>% str_sub(1, 9) %>% n_distinct()
  cat(n_eta_filtered, "etablissements restants, de", n_ent_filtered, "entreprises", "\n")

  cat("Import des libellé NAF niveaux 1 et 5 ...")

  libelle_naf <- readxl::read_excel(
    path = rprojroot::find_rstudio_root_file(file.path("..", "data-raw", "naf", "naf2008_5_niveaux.xls")),
    sheet = "naf2008_5_niveaux",
    skip = 1,
    col_names = c("code_naf_niveau5", "code_naf_niveau4", "code_naf_niveau3", "code_naf_niveau2", "code_naf_niveau1")
    ) %>%
  dplyr::select(code_naf_niveau5, code_naf_niveau1) %>%
  dplyr::left_join(
    y = readxl::read_excel(
      path = rprojroot::find_rstudio_root_file(file.path("..", "data-raw", "naf", "naf2008_liste_n5.xls")),
      sheet = "Feuil1",
      skip = 3,
      col_names = c("code_naf_niveau5", "libelle_naf_niveau5")
      ),
    by = "code_naf_niveau5"
    ) %>%
  dplyr::left_join(
    y = readxl::read_excel(
      path = rprojroot::find_rstudio_root_file(file.path("..", "data-raw", "naf", "naf2008_liste_n1.xls")),
      sheet = "Feuil1",
      skip = 3,
      col_names = c("code_naf_niveau1", "libelle_naf_niveau1")
      ),
    by = "code_naf_niveau1"
    )  %>%
  dplyr::mutate(
    code_naf_niveau5 = stringr::str_replace(
      string = code_naf_niveau5,
      pattern = "([[:digit:]]{2})\\.([[:digit:]]{2}[[:upper:]]{1})",
      replacement = "\\1\\2")
    )

  table_wholesample <- table_wholesample %>%
    left_join(libelle_naf, by = c("code_ape"= "code_naf_niveau5")) %>%
    mutate(code_ape = as.factor(code_ape))

  cat(" Fini.", "\n")

  return(table_wholesample)
}

