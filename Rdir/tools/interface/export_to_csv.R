export_to_csv <- function(database, algo, batch, fields, min_effectif){
  path_1 <- rprojroot::find_rstudio_root_file(
    "..", "dbmongo", "export", "export.sh"
    )

  path_2 <- rprojroot::find_rstudio_root_file(
    "..", "dbmongo", "export", "export_fields.txt"
  )

  ## Write fields to file path_1
  write(fields, path_2, append = FALSE, sep = "\n")

  ##

  browser()
  # FIX ME: ignore complètement min_effectif !! (par défaut, min_effectif = 10)
  system2("bash", args = c(path_1, database, algo, batch, min_effectif))
}
