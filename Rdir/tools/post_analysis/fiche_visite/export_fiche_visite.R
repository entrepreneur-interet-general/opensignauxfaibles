export_fiche_visite <- function(
  donnees,
  sirets){
  save(donnees, file = "/home/pierre/Documents/opensignauxfaibles/Rdir/RData/donnees.RData")

  for (i in seq_along(sirets)){

    raison_sociale <- unique(donnees[donnees$siret == sirets[i], "raison_sociale"])

    rmarkdown::render("/home/pierre/Documents/opensignauxfaibles/Rdir/tools/post_analysis/fiche_visite/fiche_visite.Rmd",
      params  = list(
        siret = sirets[i],
        raison_sociale = raison_sociale
      ),
      output_file = paste0("/home/pierre/Documents/opensignauxfaibles/Rdir/RData/Fiches/Fiche_visite_", raison_sociale, ".pdf")
    )
    }
}
