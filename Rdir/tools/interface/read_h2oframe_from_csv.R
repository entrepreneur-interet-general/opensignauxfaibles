read_h2oframe_from_csv <- function(){

  path <- rprojroot::find_rstudio_root_file(
    "..", "output", "features", "features.csv"
  )

  # Read csv file
  return(h2o::h2o.importFile(path))
}
