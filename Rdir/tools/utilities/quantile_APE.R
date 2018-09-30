quantile_APE <- function(data, ..., variable_names, levels = 1, noise = 0) {

  assertthat::assert_that(all(c('code_naf_niveau1', 'code_ape') %in% names(data)))
  assertthat::assert_that(all(levels >= 1 & levels <= 5))

  arguments <- c(list(data),list(...))


  # FIX ME: problème de nom dans le mutate_at si une seule variable
  assertthat::assert_that(length(variable_names) > 1)

  ## On retient join les df avec un index selon l'origine

  all_data <- as_tbl_time(bind_rows(arguments, .id = 'index'), index = periode)

  for (i in seq_along(levels)) {
    level = levels[i]


    if (level == 1) {
      all_data <- all_data %>%
        mutate(.target = as.factor(code_naf_niveau1))
    } else {
      all_data <- all_data %>%
        mutate(.target = as.factor(substr(code_ape, 1, levels)))
    }

    aux_fun <- function(x, filter){
      if (all(is.na(x[filter == 1])) || all(filter != 1)) return (NA_real_)
      mecdf = ecdf(x[filter == 1])(x) + (filter == 1) * noise * runif(1,-1,1)
      return(mecdf)
    }

    my_fun <- c("aux_fun")
    names(my_fun) <- paste0("distrib_APE",level)



    for (i in  seq_along(variable_names)){

      varname <- paste0(variable_names[i],"_distrib_APE",level)

      all_data <- all_data %>%
        group_by(.target) %>%
        mutate(!!varname := aux_fun(!!as.name(variable_names[i]), index)) %>%
        ungroup()
    }
  }
  all_data <- all_data %>%
    select(-.target)

  all_data <- all_data %>%
    mutate_if(is.character, as.factor)

  ### Split data
  splitted_data <- split(all_data %>% select(-index), all_data$index)

  return(splitted_data)


}
