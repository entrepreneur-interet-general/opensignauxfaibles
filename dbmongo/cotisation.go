package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/viper"

	"github.com/cnf/structhash"
)

// Cotisation Cotisation – fichier Urssaf
type Cotisation struct {
	NumeroCompte string  `json:"numero_compte" bson:"numero_compte"`
	PeriodeDebit string  `json:"periode_debit" bson:"periode_debit"`
	Periode      Periode `json:"period" bson:"periode"`
	Recouvrement float64 `json:"recouvrement" bson:"recouvrement"`
	Encaisse     float64 `json:"encaisse" bson:"encaisse"`
	Du           float64 `json:"du" bson:"du"`
	Ecriture     string  `json:"ecriture" bson:"ecriture"`
}

func parseCotisation(paths []string) chan *Cotisation {
	outputChannel := make(chan *Cotisation)

	field := map[string]int{
		"NumeroCompte": 0,
		"PeriodeDebit": 1,
		"Periode":      4,
		"Recouvrement": 2,
		"Encaisse":     3,
		"Du":           5,
		"Ecriture":     6,
	}

	go func() {
		for _, path := range paths {
			file, err := os.Open(viper.GetString("APP_DATA") + path)
			if err != nil {
				fmt.Println("Error", err)
			}

			reader := csv.NewReader(bufio.NewReader(file))
			reader.Comma = ';'
			reader.Read()

			for {

				row, error := reader.Read()
				if error == io.EOF {
					break
				} else if error != nil {
					// log.Fatal(error)
				}

				cotisation := Cotisation{}
				cotisation.NumeroCompte = row[field["NumeroCompte"]]
				cotisation.Periode, err = UrssafToPeriod(row[field["Periode"]])
				cotisation.PeriodeDebit = row[field["PeriodeDebit"]]
				cotisation.Recouvrement, err = strconv.ParseFloat(strings.Replace(row[field["Recouvrement"]], ",", ".", -1), 64)
				cotisation.Encaisse, err = strconv.ParseFloat(strings.Replace(row[field["Encaisse"]], ",", ".", -1), 64)
				cotisation.Du, err = strconv.ParseFloat(strings.Replace(row[field["Du"]], ",", ".", -1), 64)
				cotisation.Ecriture = row[field["Ecriture"]]

				outputChannel <- &cotisation
			}
			file.Close()
		}

		close(outputChannel)
	}()
	return outputChannel
}

func importCotisation(batch *AdminBatch) error {
	mapping, _ := getCompteSiretMapping(batch)

	for cotisation := range parseCotisation(batch.Files["cotisation"]) {
		if siret, ok := mapping[cotisation.NumeroCompte]; ok {
			hash := fmt.Sprintf("%x", structhash.Md5(*cotisation, 1))

			value := ValueEtablissement{
				Value: Etablissement{
					Siret: siret,
					Batch: map[string]Batch{
						batch.ID.Key: Batch{
							Cotisation: map[string]*Cotisation{
								hash: cotisation,
							}}}}}
			db.ChanEtablissement <- &value
		}
	}
	db.ChanEtablissement <- &ValueEtablissement{}
	return nil
}
