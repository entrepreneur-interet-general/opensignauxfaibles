context("Test wrappers")

sample   <- get_last_batch(
  database = "test_signauxfaibles",
  collection = "used_in_tests",
  last_batch = "1901_interim",
  periods = as.Date("2019-01-01"),
  fields = c("siret", "periode", "code_naf", "code_ape_niveau2", "code_ape_niveau3"),
  min_effectif = 10,
  rollback_months = 1)

connect_to_h2o()


h2o_test_frame <-
  h2o::as.h2o(datasets::mtcars %>% dplyr::mutate(
    vehicle_name = rownames(mtcars),
    outcome = round(runif(
      nrow(datasets::mtcars)
    ))))

test_that("Typing correction works as expected", {
  h2o_typed_frame <- set_h2o_types(h2o_test_frame)
  expect_equal(h2o::h2o.getTypes(h2o_typed_frame)[[12]], "enum")
  expect_equal(h2o::h2o.getTypes(h2o_typed_frame)[[13]], "enum")
})

test_that("get_last_batch works", {
  expected <- 31
  actual <- nrow(sample)
  expect_equal(nrow(sample), 31)

  actual_min_date <- min(sample$periode)
  actual_max_date <- max(sample$periode)
  expected_min_date <- as.Date("2018-12-01")
  expected_max_date <- as.Date("2019-01-01")

  expect_equal(actual_min_date, expected_min_date)
  expect_equal(actual_max_date, expected_max_date)
})

sample$outcome <- round(runif(n = nrow(sample)))
res <- prepare_frame(
  data_to_prepare = sample,
  test_or_train = "train",
  save_or_load_map = FALSE,
  outcome = "outcome" )

test_that('prepare_frame works', {

  #### Check for te_map training
  expect_equal(names(res), c("data", "te_map"))
  expect_equal(class(res[["data"]]), "H2OFrame")
  expect_true(all(c("TargetEncode_code_ape_niveau2", "TargetEncode_code_ape_niveau3") %in%
                    names(res[["data"]])))

})

test_that("predict_model_works",{
  iris.hex <-  h2o::h2o.uploadFile(
    path = system.file(
    "extdata",
    "iris.csv", package = "h2o"))
  iris.gbm <- h2o.gbm(x = 1:4, y = 5, training_frame = iris.hex, distribution = "multinomial")

  iris.hex <- h2o::h2o.cbind(iris.hex,
                             as.h2o(data.frame(siret = 1:nrow(iris.hex),
                                    periode = rep(as.Date("2018-01-01"), nrow(iris.hex)))
                             )
  )

  actual <- predict_model(model = iris.gbm,
                          iris.hex)

  expect_true("data.frame" %in% class(actual))
  expect_equal(actual$periode[1], as.Date("2018-01-01"))
  expect_true(all(c("prob", "predicted_outcome") %in% names(actual)))
})


test_that("train_model works", {
  expect_error(bar <- train_light_gradient_boosting(
    h2o_train_data = res[["data"]],
    h2o_validation_data = NULL,
    x_fields_model = "TargetEncode_code_ape_niveau2",
    outcome = "outcome",
    save_results = FALSE
  ), NA)
    expect_true("H2OBinomialModel" %in% class(bar))
})

bar <- train_light_gradient_boosting(
    h2o_train_data = res[["data"]],
    h2o_validation_data = NULL,
    x_fields_model = "TargetEncode_code_ape_niveau2",
    outcome = "outcome",
    save_results = FALSE)

test_that("predict_on_last_batch works",{
     expect_error(
       foo <- predict_on_last_batch(
         model = bar,
         database = "test_signauxfaibles",
         collection = "used_in_tests",
         te_map = res[["te_map"]],
         last_batch = "1901_interim",
         periods = as.Date("2019-01-01"),
         min_effectif = 10,
         fields = c("siret", "periode", "code_naf", "code_ape_niveau2", "code_ape_niveau3")),
       NA)

  expect_true(all(c("prob", "last_prob", "diff") %in% names(foo)))
})


test_that("train_light_gradient_boosting works", {
  expect_error(
      full_light_gradient_boosting(
        database = "test_signauxfaibles",
        collection = "used_in_tests",
        periods = as.Date("2019-01-01"),
        last_batch = "1901_interim",
        min_effectif = 10,
        retrain_model = TRUE,
        type = "dataframe"),
      NA)


})
