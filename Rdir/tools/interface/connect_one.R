connect_one  <- function(collection, siret, batch, algo="algo2"){
	cat('Connexion à la collection mongodb ...')
	dbconnection <- mongo(collection = collection, db = 'opensignauxfaibles', verbose = TRUE, url = 'mongodb://localhost:27017')
	cat(' Fini.','\n')

	cat('Import ...')

	siren  <- substr(siret,1,9)

	# FIX ME: traiter plusieurs sirets/sirens
	my_data <- dbconnection$aggregate(paste0('[{"$match":{"_id.batch":"',batch,'", "_id.algo":"',algo,'", "_id.siren":"', siren, '"}},{"$unwind":{"path": "$value"}}]'))$value

	cat(' Fini.','\n')

	if (nrow(my_data) == 0)
		warning("No data has been found")

	table_siret <- my_data %>%
		filter(siret %in% siret) %>%
		mutate(periode = as.Date(periode)) %>%
		arrange(periode) %>%
		tibbletime::as_tbl_time(periode)


	return(table_siret)
}
