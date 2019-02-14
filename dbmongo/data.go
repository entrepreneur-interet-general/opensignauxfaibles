package main

import (
	"dbmongo/lib/engine"
	"dbmongo/lib/naf"

	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

//
// @summary Lance un traitement de réduction
// @description Alimente la collection Features
// @Tags Traitements
// @accept  json
// @produce  json
// @Param algo query string true "Identifiant du traitement"
// @Param batch query string true "Identifier du batch"
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/data/reduce/{algo}/{batch}/{siret} [get]
func reduceHandler(c *gin.Context) {
	var params struct {
		BatchKey string `json:"batch"`
		Algo     string `json:"features"`
		Key      string `json:"key"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, err.Error())
	}

	var query bson.M
	if params.Key != "" {
		query = bson.M{"_id": query}
	} else {
		query = bson.M{"value.index." + params.Algo: true}
	}

	err = engine.Reduce(params.BatchKey, params.Algo, params.Key)
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, "Traitement effectué")
	}
}

// func reduce(algo string, batchKey string, key string) error {
// 	// éviter les noms d'algo essayant de pervertir l'exploration des fonctions
// 	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
// 	if !isAlphaNum(algo) {
// 		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
// 	}
// 	c.ShouldBind(params)

// 	err := engine.Reduce(params.BatchKey, params.Algo, nil)
// 	if err != nil {
// 		c.JSON(500, err.Error())
// 	} else {
// 		c.JSON(200, "Traitement effectué")
// 	}
// }

//
// @summary Lance un traitement de compactage
// @description Alimente la collection Features
// @Tags Traitements
// @accept  json
// @produce  json
// @Param algo query string true "Identifiant du traitement"
// @Param batch query string true "Identifier du batch"
// @Success 200 {string} string ""
// @Router /api/data/compact [get]
// @Security ApiKeyAuth
func compactHandler(c *gin.Context) {
	err := engine.Compact()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	c.JSON(200, "ok")
}

//
// @summary Descriptif NAF
// @description Liste tous les codes NAF, les descriptions des codes NAF et les liens entre le niveau 1 et le niveau 5
// @Tags Traitements
// @accept  json
// @produce  json
// @Param algo query string true "Identifiant du traitement"
// @Param batch query string true "Identifier du batch"
// @Success 200 {string} string ""
// @Router /api/data/compact [get]
// @Security ApiKeyAuth
func nafHandler(c *gin.Context) {
	c.JSON(200, naf.Naf)
}

//
// @summary Purge la collection RawData
// @description Suppression de tous les objets de données brutes contenus dans la collection RawData (irréversible)
// @Tags Traitements
// @accept  json
// @produce  json
// @Security ApiKeyAuth
// @Success 200 {string} string ""
// @Router /api/data/purge [get]
func purgeHandler(c *gin.Context) {
	info := engine.Purge()
	c.JSON(200, info)
}
