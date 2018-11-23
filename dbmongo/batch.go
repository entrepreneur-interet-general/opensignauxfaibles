package main

import (
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/davecgh/go-spew/spew"
	"github.com/gin-gonic/gin"
	"github.com/globalsign/mgo/bson"
)

// AdminBatch metadata Batch
type AdminBatch struct {
	ID            AdminID    `json:"id" bson:"_id"`
	Files         BatchFiles `json:"files" bson:"files"`
	Readonly      bool       `json:"readonly" bson:"readonly"`
	CompleteTypes []string   `json:"complete_types" bson:"complete_types"`
	Params        struct {
		DateDebut       time.Time `json:"date_debut" bson:"date_debut"`
		DateFin         time.Time `json:"date_fin" bson:"date_fin"`
		DateFinEffectif time.Time `json:"date_fin_effectif" bson:"date_fin_effectif"`
	} `json:"params" bson:"param"`
}

// BatchFiles fichiers mappés par type
type BatchFiles map[string][]string

func (batchFiles BatchFiles) attachFile(fileType string, file string) {
	batchFiles[fileType] = append(batchFiles[fileType], file)
}

func isBatchID(batchID string) bool {
	_, err := time.Parse("0601", batchID)
	return err == nil
}

func (batch *AdminBatch) load(batchKey string) error {
	err := db.DB.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(batch)
	return err
}

func (batch *AdminBatch) save() error {
	_, err := db.DB.C("Admin").Upsert(bson.M{"_id": batch.ID}, batch)
	return err
}

func (batch *AdminBatch) new(batchID string) error {
	if batchID == "" {
		return errors.New("Valeur de batch non autorisée")
	}
	batch.ID.Key = batchID
	batch.ID.Type = "batch"
	batch.Files = BatchFiles{}
	return nil
}

func nextBatchHandler(c *gin.Context) {
	err := nextBatch()
	if err != nil {
		c.JSON(500, fmt.Errorf("Erreur nextBatch: "+err.Error()))
	}
	batches, _ := getBatches()
	mainMessageChannel <- socketMessage{
		Batches: batches,
	}
	c.JSON(200, "nextBatch ok")
}

func nextBatch() error {
	batch := lastBatch()
	// spew.Dump(batch)
	newBatchID, err := nextBatchID(batch.ID.Key)
	if err != nil {
		return fmt.Errorf("Mauvais numéro de batch: " + err.Error())
	}
	newBatch := AdminBatch{
		ID: AdminID{
			Key:  newBatchID,
			Type: "batch",
		},
		CompleteTypes: batch.CompleteTypes,
	}

	batch.Readonly = true

	err = batch.save()
	if err != nil {
		return fmt.Errorf("Erreur readonly Batch: " + err.Error())
	}

	err = newBatch.save()
	if err != nil {
		return fmt.Errorf("Erreur newBatch: " + err.Error())
	}
	return nil
}

func nextBatchID(batchID string) (string, error) {
	batchTime, err := time.Parse("0601", batchID)
	if err != nil {
		return "", err
	}
	nextBatchTime := time.Date(batchTime.Year(), time.Month(batchTime.Month()+1), 1, 0, 0, 0, 0, time.UTC)
	return nextBatchTime.Format("0601"), err
}

func sp(s string) *string {
	return &s
}

func upsertBatch(c *gin.Context) {
	status := db.Status

	batch := AdminBatch{}
	err := c.Bind(&batch)
	if err != nil {
		c.JSON(500, err)
		return
	}

	err = batch.save()
	if err != nil {
		c.JSON(500, "Erreur à l'enregistrement")
		return
	}

	batches, _ := getBatches()
	mainMessageChannel <- socketMessage{
		Batches: batches,
	}

	status.Epoch++
	status.write()

	c.JSON(200, batch)
}

func listBatch(c *gin.Context) {
	var batch []AdminBatch
	err := db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("-_id.key").All(&batch)
	if err != nil {
		spew.Dump(err)
		c.JSON(500, err)
		return
	}
	c.JSON(200, batch)
}

func getBatchesID() []string {
	batches, _ := getBatches()
	var batchesID []string
	for _, b := range batches {
		batchesID = append(batchesID, b.ID.Key)
	}
	return batchesID
}

func getBatches() ([]AdminBatch, error) {
	var batches []AdminBatch
	err := db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batches)
	return batches, err
}

// getBatch retourne le batch correspondant à la clé batchKey
func getBatch(batchKey string) (AdminBatch, error) {
	var batch AdminBatch
	err := db.DB.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(&batch)
	return batch, err
}

