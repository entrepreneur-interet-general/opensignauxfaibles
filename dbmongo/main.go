package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {
	// Lancer Rserve en background

	InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)

	loadConfig()

	go r()

	r := gin.Default()
	r.Use(DB())
	// r.Use(Kanboard())
	// FIXME: configurer correctement CORS
	r.Use(cors.Default())

	r.Use(static.Serve("/", static.LocalFile("static/", true)))

	api := r.Group("api")
	{
		api.OPTIONS("/auth", auth)

		api.GET("/purge", purge)

		api.GET("/kanboard/get/projects", listProjects)
		api.GET("/kanboard/get/tasks", getKBTasks)

		api.PUT("/admin/batch/:batchID", registerNewBatch)
		api.POST("admin/batch", updateBatch)
		api.GET("/admin/batch", listBatch)
		api.DELETE("/admin/batch", dropBatch)

		api.GET("/admin/files", adminFiles)
		api.POST("/admin/attach", attachFileBatch)
		api.GET("/admin/types", listTypes)
		api.GET("/admin/clone/:to", cloneDB)
		api.GET("/admin/features", adminFeature)

		api.GET("/import/:batch", importBatch)

		api.GET("/compact/etablissement/:siret", compactEtablissement)
		api.GET("/compact/etablissement", compactEtablissement)
		api.GET("/compact/entreprise/:siren", compactEntreprise)
		api.GET("/compact/entreprise", compactEntreprise)

		api.GET("/reduce/:algo/:batch/:siret", reduce)
		api.GET("/reduce/:algo/:batch", reduce)

		api.POST("/R/algo1", algo1)

		api.GET("/data/prediction/:batch/:algo/:page", predictionBrowse)
		api.GET("/debug/", debug)
	}
	bind := viper.GetString("APP_BIND")
	r.Run(bind)
}

func loadConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("toml")
	viper.AddConfigPath("/etc/opensignauxfaibles")
	viper.AddConfigPath("$HOME/.opensignauxfaibles")
	viper.AddConfigPath(".")
	viper.SetDefault("APP_BIND", ":3000")
	viper.SetDefault("APP_DATA", "$HOME/data-raw/")
	viper.SetDefault("DB_HOST", "127.0.0.1")
	viper.SetDefault("DB_PORT", "27017")
	viper.SetDefault("DB", "opensignauxfaibles")
	viper.SetDefault("JWT_SECRET", "One might change this because one day it will not be sufficient")
	viper.SetDefault("KANBOARD_ENDPOINT", "http://localhost/kanboard/jsonrpc.php")
	viper.SetDefault("KANBOARD_USERNAME", "admin")
	viper.SetDefault("KANBOARD_PASSWORD", "admin")
	err := viper.ReadInConfig()
	fmt.Println(err)
}
