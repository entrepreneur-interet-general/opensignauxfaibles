package main

import (
	"fmt"
	"io/ioutil"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func adminFiles(c *gin.Context) {
	basePath := viper.GetString("APP_DATA")
	files, err := listFiles(basePath)
	if err != nil {
		c.JSON(500, err)
	} else {
		c.JSON(200, files)
	}
}

type fileSummary struct {
	Name string    `json:"name" bson:"name"`
	Size int64     `json:"size" bson:"size"`
	Date time.Time `json:"date" bson:"date"`
}

func listFiles(basePath string) ([]fileSummary, error) {
	var files []fileSummary
	basePathConf := viper.GetString("APP_DATA")
	b := len(basePathConf)

	currentFiles, err := ioutil.ReadDir(basePath)
	if err != nil {
		return []fileSummary{}, err
	}

	for _, file := range currentFiles {
		if file.IsDir() {
			subPath := fmt.Sprintf("%s/%s", basePath, file.Name())
			subFiles, err := listFiles(subPath)
			if err != nil {
				return []fileSummary{}, err
			}
			files = append(files, subFiles...)
		} else {
			files = append(files, fileSummary{
				Name: fmt.Sprintf("%s/%s", basePath, file.Name())[b:],
				Size: file.Size(),
				Date: file.ModTime(),
			})
		}
	}
	return files, nil
}

var importFunctions = map[string]func(*AdminBatch) error{
  "apconso":    importAPConso,
  "bdf":        importBDF,
  "delai":      importDelai,
  "apdemande":  importAPDemande,
  "diane":      importDiane,
  "cotisation": importCotisation,
  "dpae":       importDPAE,
  "altares":    importAltares,
  "procol":   importProcol,
  "ccsf":       importCCSF,
  "debit":      importDebit,
  "effectif":   importEffectif,
  "sirene":     importSirene,
}

func purge(c *gin.Context) {
	db.DB.C("Etablissement").RemoveAll(nil)
	db.DB.C("Entreprise").RemoveAll(nil)
	db.DB.C("Prediction").RemoveAll(nil)
	c.String(200, "Done")
}

func importBatchHandler(c *gin.Context) {
	batchKey := c.Params.ByName("batch")
	batch := AdminBatch{}
	batch.load(batchKey)
	go importBatch(&batch)
}
func importBatch(batch *AdminBatch) {
    if !batch.Readonly {
        for fnName, fn := range importFunctions {
            log(info, "importMain", "Début de l'import du type: "+fnName+" pour le batch "+batch.ID.Key)
            err := fn(batch)
            if err != nil {
                log(critical, "importMain", "Erreur à l'importation du type: "+fnName)
            }
            log(info, "importMain", "Fin de l'import du type: "+fnName+" pour le batch "+batch.ID.Key)
        }
    } else {
        log(critical, "importMain", "Le lot "+batch.ID.Key+" est fermé, import impossible.")
    }
}
