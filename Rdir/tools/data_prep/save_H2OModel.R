save_H2OModel <- function(object, path, filename){

  assertthat::assert_that(class(object) == "H2OBinomialModel",
    msg = paste("This function saves a H2OBinomialModel, not a", class(object)))

  h2o.saveModel(object, path)

  # Rename model object, which has the model_id as name at this stage
  file.rename(file.path(path, object@model_id), file.path(path, filename))
}
