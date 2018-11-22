package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/globalsign/mgo"
)

// Prediction prédiction
type Prediction struct {
	ID struct {
		Siret string `json:"siret" bson:"siret"`
		Batch string `json:"batch" bson:"batch"`
		Algo  string `json:"algo" bson:"algo"`
	} `json:"id" bson:"_id"`

	Prob          float64 `json:"prob" bson:"prob"`
	Diff          float64 `json:"diff" bson:"diff"`
	RaiSoc        string  `json:"raison_sociale" bson:"raison_sociale"`
	Departement   string  `json:"departement" bson:"departement"`
	Region        string  `json:"region" bson:"region"`
	EtatProCol    string  `json:"procol" bson:"procol"`
	DefaultUrssaf bool    `json:"default_urssaf" bson:"default_urssaf"`
	Connu         bool    `json:"connu" bson:"connu"`
	Niveau1       string  `json:"naf1" bson:"naf1"`
	Effectif      int     `json:"effectif" bson:"effectif"`
	CCSF          bool    `json:"ccsf" bson:"ccsf"`
}

func main() {
	csvFile, _ := os.Open("/home/christophe/Project/data-fake/fake-prediction.csv")

	predictionDict := make(map[string]Prediction)

	reader := csv.NewReader(bufio.NewReader(csvFile))
	reader.LazyQuotes = true
	reader.Comma = ';'
	reader.Read()
	var prediction []interface{}

	for {
		line, error := reader.Read()
		if error == io.EOF {
			break
		} else if error != nil {
			fmt.Println("niktoo")
			log.Fatal(error)
		}
		if line[0] != "" {
			proba, _ := strconv.ParseFloat(line[1], 64)
			diff, _ := strconv.ParseFloat(line[2], 64)
			effectif, _ := strconv.Atoi(line[11])
			var ccsf bool
			if line[12] == "" {
				ccsf = false
			} else {
				ccsf = true
			}

			var defaultUrssaf bool
			if line[7] == "TRUE" {
				defaultUrssaf = true
			} else {
				defaultUrssaf = false
			}

			var connu bool
			if line[8] == "0" {
				connu = false
			} else {
				connu = true
			}

			p := Prediction{
				ID: struct {
					Siret string `json:"siret" bson:"siret"`
					Batch string `json:"batch" bson:"batch"`
					Algo  string `json:"algo" bson:"algo"`
				}{
					Siret: line[0],
					Batch: "1802",
					Algo:  "algo1",
				},
				Prob:          proba,
				Diff:          diff,
				RaiSoc:        line[3],
				Departement:   line[4],
				Region:        line[5],
				EtatProCol:    line[6],
				DefaultUrssaf: defaultUrssaf,
				Connu:         connu,
				Niveau1:       line[10],
				Effectif:      effectif,
				CCSF:          ccsf,
			}

			predictionDict[line[0]] = p
		}
	}

	for _, v := range predictionDict {
		prediction = append(prediction, v)
	}

	mongodb, err := mgo.Dial("")
	if err != nil {
		fmt.Println("Insertion interrompue: " + err.Error())
		return
	}

	db := mongodb.DB("fakesignauxfaibles")

	err = db.C("Prediction").Insert(prediction...)

	if err != nil {
		fmt.Println("Insertion interrompue: " + err.Error())
		return
	}

	fmt.Println("Prédictions insérées: " + strconv.Itoa(len(prediction)))
}
