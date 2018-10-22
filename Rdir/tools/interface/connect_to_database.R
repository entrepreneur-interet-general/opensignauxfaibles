connect_to_database <- function(
  collection,
  batch,
  algo = "algo2",
  siren = NULL,
  date_inf = NULL,
  date_sup = NULL,
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

  # Filtrage effectif et date
  if (!is.null(siren)){
    eff_req <- ""
  } else {
    eff_req <- paste0(
      '{"$match":{', '"value.effectif":{"$gte":',
      min_effectif,
      '},"value.periode":{"$gte": {"$date":"', date_inf, 'T00:00:00Z"}, "$lt": {"$date":"', date_sup, 'T00:00:00Z"}}}}')
  }

  cat(eff_req, "\n")

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
  return(table_wholesample)
}
