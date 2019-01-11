connect_to_database <- function(
  database,
  collection,
  batch,
  algo = "algo2",
  siren = NULL,
  date_inf = NULL,
  date_sup = NULL,
  min_effectif = 10,
  fields = NULL,
  code_ape = NULL,
  type = "dataframe",
  reexport_csv = TRUE){

  requete <- factor_request(
    batch,
    algo,
    siren,
    date_inf,
    date_sup,
    min_effectif,
    fields,
    code_ape,
    export = (type == "csv")
  )

  assert_that(type %in% c("dataframe", "csv", "iterator"),
    msg = "connect_to_database:
    specification du type d'import/export non valide")

    cat("Connexion à la collection mongodb ...")

    dbconnection <- mongo(
      collection = collection,
      db = database,
      verbose = TRUE,
      url = "mongodb://localhost:27017")
    cat(" Fini.", "\n")

  if (type == "dataframe"){
    cat("Import ...", "\n")

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
    cat("Import de", n_eta, "etablissements issus de",
      n_ent, "entreprises", "\n")

    # Typage
    table_wholesample <- table_wholesample %>%
      mutate_if(is.POSIXct, as.Date)

    if ("numero_compte_urssaf" %in% names(table_wholesample)){
      table_wholesample$numero_compte_urssaf <-
        as.factor(paste(table_wholesample$numero_compte_urssaf))
    }

    if ("code_naf" %in% names(table_wholesample)){
      table_wholesample <- table_wholesample %>%
        mutate(
          code_naf = as.factor(code_naf),
          code_ape_niveau2 = as.factor(substr(code_ape, 1, 2)),
          code_ape_niveau3 = as.factor(substr(code_ape, 1, 3)),
          code_ape_niveau4 = as.factor(substr(code_ape, 1, 4)),
          code_ape = as.factor(code_ape)
          )
    }

    if ("siret" %in% names(table_wholesample)){
      table_wholesample <- table_wholesample %>%
        mutate(
          siret = as.factor(siret)
          )
    }


    if ("siren" %in% names(table_wholesample)){
      table_wholesample <- table_wholesample %>%
        mutate(
          siren = as.factor(siren)
          )
    }

    # Champs manquants
    champs_manquants <- fields[!fields %in% names(table_wholesample)]
    if (length(champs_manquants) >= 1){
      cat("Champs manquants: ")
      cat(champs_manquants, "\n")

      cat("Remplacements par NA", "\n")

      remplacement <- NA * double(length(champs_manquants))
      names(remplacement) <- champs_manquants
      remplacement <- as.data.frame(t(remplacement))
      table_wholesample <- cbind(
        table_wholesample,
        remplacement
        )
    }

    cat(" Fini.", "\n")
    return(table_wholesample)
  } else if (type == "csv"){
    if (reexport_csv){
        export_to_csv(database, algo, batch, fields, min_effectif)
    }
    table_wholesample <- read_h2oframe_from_csv()

    # FIX ME: code dupliqué !
    table_wholesample["code_ape_niveau2"] =
      h2o.substring(table_wholesample["code_ape"], 1, 2)

    table_wholesample["code_ape_niveau3"] =
      h2o.substring(table_wholesample["code_ape"], 1, 3)

    table_wholesample["code_ape_niveau4"] =
      h2o.substring(table_wholesample["code_ape"], 1, 4)

    return(table_wholesample)

  } else if (type == "iterator") {
    cat("not implemented yet")
  }
}
