package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/cnf/structhash"
	"github.com/spf13/viper"
)

// CCSF information urssaf ccsf
type CCSF struct {
	NumeroCompte   string    `json:"-" bson:"-"`
	DateTraitement time.Time `json:"date_traitement" bson:"date_traitement"`
	Stade          string    `json:"stade" bson:"stade"`
	Action         string    `json:"action" json:"action"`
	DateBatch      time.Time `json:"date_batch" bson:"date_batch"`
}

func parseCCSF(path string, dateBatch time.Time) chan *CCSF {
	outputChannel := make(chan *CCSF)

	file, err := os.Open(viper.GetString("APP_DATA") + path)
	if err != nil {
		fmt.Println("Error", err)
	}

	reader := csv.NewReader(bufio.NewReader(file))
	reader.Comma = ';'
	reader.Read()

	f := map[string]int{
		"NumeroCompte":   0,
		"DateTraitement": 1,
		"Stade":          2,
		"Action":         3,
	}

	go func() {
		for {
			r, error := reader.Read()
			if error == io.EOF {
				break
			} else if error != nil {
				// log.Fatal(error)
			}

			ccsf := CCSF{}
			ccsf.Action = r[f["Action"]]
			ccsf.Stade = r[f["Stade"]]
			ccsf.DateTraitement, err = UrssafToDate(r[f["DateTraitement"]])
			ccsf.NumeroCompte = r[f["NumeroCompte"]]
			ccsf.DateBatch = dateBatch
			outputChannel <- &ccsf
		}
		close(outputChannel)
		file.Close()
	}()
	return outputChannel
}

func importCCSF(batch *AdminBatch) error {

	mapping, _ := getCompteSiretMapping(batch)

	dateBatch, errDate := batchToTime(batch.ID.Key)
	if errDate != nil {
		return errDate
	}

	for _, data := range batch.Files["ccsf"] {
		for ccsf := range parseCCSF(data, dateBatch) {
			if siret, ok := mapping[ccsf.NumeroCompte]; ok {
				hash := fmt.Sprintf("%x", structhash.Md5(ccsf, 1))

				value := ValueEtablissement{
					Value: Etablissement{
						Siret: siret,
						Batch: map[string]Batch{
							batch.ID.Key: Batch{
								CCSF: map[string]*CCSF{
									hash: ccsf,
								}}}}}
				db.ChanEtablissement <- &value
			}
		}
	}
	db.ChanEtablissement <- &ValueEtablissement{}
	return nil
}
