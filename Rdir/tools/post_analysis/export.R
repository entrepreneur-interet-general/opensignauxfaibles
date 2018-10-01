export <- function(data,
                       batch,
                       destination = 'csv',
                       relative_path =  file.path('..','output')) {

  if (!tolower(destination) %in% c('csv', 'mongodb', 'json')) {

    error("Wrong destination argument. Possible destinations are 'csv', 'mongodb' or 'json'")

  } else if(tolower(destination) == 'csv') {

    dirname <- rprojroot::find_rstudio_root_file(relative_path)
    if (! dir.exists(dirname)) error('Directory not found. Check relative_path')

    file_list <- list.files(dirname)
    filename <-  paste0('prediction_batch',batch,'_exported_',Sys.Date())
    n_different <- sum(grepl(paste0('^',filename,'(_v)?[0-9]*.csv$'),file_list))
    if (n_different > 0) filename <- paste0(filename,'_v',n_different+1)
    filename <- paste0(filename,'.csv')

    fullpath <- file.path(dirname,filename)

    write.table(data,
                row.names = F,
                dec = ',',
                sep = ';',
                file = fullpath,
                quote = T,
                append = F
    )



  } else if(tolower(destination) == 'mongodb') {

    error("Export to mongodb not implemented yet !")

  } else if(tolower(destination) == 'json'){

    error("Export to json not implemented yet !")

  }
}

