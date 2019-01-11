load_H2OModel <- function(path, filename){

  object <- h2o.loadModel(file.path(path, filename))

  assertthat::assert_that(class(object) == "H2OBinomialModel",
    msg = paste("This function loads a H2OBinomialModel, not a", class(object)))

  return(object)


}
