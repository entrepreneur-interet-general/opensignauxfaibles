mark_known_sirets <- function(df, name){
  sirets <-readLines(find_rstudio_root_file('..','data-raw',name))
  sirens <- substr(sirets,1,9)

  df <- df %>%
    mutate(connu = as.numeric(substr(siret,1,9) %in% sirens))

  return(df)
}
