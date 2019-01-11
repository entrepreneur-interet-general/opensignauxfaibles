save_h2o_object <- function(
  object,
  object_name,
  extension,
  relative_path =  file.path("..", "output", "model")
  ) {

  assertthat::assert_that(extension %in% c("model", "temap"),
                          msg = "Invalid extension. Extension should be 'model' or 'temap'")

  if (extension == "model"){
    save_function <- save_H2OModel
  } else if (extension == "temap"){
    save_function <- save_H2OFrame_list
  }

  full_dir_path <- rprojroot::find_rstudio_root_file(relative_path)

  filename <- name_file(
    relative_path,
    file_detail = object_name,
    file_extension = extension
    )

  save_function(object, full_dir_path, filename)

  return(TRUE)
}
