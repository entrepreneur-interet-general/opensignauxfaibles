load_H2OFrame_list <- function(path, filename){

  object <- readRDS(file.path(path, filename))

  assertthat::assert_that(class(object) == "list",
    msg = paste("This function loads a list, not a", class(object)))

  # Convert dataframes to H2OFrames
  res <- lapply(object, as.h2o)

  return(res)
}
