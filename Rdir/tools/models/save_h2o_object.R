save_h2o_object <- function(
  object,
  object_name,
  extension,
  relative_path =  file.path("..", "output", "model")
  ) {

  assertthat::assert_that(extension %in% c("model", "temap"),
                          msg = "Invalid extension. Extension should be 'model' or 'temap'")

  if (extension == "model"){
    assertthat::assert_that(class(object) == "H2OModel",
                            msg = "Object has wrong class. Class of model object should be H2OFrame.")

    save_function <- h2o.saveModel

  } else if (extension == "temap"){
    assertthat::assert_that(class(object) == "H2OFrame",
                            msg = "Object has wrong class. Class of temap object should be H2OFrame.")
    save_function <- h2o.exportFile
  }

  fullpath <- name_file(
    relative_path,
    file_detail = object_name,
    file_extension = extension
    )

  save_function(object, fullpath)

  return(TRUE)
}
