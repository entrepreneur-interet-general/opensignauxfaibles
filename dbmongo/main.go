package main

import (
	"fmt"

	"net/http"
	"time"

	"github.com/appleboy/gin-jwt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/spf13/viper"
)

var db = initDB()

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     checkOrigin,
}

func checkOrigin(r *http.Request) bool {
	return true
}

func wshandler(w http.ResponseWriter, r *http.Request, jwt string) {
	conn, err := wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Printf("Failed to set websocket upgrade: %+v", err)
		return
	}
	channel := make(chan socketMessage)
	addClientChannel <- channel

	for event := range channel {
		conn.WriteJSON(event)
	}

}

func main() {
	// Lancer Rserve en background

	// InitLogger(ioutil.Discard, os.Stdout, os.Stdout, os.Stderr)
	go r()
	go messageSocketAddClient()

	r := gin.Default()
	r.Use(gin.Recovery())
//	r.Use(Kanboard())
	// FIXME: configurer correctement CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:8080", "http://opensignauxfaibles.fr:8080", "http://opensignauxfaibles.fr:3000", "http://opensignauxfaibles.fr"}
	config.AddAllowHeaders("Authorization")
	config.AddAllowMethods("GET", "POST", "PUT", "HEAD", "DELETE")

	r.Use(cors.New(config))

	var identityKey = "id"
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:       "test zone",
		Key:         []byte("secret key"),
		Timeout:     time.Hour,
		MaxRefresh:  time.Hour,
		IdentityKey: identityKey,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			if v, ok := data.(*User); ok {
				return jwt.MapClaims{
					identityKey: v.UserName,
				}
			}
			return jwt.MapClaims{}
		},
		IdentityHandler: func(c *gin.Context) interface{} {
			claims := jwt.ExtractClaims(c)
			return &User{
				UserName: claims["id"].(string),
			}
		},
		Authenticator: func(c *gin.Context) (interface{}, error) {
			var loginVals login
			if err := c.ShouldBind(&loginVals); err != nil {
				return "", jwt.ErrMissingLoginValues
			}
			userID := loginVals.Username
			password := loginVals.Password

			if (userID == "admin" && password == "admin") || (userID == "test" && password == "test") {
				return &User{
					UserName:  userID,
					LastName:  "Bo-Yi",
					FirstName: "Wu",
				}, nil
			}
			return nil, jwt.ErrFailedAuthentication
		},
		Authorizator: func(data interface{}, c *gin.Context) bool {
			if v, ok := data.(*User); ok && v.UserName == "admin" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},

		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		panic("JWT Error:" + err.Error())
	}

	r.Use(static.Serve("/", static.LocalFile("static/", true)))

	r.POST("/login", authMiddleware.LoginHandler)

	r.GET("/ws/:jwt", func(c *gin.Context) {
		wshandler(c.Writer, c.Request, c.Params.ByName("jwt"))
	})

	api := r.Group("api")
	api.Use(authMiddleware.MiddlewareFunc())

	{
		api.GET("/refreshToken", authMiddleware.RefreshHandler)
		api.GET("/purge", purge)
		api.GET("/kanboard/get/projects", listProjects)
		api.GET("/kanboard/get/tasks", getKBTasks)
		api.POST("/admin/batch", upsertBatch)
		api.POST("/admin/batch/addFile", addFileToBatchHandler)
		api.GET("/admin/batch", listBatch)
		api.GET("/admin/files", adminFiles)
		api.GET("/admin/types", listTypes)
		api.GET("/admin/clone/:to", cloneDB)
		api.GET("/admin/features", adminFeature)
		api.GET("/admin/status", getDBStatus)
		api.GET("/admin/getLogs", getLogsHandler)
		api.GET("/batch/revert", revertBatchHandler)
		api.GET("/batch/next", nextBatchHandler)
		api.GET("/batch/purge", purgeBatchHandler)
		api.GET("/batch/process", processBatchHandler)
		api.POST("/admin/files", addFile)
		api.GET("/data/naf", getNAF)
		api.GET("/data/features", getFeatures)
		api.GET("/admin/epoch", epoch)
		api.GET("/data/prediction", predictionBrowse)
		api.GET("/data/etablissement/:batchKey/:siret", etablissementBrowseHandler)
		api.GET("/data/etablissement/:batchKey", etablissementBrowseHandler)
		api.GET("/import/:batch", importBatchHandler)
		api.GET("/compact/etablissement/:siret", compactEtablissementHandler)
		api.GET("/compact/etablissement", compactEtablissementHandler)
		api.GET("/compact/entreprise/:siren", compactEntrepriseHandler)
		api.GET("/compact/entreprise", compactEntrepriseHandler)
		api.GET("/reduce/:algo/:batch/:siret", reduceHandler)
		api.GET("/reduce/:algo/:batch", reduceHandler)
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
	viper.SetDefault("JWT_SECRET", "Secret à changer")
	viper.SetDefault("KANBOARD_ENDPOINT", "http://localhost/kanboard/jsonrpc.php")
	viper.SetDefault("KANBOARD_USERNAME", "admin")
	viper.SetDefault("KANBOARD_PASSWORD", "admin")
	err := viper.ReadInConfig()
	if err != nil {
		panic("Erreur à la lecture de la configuration")
	}
}
