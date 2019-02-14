package engine

import (
	"dbmongo/lib/misc"
	"dbmongo/lib/naf"
	"errors"
	"io/ioutil"
	"regexp"

	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
)

func loadJSFunctions(path string) (map[string]bson.JavaScript, error) {
	files, err := ioutil.ReadDir(path)
	r := regexp.MustCompile(`(.*)\.js$`)

	functions := make(map[string]bson.JavaScript)

	for _, f := range files {
		if name := r.FindStringSubmatch(f.Name()); len(name) > 0 {
			b, err := ioutil.ReadFile(path + f.Name())
			if err == nil {
				functions[name[1]] = bson.JavaScript{
					Code: string(b),
				}
			}
		}
	}
	return functions, err
}

// PurgeBatch permet de supprimer un batch dans les objets de RawData
func PurgeBatch(batchKey string) error {

	functions, err := loadJSFunctions("js/purgeBatch/")
	if err != nil {
		return err
	}
	scope := bson.M{
		"currentBatch": batchKey,
		"f":            functions,
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "RawData"},
		Scope:    scope,
	}

	_, err = Db.DB.C("RawData").Find(nil).MapReduce(job, nil)
	return err
}

// Compact traite le compactage de la base RawData
func Compact() error {
	batches, _ := GetBatches()

	// Détermination scope traitement
	var completeTypes = make(map[string][]string)
	var batchesID []string

	for _, b := range batches {
		completeTypes[b.ID.Key] = b.CompleteTypes
		batchesID = append(batchesID, b.ID.Key)
	}

	functions, err := loadJSFunctions("js/compact/")
	if err != nil {
		return err
	}
	// Traitement MR
	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "RawData"},
		Scope: bson.M{
			"functions":     functions,
			"batches":       GetBatchesID(),
			"types":         GetTypes(),
			"completeTypes": completeTypes,
		},
	}

	_, err = Db.DB.C("RawData").Find(nil).MapReduce(job, nil)

	return err
}

// Reduce alimente la base Features
func Reduce(batchKey string, algo string, query interface{}) error {
	// éviter les noms d'algo essayant de pervertir l'exploration des fonctions
	isAlphaNum := regexp.MustCompile(`^[A-Za-z0-9]+$`).MatchString
	if !isAlphaNum(algo) {
		return errors.New("nom d'algorithme invalide, alphanumérique sans espace exigé")
	}

	functions, err := loadJSFunctions("js/reduce." + algo + "/")

	naf, err := naf.LoadNAF()
	if err != nil {
		return err
	}

	batch, err := GetBatch(batchKey)
	if err != nil {
		return err
	}

	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateDebut, batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf,
		"f":                      functions,
		"batches":                GetBatchesID(),
		"types":                  GetTypes(),
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"merge": "Features"},
		Scope:    scope,
	}

	_, err = Db.DB.C("RawData").Find(query).MapReduce(job, nil)

	return err
}

// Publish alimente la collection Public avec les objets destinés à la diffusion
func Public(batch AdminBatch) error {
	functions, err := loadJSFunctions("js/public/")

	scope := bson.M{
		"date_debut":             batch.Params.DateDebut,
		"date_fin":               batch.Params.DateFin,
		"date_fin_effectif":      batch.Params.DateFinEffectif,
		"serie_periode":          misc.GenereSeriePeriode(batch.Params.DateDebut, batch.Params.DateFin),
		"serie_periode_annuelle": misc.GenereSeriePeriodeAnnuelle(batch.Params.DateDebut, batch.Params.DateFin),
		"offset_effectif":        (batch.Params.DateFinEffectif.Year()-batch.Params.DateFin.Year())*12 + int(batch.Params.DateFinEffectif.Month()-batch.Params.DateFin.Month()),
		"actual_batch":           batch.ID.Key,
		"naf":                    naf.Naf,
		"f":                      functions,
		"batches":                GetBatchesID(),
		"types":                  GetTypes(),
	}

	job := &mgo.MapReduce{
		Map:      functions["map"].Code,
		Reduce:   functions["reduce"].Code,
		Finalize: functions["finalize"].Code,
		Out:      bson.M{"replace": "Features"},
		Scope:    scope,
	}
	// exécution

	_, err = Db.DB.C("DataRaw").Find(nil).MapReduce(job, nil)

	if err != nil {
		return errors.New("Erreur dans l'exécution des jobs MapReduce" + err.Error())
	}
	return nil
}

// BrowsePublic selectionne et retourne les objets de la collection Public
// Cette selection tient compte du scope et des tris demandés pour aggréger le résultat
func BrowsePublic(query interface{}) []Browseable {
	return []Browseable{
		Browseable{
			ID: struct {
				Key   string   `json:"key" bson:"key"`
				Scope []string `json:"scope" bson:"scope"`
			}{
				Key:   "test",
				Scope: []string{"test", "test2"},
			},
			Value: map[string]interface{}{
				"test": "test",
			},
		},
	}
}

// GetBatches retourne tous les objets AdminBatch de la base triés par ID
func GetBatches() ([]AdminBatch, error) {
	var batches []AdminBatch
	err := Db.DB.C("Admin").Find(bson.M{"_id.type": "batch"}).Sort("_id.key").All(&batches)
	return batches, err
}

// GetBatchesID retourne les clés des batches contenus en base
func GetBatchesID() []string {
	batches, _ := GetBatches()
	var batchesID []string
	for _, b := range batches {
		batchesID = append(batchesID, b.ID.Key)
	}
	return batchesID
}

// GetBatch retourne le batch correspondant à la clé batchKey
func GetBatch(batchKey string) (AdminBatch, error) {
	var batch AdminBatch
	err := Db.DB.C("Admin").Find(bson.M{"_id.type": "batch", "_id.key": batchKey}).One(&batch)
	return batch, err
}

// Purge réinitialise la base, à utiliser avec modération
func Purge() interface{} {
	infoRawData, errRawData := Db.DB.C("RawData").RemoveAll(nil)
	infoJournal, errJournal := Db.DB.C("Journal").RemoveAll(nil)
	infoFeatures, errFeatures := Db.DB.C("Features").RemoveAll(nil)
	infoPublic, errPublic := Db.DB.C("Public").RemoveAll(nil)

	returnData := map[string]map[string]interface{}{
		"RawData": map[string]interface{}{
			"info":  infoRawData,
			"error": errRawData,
		},
		"Journal": map[string]interface{}{
			"info":  infoJournal,
			"error": errJournal,
		},
		"Features": map[string]interface{}{
			"info":  infoFeatures,
			"error": errFeatures,
		},
		"Public": map[string]interface{}{
			"info":  infoPublic,
			"error": errPublic,
		},
	}

	return returnData
}