// batchToTime calcule la date de référence à partir de la référence de batch
func batchToTime(batch string) (time.Time, error) {
	year, err := strconv.Atoi(batch[0:2])
	if err != nil {
		return time.Time{}, err
	}

	month, err := strconv.Atoi(batch[2:4])
	if err != nil {
		return time.Time{}, err
	}

	date := time.Date(2000+year, time.Month(month), 1, 0, 0, 0, 0, time.UTC)
	return date, err
}

func processBatchHandler(c *gin.Context) {
	go func() {
		processBatch()
	}()
	c.JSON(200, "ok !")
}

func processBatch() {
	log("info", "processBatch", "Lancement de l'intégration du batch")
	status := db.Status
	batch := lastBatch()
	status.setDBStatus(sp("Import des fichiers"))
	importBatch(&batch)
	compactEntreprise("")
	compactEtablissement("")
	// for _, algo := range []string{"algo1", "algo2"} {
	// 	_, err := reduce(batch, algo, "")
	// 	fmt.Println(err)
	// }
	status.setDBStatus(nil)
}

func lastBatch() AdminBatch {
	batches, _ := getBatches()
	l := len(batches)
	batch := batches[l-1]
	return batch
}

func createNextBatch() error {
	batchID, _ := nextBatchID(lastBatch().ID.Key)
	batch := AdminBatch{
		ID: AdminID{
			Key:  batchID,
			Type: "batch",
		},
	}
	err := batch.save()
	return err
}

type newFile struct {
	FileName string `json:"filename"`
	Type     string `json:"type"`
	BatchKey string `json:"batch"`
}

func addFileToBatchHandler(c *gin.Context) {
	var file newFile
	err := c.Bind(&file)
	if err != nil {
		c.JSON(500, err.Error())
	}
	addFileChannel <- file

	c.JSON(200, "Demande d'ajout prise en compte")
}

func addFileToBatch() chan newFile {
	channel := make(chan newFile)

	go func() {
		for file := range channel {
			batch, _ := getBatch(file.BatchKey)
			batch.Files[file.Type] = append(batch.Files[file.Type], file.FileName)
			batch.save()
			batches, _ := getBatches()
			db.Status.Epoch++
			db.Status.write()
			mainMessageChannel <- socketMessage{
				JournalEvent: log(info, "addFileToBatch", "Fichier "+file.FileName+"du type "+file.Type+" ajouté au batch "+file.BatchKey),
				Batches:      batches,
			}
		}
	}()

	return channel
}

func purgeBatchHandler(c *gin.Context) {
	err := purgeBatch()
	if err != nil {
		c.JSON(500, "Erreur dans la purge du batch: "+err.Error())
	}
}

func purgeBatch() error {
	batch := lastBatch()

	// prepareMRJob charge les fichiers MapReduce et fournit les paramètres pour l'exécution
	MREntreprise, errEn := loadMR("purgeBatch", "entreprise")
	MREtablissement, errEt := loadMR("purgeBatch", "etablissement")

	if errEn != nil || errEt != nil {
		return fmt.Errorf("Erreur de chargement MapReduce")
	}

	MREntreprise.Out = &bson.M{"replace": "Entreprise"}
	MREtablissement.Out = &bson.M{"replace": "Etablissement"}

	MREntreprise.Scope = &bson.M{
		"currentBatch": batch.ID.Key,
	}
	MREtablissement.Scope = &bson.M{
		"currentBatch": batch.ID.Key,
	}

	_, errEn = db.DB.C("Entreprise").Find(nil).MapReduce(MREntreprise, nil)
	_, errEt = db.DB.C("Etablissement").Find(nil).MapReduce(MREtablissement, nil)

	return nil
}

func revertBatchHandler(c *gin.Context) {
	err := revertBatch()
	if err != nil {
		c.JSON(500, err)
	}
	batches, _ := getBatches()
	mainMessageChannel <- socketMessage{
		Batches: batches,
	}
	c.JSON(200, "ok")
}

func dropLastBatch() error {
	batch := lastBatch()
	_, err := db.DB.C("Admin").RemoveAll(bson.M{"_id.key": batch.ID.Key, "_id.type": "batch"})
	return err
}

// revertBatch purge le batch et supprime sa référence dans la collection Admin
func revertBatch() error {
	err := purgeBatch()
	if err != nil {
		return fmt.Errorf("Erreur lors de la purge: " + err.Error())
	}
	err = dropLastBatch()
	if err != nil {
		return fmt.Errorf("Erreur lors de la purge: " + err.Error())
	}

	return nil
}

var addFileChannel = addFileToBatch()
