#' Découper l'échantillon en échantillon d'entraînement, de test et de validation
#'
#' Découpe l'échantillon en échantillon d'entraînement, de validation et de test, de telle sorte à ce que deux entreprises de la même entreprise soient dans le même échantillon.
#' Suréchantillone chaque entreprise en l'observant à des périodes différentes, en prenant aléatoirement \code{frac_periods} periodes.
#'
#' La fraction de l'échantillon de test est égale à 1-\code{frac_val}-\code{frac_train}
#'
#'
#' @param data_sample Un dataframe, avec des champs \code{periode} et \code{siret}.
#' @param frac_train Taille de l'échantillon d'entraînement (fraction du total)
#' @param frac_val Taille de l'échantillon de validation (fraction du total)
#' @param frac_periods Nombre de périodes à échantillonner aléatoirement (fraction du total)
#'
#' @return Liste de trois vecteurs: \code{train}, \code{validation} and \code{test} avec les indices des échantillons
#' @export
#' @examples


split_snapshot_rdm_month <- function(
  data_sample,
  frac_train = 0.60,
  frac_val = 0.20) {

  if (data_sample %>% inherits("tbl_spark")){

    # split <- function(data_sample, split_vect){
    #   browser()
    #   data_sample <- data_sample %>%
    #     sparklyr::sdf_bind_cols(split_vect)
    #   data_sample <- spark_dataframe(data_sample) %>%
    #     sparklyr::invoke("partitionBy", "ss") %>%
    #     register()
    # }
    split <- function(data_sample, split_vect){

      data_sample <- data_sample %>%
        sparklyr::sdf_bind_cols(split_vect)


      res <-  list()
      categories <- pull(split_vect %>% distinct(), ss) %>%
        sort()

      for (i in categories){
        res[[i]] <- data_sample %>% filter(ss == i) %>% select(-ss)
      }
      return (res)
    }
  }




  assertthat::assert_that(
    frac_train > 0 ,
    frac_val > 0,
    frac_train + frac_val <= 1,
    msg = "Fractions must be positive and not exceed 1"
  )

  frac_test = 1 - (frac_train + frac_val)


  data_sample <- data_sample %>%
    select(siret,periode) %>%
    arrange(siret, periode) %>%
    mutate(siren = substr(siret, 1, 9))

  sirens <- data_sample %>% select(siren) %>%
    distinct()
  if (sirens %>% inherits("tbl_spark")){
    sirens <- sirens %>%
      mutate(ss = rand()) %>%
      mutate(ss = ifelse(ss < frac_train, 1, ss)) %>%
      mutate(ss = ifelse(ss < frac_train + frac_val, 2, ss)) %>%
      mutate(ss = ifelse(ss < frac_train + frac_val + frac_test, 3, ss))
  } else {
    sirens <- sirens %>%
      mutate(ss = sample(1:3, size = nrow(sirens), replace = TRUE,
                         prob = c(frac_train, frac_val, frac_test)))
  }

  data_sample <- data_sample %>%
    left_join(y = sirens, by = "siren") %>%
    select(siret, periode, ss)

  result <- setNames(split(data_sample %>% select(-ss), data_sample %>% select(ss)), c("train", "validation", "test"))

  return(result)
}
