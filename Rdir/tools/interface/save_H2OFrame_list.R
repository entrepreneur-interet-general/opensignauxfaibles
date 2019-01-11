save_H2OFrame_list <- function(object, path, filename){

  assertthat::assert_that(class(object) == "list",
    msg = paste("This function saves a list, not a", class(object)))

  # Convert H2OFrames to dataframes
  res <- lapply(object, as.data.frame)

  saveRDS(res, file.path(path, filename))
}
