package main

import (
	"dbmongo/lib/altares"
	"dbmongo/lib/apartconso"
	"dbmongo/lib/apartdemande"
	"dbmongo/lib/bdf"
	"dbmongo/lib/diane"
	"dbmongo/lib/engine"
	"dbmongo/lib/files"
	"dbmongo/lib/interim"
	"dbmongo/lib/sirene"
	"dbmongo/lib/urssaf"
	"fmt"
	"io"
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
	"github.com/spf13/viper"
)

//
func adminFeature(c *gin.Context) {
	c.JSON(200, []string{"algo1", "algo2"})
}

//
func listTypesHandler(c *gin.Context) {
	c.JSON(200, engine.GetTypes)
}

//
func addFile(c *gin.Context) {
	file, err := c.FormFile("file")
	batch := c.Request.FormValue("batch")
	fileType := c.Request.FormValue("type")

	source, err := file.Open()
	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	defer source.Close()

	os.MkdirAll(viper.GetString("APP_DATA")+"/"+batch+"/"+fileType+"/", os.ModePerm)
	destination, err := os.Create(viper.GetString("APP_DATA") + "/" + batch + "/" + fileType + "/" + file.Filename)
	if err != nil {
		c.JSON(500, err.Error())
		return
	}
	defer destination.Close()

	_, err = io.Copy(destination, source)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	basePath := viper.GetString("APP_DATA")
	newFiles, err := files.ListFiles(basePath)

	if err != nil {
		c.JSON(500, err.Error())
		return
	}

	engine.MainMessageChannel <- engine.SocketMessage{
		Files: newFiles,
	}

	c.JSON(200, "ok")
}

//
func nextBatchHandler(c *gin.Context) {
	err := engine.NextBatch()
	if err != nil {
		c.JSON(500, fmt.Errorf("Erreur lors de la création du batch suivant: "+err.Error()))
	}
	batches, _ := engine.GetBatches()
	engine.MainMessageChannel <- engine.SocketMessage{
		Batches: batches,
	}
	c.JSON(200, batches)
}

func sp(s string) *string {
	return &s
}

//
func upsertBatchHandler(c *gin.Context) {
	var batch engine.AdminBatch
	err := c.Bind(&batch)
	if err != nil {
		c.JSON(400, err.Error)
		return
	}

	err = batch.Save()
	if err != nil {
		c.JSON(500, "Erreur à l'enregistrement: "+err.Error())
		return
	}

	batches, _ := engine.GetBatches()
	engine.MainMessageChannel <- engine.SocketMessage{
		Batches: batches,
	}

	c.JSON(200, batch)
}

//
func listBatchHandler(c *gin.Context) {
	var batch []engine.AdminBatch
	err := engine.Db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("-_id.key").All(&batch)
	if err != nil {
		spew.Dump(err)
		c.JSON(500, err)
		return
	}
	c.JSON(200, batch)
}

//
func processBatchHandler(c *gin.Context) {
	var query struct {
		Batches []string `json:"batches"`
	}
	err := c.ShouldBind(query)
	if err != nil {
		c.JSON(400, err.Error)
	}

	// TODO: valider que tous les batches demandés existent
	err = engine.ProcessBatch(query.Batches)
	if err != nil {
		c.JSON(500, err.Error())
	}
	c.JSON(200, "ok !")
}

//
func purgeBatchHandler(c *gin.Context) {
	batch := engine.LastBatch()
	err := engine.PurgeBatch(batch.ID.Key)

	if err != nil {
		c.JSON(500, "Erreur dans la purge du batch: "+err.Error())
	} else {
		c.JSON(200, "ok")
	}
}

func revertBatchHandler(c *gin.Context) {
	err := engine.RevertBatch()
	if err != nil {
		c.JSON(500, err)
	}
	batches, _ := engine.GetBatches()
	engine.MainMessageChannel <- engine.SocketMessage{
		Batches: batches,
	}
	c.JSON(200, "ok")
}

// RegisteredParsers liste des parsers disponibles
var registeredParsers = map[string]engine.Parser{
	"urssaf":    urssaf.Parser,
	"apconso":   apartconso.Parser,
	"apdemande": apartdemande.Parser,
	"bdf":       bdf.Parser,
	"altares":   altares.Parser,
	"sirene":    sirene.Parser,
	"diane":     diane.Parser,
	"interim":   interim.Parser,
}

//
func adminFilesHandler(c *gin.Context) {
	basePath := viper.GetString("APP_DATA")
	files, err := files.ListFiles(basePath)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, files)
	}
}

//
func importBatchHandler(c *gin.Context) {
	var params struct {
		BatchKey string   `json:"batch"`
		Parsers  []string `json:"parsers"`
	}
	err := c.ShouldBind(&params)
	if err != nil {
		c.JSON(400, "Requête malformée: "+err.Error())
	}
	batch := engine.AdminBatch{}
	batch.Load(params.BatchKey)

	var parsers []engine.Parser
	if params.Parsers == nil {
		for _, f := range registeredParsers {
			parsers = append(parsers, f)
		}
	} else {
		for _, p := range params.Parsers {
			if f, ok := registeredParsers[p]; ok {
				parsers = append(parsers, f)
			} else {
				c.JSON(404, p+" n'est pas un parser reconnu.")
				return
			}
		}
	}

	engine.ImportBatch(batch, parsers)
}

//
func eventsHandler(c *gin.Context) {
	logs, err := engine.GetEventsFromDB(nil, 250)
	if err != nil {
		c.JSON(500, err.Error())
	} else {
		c.JSON(200, logs)
	}
}

func deleteHandler(c *gin.Context) {
	var result []interface{}
	engine.Db.DB.C("RawData").RemoveAll(bson.M{"_id": bson.M{"$type": "objectId"}})
	c.JSON(200, result)
}

func getTasksHandler(c *gin.Context) {
	c.JSON(501, "Not implemented (for the moment)")
}

func browsePublicHandler(c *gin.Context) {
	data := engine.BrowsePublic(nil)
	c.JSON(200, data)
}
