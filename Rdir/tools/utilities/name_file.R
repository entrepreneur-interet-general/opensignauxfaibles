name_file <- function(
  relative_path,
  file_detail,
  file_extension = "",
  full_name = FALSE
  ){

  full_dir_path <- rprojroot::find_rstudio_root_file(relative_path)

  assertthat::assert_that(dir.exists(full_dir_path),
    msg = "Directory not found. Check relative path")

  file_list <- list.files(full_dir_path)


  n_different <- grepl(
    paste0("^", Sys.Date(), "_v[0-9]*_",
      file_detail, "\\.", file_extension, "$"),
    file_list) %>%
  sum()


filename <- paste0(Sys.Date(),
  "_v",
  n_different + 1,
  "_",
  file_detail,
  ".",
  file_extension)

if (full_name){
  full_file_path <- file.path(full_dir_path, filename)

  return(full_file_path)
} else {
  return(filename)
}

}
