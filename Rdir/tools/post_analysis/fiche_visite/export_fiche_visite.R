export_fiche_visite <- function(
  sirets,
  database,
  batch,
  with_urssaf,
  folder = batch){

  field_names <- c(
    "siret",
    "siren",
    "periode",
    "raison_sociale",
    "effectif",
    "cotisation_moy12m",
    "montant_part_patronale",
    "montant_part_ouvriere",
    "apart_heures_consommees",
    "apart_heures_autorisees",
    "code_naf",
    "code_ape",
    "libelle_naf",
    "libelle_ape2",
    "libelle_ape3",
    "libelle_ape4",
    "libelle_ape5",
    "CA",
    "resultat_net_consolide",
    "resultat_expl",
    "poids_frng",
    "taux_marge",
    "delai_fournisseur",
    "dette_fiscale",
    "financier_court_terme",
    "frais_financier",
    "age",
    "departement",
    "region",
    "CA_past_1",
    "resultat_expl_past_1",
    "resultat_net_consolide_past_1",
    "ratio_apart"
  )



  for (i in seq_along(sirets)){

    one_company <- connect_to_database(
      database = database,
      collection = "Features",
      batch = batch,
      date_inf = "2014-01-01",
      date_sup = "2999-12-31",
      algo = "algo2",
      siren = substr(sirets[i], 1, 9),
      min_effectif = 0,
      fields = field_names
    )

    cat("Filtrage par code ape ", substring(one_company$code_ape[1], 1, 3),"\n")

    donnees <- connect_to_database(
      database = database,
      collection = "Features",
      batch = batch,
      date_inf = "2014-01-01",
      date_sup = "2999-12-31",
      algo = "algo2",
      siren = substr(sirets[i], 1, 9),
      min_effectif = 0,
      fields = field_names,
      code_ape = substring(one_company$code_ape[1], 1, 3)
    )


    save(donnees, file = "/home/pierre/Documents/opensignauxfaibles/Rdir/RData/donnees.RData")

    raison_sociale <- unique(donnees[donnees$siret == sirets[i], "raison_sociale"])

    rmarkdown::render("/home/pierre/Documents/opensignauxfaibles/Rdir/tools/post_analysis/fiche_visite/fiche_visite.Rmd",
      params  = list(
        siret = sirets[i],
        batch = batch,
        raison_sociale = raison_sociale,
        with_urssaf = with_urssaf
      ),
      output_file = paste0("/home/pierre/Documents/opensignauxfaibles/output/Fiches/", folder, "/Fiche_visite_", raison_sociale, ".pdf"),
      clean = TRUE
    )
    }
}
