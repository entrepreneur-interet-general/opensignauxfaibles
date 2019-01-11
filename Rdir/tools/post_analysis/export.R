export <- function(
  donnees,
  batch,
  destination = "csv",
  relative_path =  file.path("..", "output")) {


  assertthat::assert_that(tolower(destination) %in% c("csv", "mongodb", "json"),
    msg = "Wrong export destination argument.
    Possible destinations are 'csv', 'mongodb' or 'json'")

  if (tolower(destination) == "csv") {


    fullpath <- name_file(
      relative_path,
      file_detail = paste0("detection", batch),
      file_extension = "csv",
      full_name = TRUE
      )

    write.table(donnees,
      row.names = F,
      dec = ",",
      sep = ";",
      file = fullpath,
      quote = T,
      append = F
      )



  } else if (tolower(destination) == "mongodb") {

    error("Export to mongodb not implemented yet !")

  } else if (tolower(destination) == "json"){

    error("Export to json not implemented yet !")

  }
}
