connect_to_database <- function(
  collection,
  batch,
  algo = "algo2",
  siren = NULL,
  min_effectif = 10,
  fields = NULL){

  cat("Connexion à la collection mongodb ...")
  dbconnection <- mongo(
    collection = collection,
    db = "opensignauxfaibles",
    verbose = TRUE,
    url = "mongodb://localhost:27017")
  cat(" Fini.", "\n")

  # Construction du match
  match_req  <- paste0('"_id.batch":"', batch, '","_id.algo":"', algo, '"')
  if (!is.null(siren)){
    match_siren  <- c()
    for (i in seq_along(siren)){
      match_siren  <- c(
        match_siren,
        paste0('{"_id.siren":"', siren[i], '"}')
        )
    }
    match_siren <- paste0('"$or":[', paste(match_siren, collapse = ","), "]")
    match_req <- paste(match_req, match_siren, sep = ", ")
  }
  match_req <- paste0('{"$match":{', match_req, "}}")

  # Construction de l'unwind
  unwind_req <- '{"$unwind":{"path": "$value"}}'

  # Filtrage effectif
  if (!is.null(siren)){
    eff_req <- ""
  } else {
    eff_req <- paste0(
      '{"$match":{"value.effectif":{"$gte":',
      min_effectif,
      '}}}')
  }

  # Construction de la projection
  if (is.null(fields)){
    projection_req  <- ""
  } else {
    projection_req  <- paste0('"value.',fields,'":1')
    projection_req  <- paste(projection_req, collapse = ",")
    projection_req  <- paste0('{"$project":{', projection_req, '}}')
  }

  reqs <- c(
    match_req,
    unwind_req,
    eff_req,
    projection_req)

  requete  <- paste(
    reqs[reqs != ""],
    collapse = ", ")
  requete <- paste0(
    "[",
    requete,
    "]")


  cat("Import ...")

  donnees <- dbconnection$aggregate(requete)$value

  cat(" Fini.", "\n")

  assertthat::assert_that(nrow(donnees) > 0,
    msg = "La requête ne retourne aucun résultat")
  assertthat::assert_that(
    all(c("periode", "siret") %in% names(donnees))
    )
  assertthat::assert_that(
    anyDuplicated(donnees %>% select(siret, periode)) == 0
    )

  table_wholesample <- donnees %>%
    mutate(periode = as.Date(periode)) %>%
    arrange(periode) %>%
    tibbletime::as_tbl_time(periode)

  n_eta <- table_wholesample$siret %>%
    n_distinct()
  n_ent <- table_wholesample$siret %>%
    str_sub(1, 9) %>%
    n_distinct()
  cat("Import de", n_eta, "etablissements issus de", n_ent, "entreprises", "\n")


  cat(" Fini.", "\n")

#  cat("Import des libellé NAF niveaux 1 et 5 ...")
#
#  libelle_naf <- readxl::read_excel(
#    path = rprojroot::find_rstudio_root_file(file.path("..", "data-raw", "naf", "naf2008_5_niveaux.xls")),
#    sheet = "naf2008_5_niveaux",
#    skip = 1,
#    col_names = c("code_naf_niveau5", "code_naf_niveau4", "code_naf_niveau3", "code_naf_niveau2", "code_naf_niveau1")
#    ) %>%
#  dplyr::select(code_naf_niveau5, code_naf_niveau1) %>%
#  dplyr::left_join(
#    y = readxl::read_excel(
#      path = rprojroot::find_rstudio_root_file(file.path("..", "data-raw", "naf", "naf2008_liste_n5.xls")),
#      sheet = "Feuil1",
#      skip = 3,
#      col_names = c("code_naf_niveau5", "libelle_naf_niveau5")
#      ),
#    by = "code_naf_niveau5"
#    ) %>%
#  dplyr::left_join(
#    y = readxl::read_excel(
#      path = rprojroot::find_rstudio_root_file(file.path("..", "data-raw", "naf", "naf2008_liste_n1.xls")),
#      sheet = "Feuil1",
#      skip = 3,
#      col_names = c("code_naf_niveau1", "libelle_naf_niveau1")
#      ),
#    by = "code_naf_niveau1"
#    )  %>%
#  dplyr::mutate(
#    code_naf_niveau5 = stringr::str_replace(
#      string = code_naf_niveau5,
#      pattern = "([[:digit:]]{2})\\.([[:digit:]]{2}[[:upper:]]{1})",
#      replacement = "\\1\\2")
#    )
#
#  table_wholesample <- table_wholesample %>%
#    left_join(libelle_naf, by = c("code_ape" = "code_naf_niveau5")) %>%
#    mutate(code_ape = as.factor(code_ape))
#
#  cat(" Fini.", "\n")

  return(table_wholesample)
}
