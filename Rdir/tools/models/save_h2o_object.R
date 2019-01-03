save_h2o_object <- function(
  object,
  object_name,
  extension,
  relative_path =  file.path("..", "output", "model")
  ) {

  fullpath <- name_file(
    relative_path,
    file_detail = model_name,
    file_extension = extension
    )

  h2o.saveModel(object = model, path = fullpath)

  return(TRUE)
}
