quantile_APE_apply <- function(
  ref_quantile_APE, ...){

  assertthat::assert_that(all(
      c("niveau",
        "code",
        "variable",
        "moy",
        "std") %in%
      names(ref_quantile_APE)
    ))

  ape_levels  <- unique(ref_quantile_APE$niveau)
  variable_names <- unique(ref_quantile_APE$variable)

  aux_fun  <- function(my_data){

    out  <- my_data

    for (i in seq_along(ape_levels)) {
      level <- ape_levels[i]

      if (level == 1) {
        out <- out %>%
          mutate(.target = as.character(code_naf))
      } else {
        out <- out %>%
          mutate(.target = as.character(substr(code_ape, 1, ape_levels)))
      }


      for (i in seq_along(variable_names)){

        new_var  <- paste0(variable_names[i], "_distrib_APE", level)
        out  <- out %>%
          left_join(ref_quantile_APE %>%
            filter(variable == variable_names[i]) %>%
            select(code, moy, std),
          by = c(".target" = "code")) %>%
        mutate(!!new_var := (!!rlang::sym(variable_names[i]) - moy) / std) %>%
        select(-c("moy", "std"))
      }
      out <- out %>%
        select(-.target)
    }
    return(out)
  }

  return(
    lapply(..., aux_fun)
  )
}

