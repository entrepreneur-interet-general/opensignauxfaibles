factor_request <- function(
  batch,
  algo,
  siren,
  date_inf,
  date_sup,
  min_effectif,
  fields,
  code_ape,
  export = FALSE
) {

  assert_that(
    !is.null(date_inf) &&
      !is.null(date_sup) &&
      !is.na(date_inf) &&
      !is.na(date_sup),
    msg = "factor_request: les dates spécifiés sont invalides")

  ## Construction de la requête ##
  # Filtrage siren
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

  # Filtrage code APE

  if (!is.null(code_ape)){
    niveau_code_ape <- nchar(code_ape)
    if (niveau_code_ape >= 2){
      field <- "code_ape"
    } else {
      field <- "code_naf"
    }
    match_APE <- paste0('"value.0.', field, '":
                        {"$regex":"^', code_ape, '", "$options":"i"}
                        ')
    match_req <- paste(match_req, match_APE, sep = ", ")
  }

  match_req <- paste0('{"$match":{', match_req, "}}")

  # Unwind du tableau
  unwind_req <- '{"$unwind":{"path": "$value"}}'

  # Filtrage effectif et date
  if (!is.null(siren)){
    eff_req <- ""
  } else {
    eff_req <- paste0(
      '{"$match":{', '"value.effectif":{"$gte":',
      min_effectif,
      '},"value.periode":{
        "$gte": {"$date":"', date_inf, 'T00:00:00Z"},
        "$lt": {"$date":"', date_sup, 'T00:00:00Z"}
}}}')
  }

  # Construction de la projection
  if (is.null(fields)){
    projection_req  <- ""
  } else {
    projection_req  <- paste0('"value.', fields, '":1')
    projection_req  <- paste(projection_req, collapse = ",")
    projection_req  <- paste0('{"$project":{', projection_req, "}}")
  }

  # Construction de l'export
  if (!export){
    export_req <- ""
  } else {
    export_req <- '{"$replaceRoot": { "newRoot": "$value" }},
    { "$project" : {"_id" : 0.0} },
    { "$out" : "to_export"}'
  }

  reqs <- c(
    match_req,
    unwind_req,
    eff_req,
    projection_req,
    export_req)

  requete  <- paste(
    reqs[reqs != ""],
    collapse = ", ")
  requete <- paste0(
    "[",
    requete,
    "]")


  return(requete)

}
