load_h2o_object <- function(
    name,
    extension,
    relative_path = file.path("..", "output", "model"),
    last = TRUE,
    file_name = ""
    ){


  assertthat::assert_that(extension %in% c("model", "temap"),
    msg = 'Unsupported extension. Supported extensions are "model" and "temap"')

  if (extension == "model") {
    load_function <- load_H2OModel
  } else if (extension == "temap") {
    load_function <- load_H2OFrame_list
  }

  full_dir_path <- rprojroot::find_rstudio_root_file(relative_path)
  assertthat::assert_that(dir.exists(full_dir_path),
    msg = "Directory not found. Check relative path")

  if (last){
    file_candidates <- list.files(full_dir_path) %>%
      grep(pattern = paste0(name, ".", extension), value = TRUE)

    assertthat::assert_that(length(file_candidates) > 0,
      msg = "No such file, please check name and extension")

    file_name <-  file_candidates %>%
      sort(decreasing = TRUE) %>%
      .[1]
  }


  full_path <- file.path(full_dir_path, file_name)

  assertthat::assert_that(file.exists(full_path),
      msg = "No such file, please check file_name")

  res <- load_function(full_dir_path, file_name)

  return(res)
}
