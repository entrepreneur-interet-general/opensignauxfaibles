add_past_trends <- function(
  data,
  variables,
  lookback_months,
  type = 'lag') {

    assertthat::assert_that(type %in% c('lag','slope','rate','mean_unique'))

   #  check_zeros <- function(var, data) {
   #    return(sum(data[var] == 0, na.rm = TRUE) > 0)
   #  }

   #  with_zeros <-
   #    map_int(.x = variables, .f = check_zeros, data = data)
   #  names(with_zeros) <- variables


    assertthat::assert_that(all(c('siret', 'periode', variables) %in% names(data)))

    grid <-
      expand.grid(
        variables = variables,
        lookback = lookback_months,
        stringsAsFactors = FALSE
      )
    variables <- grid$variables
    lookback_months <- grid$lookback

    data_nested <- data %>%
      arrange(siret, periode) %>%
      group_by(siret) %>%
      nest()

    if (type == 'slope') {
    #   aux_slope <- function(x, y) {
    #     n <- length(x)
    #     sdx <- sqrt((n ^ 2 - 1) / 12) * sqrt(n / (n - 1))
    #     slope <- cor(x, y) * sd(y) / sdx
    #     return(slope)
    #   }
    #   #### Coefficient de variation de la régression ####
    #   aux_past <- function(data, variable, last_n, with_zero) {
    #     y <- data[[variable]]
    #     slopes <- array(dim = length(y))
    #     if (last_n <= length(y))
    #       for (i in last_n:length(y)) {
    #         sub_y <- tail(y[1:i], n = last_n)
    #         x <- 1:length(sub_y)
    #         slopes[i] <- aux_slope(x, sub_y)
    #       }

    #     assertthat::assert_that(length(slopes) == length(y))

    #     return(slopes)
    # }

    } else if (type == 'rate') {
      #### Taux de variation ####
     # aux_past <-  function(data, variable, last_n, with_zero) {
     #   y <- data[[variable]]
     #   y_lag <- lag(y, last_n)
     #   if (with_zero) {
     #     tv <- sign(y - y_lag[1:length(y)])
     #   } else {
     #     tv <- (y - y_lag[1:length(y)]) / y
     #   }
     #   return(tv)
     # }

    } else if (type == 'lag')  {
      #### Valeur reportee ####
      aux_past <- function(data, variable, last_n){
        y_lag <- lag(y,last_n)
        return(y_lag[1:length(y)])
      }

    } else if(type == 'mean_unique'){
      ### Moyenne de la difference après suppression de doublons ##
      aux_past <- function(data, variable, last_n){
        res <- sapply(1:length(y), FUN = function(x){
          out <- unique(tail(y[1:x],last_n+1))
        })
      }

    }

    #cat('/!\ FIX ME: How are slopes on 6 months computed with only  3 available months? + Sensitive to holes','\n')

    for (i in seq_along(lookback_months)) {
      data_nested[[paste0(variables[i], '_variation_', lookback_months[i])]] <-
        purrr::map(
          data_nested$data,
          aux_past,
          variable = variables[i],
          last_n = lookback_months[i]
        )
    }


    data <- data_nested %>%
      unnest()

    return(data)

  }
