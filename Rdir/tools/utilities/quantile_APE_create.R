quantile_APE_create <- function(
  my_data,
  variable_names,
  ape_levels = 1
){

  assertthat::assert_that(all(c("code_naf", "code_ape") %in%
                                names(my_data)))

  assertthat::assert_that(all(ape_levels >= 1 & ape_levels <= 5))


  out <- data.frame(
    niveau = numeric(0),
    code = character(0),
    variable = character(0),
    moy = numeric(0),
    std = numeric(0)
  )

  for (i in seq_along(ape_levels)) {
    level <- ape_levels[i]

    if (level == 1) {
      my_data <- my_data %>%
        mutate(.target = as.character(code_naf))
    } else {
      my_data <- my_data %>%
        mutate(.target = as.character(substr(code_ape, 1, ape_levels)))
    }

    for (i in seq_along(variable_names)){
      out_add  <- my_data %>%
        group_by(.target) %>%
        summarize(
          niveau = level,
          code = unique(.target),
          variable = variable_names[i],
          moy = mean(!!rlang::sym(variable_names[i]), na.rm = TRUE),
          std = sd(!!rlang::sym(variable_names[i]), na.rm = TRUE)
        )
      out <- rbind(out, out_add)
    }
  }

  out <- out %>%
    select(-c(".target"))

  out[!is.finite(out$moy), 'moy'] <- NA
  out[!is.finite(out$std), 'std'] <- NA
  return(out)
}
