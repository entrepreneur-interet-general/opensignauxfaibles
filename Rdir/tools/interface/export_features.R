export_features <- function(algo, batch){
  path <- rprojroot::find_rstudio_root_file(
    "..", "dbmongo", "export", "export.sh"
    )
  system2("bash", args = c(path, batch, algo))
}
